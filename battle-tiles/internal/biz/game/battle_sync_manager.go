package game

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	model "battle-tiles/internal/dal/model/game"
	repo "battle-tiles/internal/dal/repo/game"
	plazautils "battle-tiles/internal/utils/plaza"

	"github.com/go-kratos/kratos/v2/log"
)

// BattleSyncManager 管理所有会话的战绩同步
type BattleSyncManager struct {
	mu         sync.RWMutex
	syncers    map[string]*battleSyncer // key: "userID:houseGID"
	repo       repo.BattleRecordRepo
	logger     *log.Helper
}

// NewBattleSyncManager 创建战绩同步管理器
func NewBattleSyncManager(battleRepo repo.BattleRecordRepo, logger log.Logger) *BattleSyncManager {
	return &BattleSyncManager{
		syncers: make(map[string]*battleSyncer),
		repo:    battleRepo,
		logger:  log.NewHelper(logger),
	}
}

// StartSync 启动战绩同步
func (m *BattleSyncManager) StartSync(userID int, houseGID int) {
	key := fmt.Sprintf("%d:%d", userID, houseGID)
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// 如果已经存在，先停止
	if syncer, exists := m.syncers[key]; exists {
		syncer.stop()
		delete(m.syncers, key)
	}
	
	// 创建新的同步器
	syncer := newBattleSyncer(userID, houseGID, m.repo, m.logger)
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
	userID       int
	houseGID     int
	battleRepo   repo.BattleRecordRepo
	logger       *log.Helper
	stopChan     chan struct{}
	wg           sync.WaitGroup
	syncInterval time.Duration
	isFirstSync  bool
}

func newBattleSyncer(userID int, houseGID int, battleRepo repo.BattleRecordRepo, logger *log.Helper) *battleSyncer {
	return &battleSyncer{
		userID:       userID,
		houseGID:     houseGID,
		battleRepo:   battleRepo,
		logger:       logger,
		stopChan:     make(chan struct{}),
		syncInterval: 3 * time.Second,
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
	ctx := context.Background()

	// 第一次同步获取最近1小时的数据，之后获取最近3分钟的数据
	typeid := 0 // 最近3分钟
	typeDesc := "last 3 minutes"
	if s.isFirstSync {
		typeid = 2 // 最近1小时
		typeDesc = "last 1 hour"
		s.isFirstSync = false
	}

	s.logger.Infof("Fetching battle records for house %d (%s)...", s.houseGID, typeDesc)

	// 从 plaza 获取战绩数据
	// 使用现有的 HTTP API 函数
	httpClient := &http.Client{Timeout: 10 * time.Second}
	baseURL := "http://phone2.foxuc.com/Ashx/GroService.ashx"
	battles, err := plazautils.GetGroupNewBattleInfoCtx(ctx, httpClient, baseURL, s.houseGID, typeid)
	if err != nil {
		s.logger.Errorf("Failed to fetch battle info for house %d: %v", s.houseGID, err)
		return
	}

	s.logger.Infof("Fetched %d battle records from plaza for house %d", len(battles), s.houseGID)

	if len(battles) == 0 {
		return
	}
	
	// 转换为数据库模型
	records := make([]*model.GameBattleRecord, 0)
	for _, bi := range battles {
		// 验证数据完整性：零和游戏的总分应该为0
		totalScore := 0
		for _, p := range bi.Players {
			totalScore += p.Score
		}
		if totalScore != 0 {
			s.logger.Warnf("Invalid battle data: room %d, total score %d != 0", bi.RoomID, totalScore)
			continue
		}

		// 为每个玩家创建一条记录
		for _, player := range bi.Players {
			playerGameID := int32(player.UserGameID)
			record := &model.GameBattleRecord{
				HouseGID:      int32(s.houseGID),
				GroupID:       int32(s.houseGID), // 使用 houseGID 作为 groupID
				RoomUID:       int32(bi.RoomID),
				KindID:        int32(bi.KindID),
				BaseScore:     int32(bi.BaseScore),
				BattleAt:      time.Unix(int64(bi.CreateTime), 0),
				PlayersJSON:   "",             // TODO: 如果需要保存完整玩家列表，可以序列化为 JSON
				PlayerID:      nil,            // TODO: 需要从 game_user_id 映射到内部 player_id
				PlayerGameID:  &playerGameID,  // 玩家游戏 ID
				Score:         int32(player.Score),
				Fee:           0,   // TODO: 根据游戏规则计算手续费
				Factor:        1.0, // TODO: 根据游戏规则设置倍率
				PlayerBalance: 0,   // TODO: 需要查询玩家当前余额
			}
			records = append(records, record)
		}
	}
	
	// 批量保存到数据库（带去重）
	saved, err := s.battleRepo.SaveBatchWithDedup(ctx, records)
	if err != nil {
		s.logger.Errorf("Failed to save battle records for house %d: %v", s.houseGID, err)
		return
	}
	
	if saved > 0 {
		s.logger.Infof("Synced %d battle records for house %d", saved, s.houseGID)
	}
}

