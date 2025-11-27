package service

import (
	biz "battle-tiles/internal/biz/game"
	"battle-tiles/internal/conf"
	cloudRepo "battle-tiles/internal/dal/repo/cloud"
	gameRepo "battle-tiles/internal/dal/repo/game"
	"battle-tiles/internal/infra"
	"battle-tiles/internal/infra/plaza"
	"context"
	"time"

	pdb "battle-tiles/pkg/plugin/dbx"

	"github.com/go-kratos/kratos/v2/log"
)

// SessionMonitor 定时同步 plaza 会话到 DB（按店铺维度，一店一条）
type SessionMonitor struct {
	log            *log.Helper
	mgr            plaza.Manager
	cloud          cloudRepo.BasePlatformRepo
	link           gameRepo.GameCtrlAccountHouseRepo
	sess           gameRepo.SessionRepo
	ctrlRepo       gameRepo.GameCtrlAccountRepo // 新增：用于检查中控账号状态
	ctrl           *biz.CtrlSessionUseCase
	syncMgr        *biz.BattleSyncManager // 新增：战绩同步管理器
	tick           time.Duration
	stopC          chan struct{}
	verboseTaskLog bool // 是否显示详细的任务日志
}

func NewSessionMonitor(logConf *conf.Log, logger log.Logger, mgr plaza.Manager, cloud cloudRepo.BasePlatformRepo, link gameRepo.GameCtrlAccountHouseRepo, sess gameRepo.SessionRepo, ctrlRepo gameRepo.GameCtrlAccountRepo, syncMgr *biz.BattleSyncManager, uc *biz.CtrlSessionUseCase) *SessionMonitor {
	verboseTaskLog := false
	if logConf != nil && logConf.VerboseTaskLog {
		verboseTaskLog = true
	}
	m := &SessionMonitor{
		log:            log.NewHelper(log.With(logger, "module", "service/session_monitor")),
		mgr:            mgr,
		cloud:          cloud,
		link:           link,
		sess:           sess,
		ctrlRepo:       ctrlRepo,
		syncMgr:        syncMgr,
		tick:           10 * time.Second,
		stopC:          make(chan struct{}),
		verboseTaskLog: verboseTaskLog,
	}
	if uc != nil {
		m.WithCtrlUseCase(uc)
	}
	return m
}

func (m *SessionMonitor) Start(ctx context.Context) {
	go func() {
		t := time.NewTicker(m.tick)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-m.stopC:
				return
			case <-t.C:
				m.syncOnce(ctx)
			}
		}
	}()
}

func (m *SessionMonitor) Stop() { close(m.stopC) }

// WithCtrlUseCase 注入会话启动用例（可选）
func (m *SessionMonitor) WithCtrlUseCase(uc *biz.CtrlSessionUseCase) *SessionMonitor {
	m.ctrl = uc
	return m
}

func (m *SessionMonitor) syncOnce(ctx context.Context) {
	platforms, err := m.cloud.ListPlatform()
	if err != nil {
		m.log.Errorf("list platforms err: %v", err)
		return
	}
	for _, p := range platforms {
		if p == nil || p.Platform == "" {
			continue
		}
		// 为当前平台设置数据库上下文
		pctx := context.WithValue(ctx, pdb.CtxDBKey, p.Platform)
		// 如果不显示详细日志，则使用静默模式
		if !m.verboseTaskLog {
			pctx = infra.WithQuietDB(pctx)
		}
		houses, err := m.link.ListDistinctHouses(pctx)
		if err != nil {
			m.log.Errorf("list distinct houses err: %v", err)
			continue
		}
		for _, hg := range houses {
			// 取绑定的中控 ctrlID（优先用于自动拉起）
			ctrls, _ := m.link.ListByHouse(pctx, hg)
			var ctrlID int32 = 0
			if len(ctrls) > 0 && ctrls[0] != nil {
				ctrlID = ctrls[0].Id
			}

			// 检查中控账号状态，如果停用则停止会话
			if ctrlID > 0 {
				ctrl, err := m.ctrlRepo.Get(pctx, ctrlID)
				if err == nil && ctrl.Status != 1 {
					// 中控账号已停用，停止会话和战绩同步
					if _, ok := m.mgr.GetAnyByHouse(int(hg)); ok {
						if m.verboseTaskLog {
							m.log.Infof("ctrl account %d is disabled, stopping session and sync for house %d", ctrlID, hg)
						}
						m.mgr.StopUser(1, int(hg)) // 停止超管的会话
						// 停止战绩同步
						if m.syncMgr != nil {
							m.syncMgr.StopSync(1, int(hg))
						}
						if err := m.sess.SetOfflineByHouse(pctx, hg); err != nil {
							m.log.Errorf("set offline house=%d err=%v", hg, err)
						}
					}
					continue
				}
			}

			if _, ok := m.mgr.GetAnyByHouse(int(hg)); ok {
				// 已有会话（任意用户）→ 什么都不做,会话启动时已经更新了数据库
				// 不要在这里频繁调用 EnsureOnlineByHouse,会导致每10秒都更新数据库!
			} else {
				// 无任何会话 → 若有绑定中控且注入了 UseCase，尝试自动拉起（以超管 1 身份）
				started := false
				if m.ctrl != nil && ctrlID > 0 {
					// 再次检查中控账号状态
					ctrl, err := m.ctrlRepo.Get(pctx, ctrlID)
					if err == nil && ctrl.Status == 1 {
						if err := m.ctrl.StartSession(pctx, 1, ctrlID, hg); err != nil {
							if m.verboseTaskLog {
								m.log.Errorf("auto start session failed house=%d ctrl=%d err=%v", hg, ctrlID, err)
							}
						} else {
							started = true
							if m.verboseTaskLog {
								m.log.Infof("auto start session ok house=%d ctrl=%d", hg, ctrlID)
							}
						}
					}
				}
				// 若本轮已尝试并成功发起启动，则不在本 tick 立即写离线，下一轮再校正
				if !started {
					if _, ok2 := m.mgr.GetAnyByHouse(int(hg)); !ok2 {
						if err := m.sess.SetOfflineByHouse(pctx, hg); err != nil {
							m.log.Errorf("set offline house=%d err=%v", hg, err)
						}
					}
				}
			}
		}
	}
}
