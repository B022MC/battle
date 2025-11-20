package game

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	model "battle-tiles/internal/dal/model/game"
	"battle-tiles/internal/infra"

	"github.com/go-kratos/kratos/v2/log"
)

// BattleSyncManager 管理所有会话的战绩同步
type BattleSyncManager struct {
	mu       sync.RWMutex
	syncers  map[string]*battleSyncer // key: "userID:houseGID"
	battleUC *BattleRecordUseCase     // 使用 UseCase 统一处理战绩逻辑
	data     *infra.Data              // 用于记录同步日志
	logger   *log.Helper
}

// NewBattleSyncManager 创建战绩同步管理器
func NewBattleSyncManager(battleUC *BattleRecordUseCase, data *infra.Data, logger log.Logger) *BattleSyncManager {
	return &BattleSyncManager{
		syncers:  make(map[string]*battleSyncer),
		battleUC: battleUC,
		data:     data,
		logger:   log.NewHelper(logger),
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

	// 创建新的同步器，传入带 platform 的 context
	syncer := newBattleSyncer(ctx, userID, houseGID, m.battleUC, m.data, m.logger)
	m.syncers[key] = syncer
	syncer.start()

	m.logger.Infof("Started battle syncer for user %d, house %d", userID, houseGID)
}

// StopSync 停止战绩同步
func (m *BattleSyncManager) StopSync(userID int, houseGID int) {
	key := fmt.Sprintf("%d:%d", userID, houseGID)

	m.mu.Lock()
	defer m.mu.Unlock()

	if syncer, exists := m.syncers[key]; exists {
		syncer.stop()
		delete(m.syncers, key)
		m.logger.Infof("Stopped battle syncer for user %d, house %d", userID, houseGID)
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
	ctx          context.Context // 保存带 platform 的 context
	userID       int
	houseGID     int
	battleUC     *BattleRecordUseCase // 使用 UseCase 处理战绩
	data         *infra.Data          // 用于记录同步日志
	logger       *log.Helper
	stopChan     chan struct{}
	wg           sync.WaitGroup
	syncInterval time.Duration
	isFirstSync  bool
	sessionID    int32 // 会话 ID，用于记录同步日志
}

func newBattleSyncer(ctx context.Context, userID int, houseGID int, battleUC *BattleRecordUseCase, data *infra.Data, logger *log.Helper) *battleSyncer {
	return &battleSyncer{
		ctx:          ctx, // 保存 context
		userID:       userID,
		houseGID:     houseGID,
		battleUC:     battleUC,
		data:         data,
		logger:       logger,
		stopChan:     make(chan struct{}),
		syncInterval: 10 * time.Second, // 改为10秒一次
		isFirstSync:  true,
	}
}

func (s *battleSyncer) start() {
	s.logger.Infof("Starting battle syncer for house %d, isFirstSync=%v", s.houseGID, s.isFirstSync)
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

	s.logger.Infof("Fetching battle records for house %d (%s)...", s.houseGID, typeDesc)

	// 使用 BattleRecordUseCase.PullAndSave 统一处理战绩同步
	httpClient := &http.Client{Timeout: 10 * time.Second}
	baseURL := "http://phone2.foxuc.com/Ashx/GroService.ashx"

	saved, err := s.battleUC.PullAndSave(s.ctx, httpClient, baseURL, s.houseGID, 0, typeid)
	if err != nil {
		s.logger.Errorf("Failed to sync battle records for house %d: %v", s.houseGID, err)
		// 记录失败日志
		s.recordSyncLog(s.ctx, startTime, 0, model.SyncStatusFailed, err.Error())
		return
	}

	if saved > 0 {
		s.logger.Infof("Synced %d battle records for house %d", saved, s.houseGID)
	}

	// 记录成功日志
	s.recordSyncLog(s.ctx, startTime, int32(saved), model.SyncStatusSuccess, "")
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
