package game

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"battle-tiles/internal/conf"
	model "battle-tiles/internal/dal/model/game"
	"battle-tiles/internal/infra"

	"github.com/go-kratos/kratos/v2/log"
)

// BattleSyncManager 管理所有会话的战绩同步
type BattleSyncManager struct {
	mu              sync.RWMutex
	syncers         map[string]*battleSyncer // key: "userID:houseGID"
	battleUC        *BattleRecordUseCase     // 使用 UseCase 统一处理战绩逻辑
	data            *infra.Data              // 用于记录同步日志
	logger          *log.Helper
	syncInterval    time.Duration // 战绩同步间隔
	compensateEvery int           // 每隔多少次同步执行一次补偿
}

// NewBattleSyncManager 创建战绩同步管理器
func NewBattleSyncManager(battleUC *BattleRecordUseCase, data *infra.Data, syncConf *conf.Sync, logger log.Logger) *BattleSyncManager {
	// 默认值
	syncInterval := 10 * time.Second
	compensateEvery := 30

	if syncConf != nil {
		if syncConf.BattleSyncInterval > 0 {
			syncInterval = time.Duration(syncConf.BattleSyncInterval) * time.Second
		}
		if syncConf.BalanceCompensateEvery > 0 {
			compensateEvery = int(syncConf.BalanceCompensateEvery)
		}
	}

	return &BattleSyncManager{
		syncers:         make(map[string]*battleSyncer),
		battleUC:        battleUC,
		data:            data,
		logger:          log.NewHelper(logger),
		syncInterval:    syncInterval,
		compensateEvery: compensateEvery,
	}
}

// StartSync 启动战绩同步
func (m *BattleSyncManager) StartSync(ctx context.Context, userID int, houseGID int) {
	key := fmt.Sprintf("%d:%d", userID, houseGID)

	m.mu.Lock()
	defer m.mu.Unlock()

	// 如果已经存在，先停止
	if syncer, exists := m.syncers[key]; exists {
		syncer.stop()
		delete(m.syncers, key)
	}

	// 创建新的同步器，传入带 platform 的 context 和配置
	syncer := newBattleSyncer(ctx, userID, houseGID, m.battleUC, m.data, m.logger, m.syncInterval, m.compensateEvery)
	m.syncers[key] = syncer
	syncer.start()

	if syncer.verboseTaskLog {
		m.logger.Infof("Started battle syncer for user %d, house %d", userID, houseGID)
	}
}

// StopSync 停止战绩同步
func (m *BattleSyncManager) StopSync(userID int, houseGID int) {
	key := fmt.Sprintf("%d:%d", userID, houseGID)

	m.mu.Lock()
	defer m.mu.Unlock()

	if syncer, exists := m.syncers[key]; exists {
		verbose := syncer.verboseTaskLog
		syncer.stop()
		delete(m.syncers, key)
		if verbose {
			m.logger.Infof("Stopped battle syncer for user %d, house %d", userID, houseGID)
		}
	}
}

// StopAll 停止所有战绩同步
func (m *BattleSyncManager) StopAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for key, syncer := range m.syncers {
		syncer.stop()
		delete(m.syncers, key)
	}

	m.logger.Info("Stopped all battle syncers")
}

// battleSyncer 单个会话的战绩同步器
type battleSyncer struct {
	ctx             context.Context // 保存带 platform 的 context
	userID          int
	houseGID        int
	battleUC        *BattleRecordUseCase // 使用 UseCase 处理战绩
	data            *infra.Data          // 用于记录同步日志
	logger          *log.Helper
	stopChan        chan struct{}
	wg              sync.WaitGroup
	syncInterval    time.Duration
	isFirstSync     bool
	sessionID       int32 // 会话 ID，用于记录同步日志
	verboseTaskLog  bool  // 是否显示详细的任务日志
	syncCount       int   // 同步计数器，用于触发补偿任务
	compensateEvery int   // 每隔多少次同步执行一次补偿
}

func newBattleSyncer(ctx context.Context, userID int, houseGID int, battleUC *BattleRecordUseCase, data *infra.Data, logger *log.Helper, syncInterval time.Duration, compensateEvery int) *battleSyncer {
	return &battleSyncer{
		ctx:             ctx, // 保存 context
		userID:          userID,
		houseGID:        houseGID,
		battleUC:        battleUC,
		data:            data,
		logger:          logger,
		stopChan:        make(chan struct{}),
		syncInterval:    syncInterval, // 从配置读取
		isFirstSync:     true,
		verboseTaskLog:  battleUC.verboseTaskLog, // 使用 UseCase 的配置
		syncCount:       0,
		compensateEvery: compensateEvery, // 从配置读取
	}
}

func (s *battleSyncer) start() {
	if s.verboseTaskLog {
		s.logger.Infof("Starting battle syncer for house %d, isFirstSync=%v", s.houseGID, s.isFirstSync)
	}
	s.wg.Add(1)
	go s.syncLoop()
}

func (s *battleSyncer) stop() {
	close(s.stopChan)
	s.wg.Wait()
}

func (s *battleSyncer) syncLoop() {
	defer s.wg.Done()

	ticker := time.NewTicker(s.syncInterval)
	defer ticker.Stop()

	// 立即执行第一次同步
	s.syncOnce()

	for {
		select {
		case <-ticker.C:
			s.syncOnce()
		case <-s.stopChan:
			return
		}
	}
}

func (s *battleSyncer) syncOnce() {
	startTime := time.Now()

	// 获取 session_id（用于记录同步日志）
	// 第一次同步获取最近1小时的数据，之后获取最近3分钟的数据
	typeid := 0 // 最近3分钟
	typeDesc := "last 3 minutes"
	if s.isFirstSync {
		typeid = 2 // 最近1小时
		typeDesc = "last 1 hour"
		s.isFirstSync = false
	}

	if s.verboseTaskLog {
		s.logger.Infof("Fetching battle records for house %d (%s)...", s.houseGID, typeDesc)
	}

	// 使用 BattleRecordUseCase.PullAndSave 统一处理战绩同步
	httpClient := &http.Client{Timeout: 10 * time.Second}
	baseURL := "http://phone2.foxuc.com/Ashx/GroService.ashx"

	// groupID 参数应该传 houseGID，因为 foxuc API 的 groupid 就是店铺 GID
	saved, err := s.battleUC.PullAndSave(s.ctx, httpClient, baseURL, s.houseGID, s.houseGID, typeid)
	if err != nil {
		s.logger.Errorf("Failed to sync battle records for house %d: %v", s.houseGID, err)
		// 记录失败日志
		s.recordSyncLog(s.ctx, startTime, 0, model.SyncStatusFailed, err.Error())
		return
	}

	if saved > 0 && s.verboseTaskLog {
		s.logger.Infof("Synced %d battle records for house %d", saved, s.houseGID)
	}

	// 记录成功日志
	s.recordSyncLog(s.ctx, startTime, int32(saved), model.SyncStatusSuccess, "")

	// 定期执行补偿任务（每 compensateEvery 次同步执行一次）
	s.syncCount++
	if s.syncCount >= s.compensateEvery {
		s.syncCount = 0
		go s.runCompensate()
	}
}

// runCompensate 执行补偿任务
func (s *battleSyncer) runCompensate() {
	if s.verboseTaskLog {
		s.logger.Infof("Running compensate task for house %d...", s.houseGID)
	}

	compensated, err := s.battleUC.CompensateUnSettledBattles(s.ctx, int32(s.houseGID))
	if err != nil {
		s.logger.Errorf("Compensate task failed for house %d: %v", s.houseGID, err)
		return
	}

	if compensated > 0 {
		s.logger.Infof("Compensate task completed for house %d: %d records", s.houseGID, compensated)
	}
}

// recordSyncLog 记录同步日志
func (s *battleSyncer) recordSyncLog(ctx context.Context, startTime time.Time, recordsSynced int32, status string, errorMsg string) {
	// 如果没有 sessionID，跳过日志记录
	if s.sessionID == 0 || status == "success" {
		return
	}

	completedAt := time.Now()
	syncLog := &model.GameSyncLog{
		SessionID:     s.sessionID,
		SyncType:      model.SyncTypeBattleRecord,
		Status:        status,
		RecordsSynced: recordsSynced,
		ErrorMessage:  errorMsg,
		StartedAt:     startTime,
		CompletedAt:   &completedAt,
	}

	if err := s.data.GetDBWithContext(ctx).Create(syncLog).Error; err != nil {
		s.logger.Errorf("Failed to create sync log: %v", err)
	}
}
