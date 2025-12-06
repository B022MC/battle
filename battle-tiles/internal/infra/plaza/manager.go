// internal/infra/plaza/manager.go
package plaza

import (
	"battle-tiles/internal/conf"
	"battle-tiles/internal/consts"
	gamevo "battle-tiles/internal/dal/vo/game"
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	utilsplaza "battle-tiles/internal/utils/plaza"

	"github.com/go-kratos/kratos/v2/log"
)

// 让外部仍可传入 utilsplaza 的处理器
type Handler = utilsplaza.IPlazaHandler

type Config struct {
	Server82      string        // 例: "androidsc.foxuc.com:8200"
	Server87Host  string        // 例: "newbgp.foxuc.com"
	KeepAlive     time.Duration // 例: 30 * time.Second
	AutoReconnect bool
}

type Manager interface {
	// 启动/替换某个用户在某个 House 的会话
	StartUser(ctx context.Context, userID, houseGID int, mode consts.GameLoginMode, identifier, pwdMD5 string, gameUserID int, h Handler) error
	// 获取指定用户在指定 House 的会话
	Get(userID, houseGID int) (*utilsplaza.Session, bool)
	// 获取任意用户在该 House 下的会话（用于共享读取场景）
	GetAnyByHouse(houseGID int) (*utilsplaza.Session, bool)
	ListTables(userID, houseGID int) ([]*utilsplaza.TableInfo, error)
	// 列出某个用户的全部会话（不同 House）
	GetByUser(userID int) map[int]*utilsplaza.Session

	// 在线状态
	IsOnline(userID, houseGID int) bool
	WaitOnline(ctx context.Context, userID, houseGID int) error

	// 便捷封装（都要求明确 houseGID）
	ForbidMembers(userID, houseGID int, key string, members []int, forbid bool) error
	GetGroupMembers(userID, houseGID int) error
	DismissTable(userID, houseGID int, kindID, mappedNum int) error
	KickMember(userID, houseGID, memberID int) error

	// 关闭
	StopUser(userID, houseGID int) // 关闭该用户在该 House 的会话
	StopUserAll(userID int)        // 关闭该用户的所有会话
	StopAll()                      // 关闭全部会话（危险操作）

	// 只做登录探测，不建长连接、不入 sessions
	ProbeLogin(ctx context.Context, mode consts.GameLoginMode, identifier, pwdMD5 string) error
	// 探测并返回游戏端用户信息（用于持久化 game_user_id）
	ProbeLoginWithInfo(ctx context.Context, mode consts.GameLoginMode, identifier, pwdMD5 string) (*gamevo.UserLogonInfo, error)

	// 尝试通过协议捕获可见店铺列表（best-effort）
	ListHousesByLogin(ctx context.Context, mode consts.GameLoginMode, identifier, pwdMD5 string) ([]int, error)

	// 指标与健康
	Metrics() Metrics
	Health() HealthStatus

	// 设置重连失败回调
	SetReconnectFailedCallback(callback func(houseGID int, retryCount int) int32)
}

type manager struct {
	cfg        Config
	globalConf *conf.Global
	logger     *log.Helper

	mu       sync.RWMutex
	sessions map[string]*utilsplaza.Session // key: userKey(userID, houseGID)
	online   map[string]bool                // 在线状态

	// 指标
	restartCount  map[string]int       // 重启次数，按 userID:houseGID 统计
	lastRestartAt map[string]time.Time // 最近一次重启时间

	// 重连失败回调 (houseGID, retryCount) -> ctrlAccountID
	onReconnectFailedCallback func(houseGID int, retryCount int) int32
}

// 复合 key
func userKey(userID, houseGID int) string {
	return fmt.Sprintf("%d:%d", userID, houseGID)
}

func NewManager(globalConf *conf.Global, logger log.Logger) Manager {
	cfg := Config{
		Server82:      globalConf.Game.Plaza.Server82,
		Server87Host:  globalConf.Game.Plaza.Server87Host,
		KeepAlive:     time.Duration(globalConf.Game.Plaza.KeepaliveSeconds) * time.Second, // 配置是秒，这里转为 Duration
		AutoReconnect: globalConf.Game.Plaza.AutoReconnect,
	}
	return &manager{
		cfg:           cfg,
		globalConf:    globalConf,
		logger:        log.NewHelper(log.With(logger, "module", "infra/plaza/manager")),
		sessions:      make(map[string]*utilsplaza.Session),
		online:        make(map[string]bool),
		restartCount:  make(map[string]int),
		lastRestartAt: make(map[string]time.Time),
	}
}

// --- handler 包装器：实现并转发 IPlazaHandler 的全部方法，附带钩子 ---

type handlerWrapper struct {
	inner             Handler
	onLogin           func(ok bool)
	onRooms           func(tables []*utilsplaza.TableInfo)
	onRestart         func(s *utilsplaza.Session)
	onReconnectFailed func(houseGID int, retryCount int)
}

func (w *handlerWrapper) OnSessionRestarted(s *utilsplaza.Session) {
	if w.onRestart != nil {
		w.onRestart(s)
	}
	// 重启后，主动触发一次获取成员列表，促使服务端推送房间列表
	if s != nil {
		go s.GetGroupMembers()
	}
	if w.inner != nil {
		w.inner.OnSessionRestarted(s)
	}
}
func (w *handlerWrapper) OnMemberListUpdated(ms []*utilsplaza.GroupMember) {
	if w.inner != nil {
		w.inner.OnMemberListUpdated(ms)
	}
}
func (w *handlerWrapper) OnMemberInserted(m *utilsplaza.MemberInserted) {
	if w.inner != nil {
		w.inner.OnMemberInserted(m)
	}
}
func (w *handlerWrapper) OnMemberDeleted(m *utilsplaza.MemberDeleted) {
	if w.inner != nil {
		w.inner.OnMemberDeleted(m)
	}
}
func (w *handlerWrapper) OnMemberRightUpdated(account string, right int, allow bool) {
	if w.inner != nil {
		w.inner.OnMemberRightUpdated(account, right, allow)
	}
}
func (w *handlerWrapper) OnLoginDone(ok bool) {
	if w.onLogin != nil {
		w.onLogin(ok)
	}
	if w.inner != nil {
		w.inner.OnLoginDone(ok)
	}
}
func (w *handlerWrapper) OnRoomListUpdated(ts []*utilsplaza.TableInfo) {
	if w.onRooms != nil {
		w.onRooms(ts)
	}
	if w.inner != nil {
		w.inner.OnRoomListUpdated(ts)
	}
}
func (w *handlerWrapper) OnUserSitDown(e *utilsplaza.UserSitDown) {
	if w.inner != nil {
		w.inner.OnUserSitDown(e)
	}
}
func (w *handlerWrapper) OnUserStandUp(e *utilsplaza.UserStandUp) {
	if w.inner != nil {
		w.inner.OnUserStandUp(e)
	}
}
func (w *handlerWrapper) OnTableRenew(e *utilsplaza.TableRenew) {
	if w.inner != nil {
		w.inner.OnTableRenew(e)
	}
}
func (w *handlerWrapper) OnDismissTable(mappedNum int) {
	if w.inner != nil {
		w.inner.OnDismissTable(mappedNum)
	}
}
func (w *handlerWrapper) OnAppliesForHouse(list []*utilsplaza.ApplyInfo) {
	if w.inner != nil {
		w.inner.OnAppliesForHouse(list)
	}
}
func (w *handlerWrapper) OnReconnectFailed(houseGID int, retryCount int) {
	if w.onReconnectFailed != nil {
		w.onReconnectFailed(houseGID, retryCount)
	}
	if w.inner != nil {
		w.inner.OnReconnectFailed(houseGID, retryCount)
	}
}

// --- Manager 方法实现 ---

func (m *manager) StartUser(ctx context.Context, userID, houseGID int, mode consts.GameLoginMode, identifier, pwdMD5 string, gameUserID int, h Handler) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := userKey(userID, houseGID)

	// 若已存在该 House 的会话，先停掉老的
	if old, ok := m.sessions[key]; ok && old != nil {
		old.Shutdown()
		delete(m.sessions, key)
		delete(m.online, key)
	}

	// 包装 handler，用于写入在线状态/房间回调
	w := &handlerWrapper{
		inner: h,
		onLogin: func(ok bool) {
			m.mu.Lock()
			m.online[key] = ok
			m.mu.Unlock()
			// 登录成功后，主动拉取成员列表，促使服务端推送房间列表快照
			if ok {
				m.mu.RLock()
				s := m.sessions[key]
				m.mu.RUnlock()
				if s != nil {
					go s.GetGroupMembers()
				}
			}
		},
		onRooms: func(tables []*utilsplaza.TableInfo) {
			// manager 层只做转发/状态，不直接落库；如需落库，请在业务层 handler 处理
		},
		onRestart: func(s *utilsplaza.Session) {
			m.mu.Lock()
			defer m.mu.Unlock()
			// 更新会话指针与指标
			m.sessions[key] = s
			m.restartCount[key] = m.restartCount[key] + 1
			m.lastRestartAt[key] = time.Now()
		},
		onReconnectFailed: func(houseGID int, retryCount int) {
			m.logger.Errorf("重连失败 house=%d retries=%d, 将自动停用中控账号", houseGID, retryCount)
			// 调用回调获取 ctrlAccountID 并停用
			if m.onReconnectFailedCallback != nil {
				ctrlID := m.onReconnectFailedCallback(houseGID, retryCount)
				m.logger.Infof("自动停用中控账号 ctrlID=%d house=%d", ctrlID, houseGID)
			}
		},
	}

	cfg := utilsplaza.SessionConfig{
		Server82:      m.cfg.Server82,
		Server87Host:  m.cfg.Server87Host,
		KeepAlive:     m.cfg.KeepAlive,
		AutoReconnect: m.cfg.AutoReconnect,
		Identifier:    identifier, // 账号或手机号
		UserPwdMD5:    pwdMD5,
		UserID:        gameUserID,
		HouseGID:      houseGID,
		Handler:       w,
		LoginMode:     mode,
	}

	s, err := utilsplaza.NewSessionWithConfig(cfg)
	if err != nil {
		return err
	}
	m.sessions[key] = s

	// 在线状态由 OnLoginDone 回调置位
	return nil
}

// Metrics 指标快照
type Metrics struct {
	TotalSessions int
	OnlineCount   int
	RestartTotal  int
	RestartsByKey map[string]int
	LastRestartAt map[string]time.Time
}

// HealthStatus 健康检查结果
type HealthStatus struct {
	OK     bool
	Reason string
}

// Metrics 返回当前指标快照
func (m *manager) Metrics() Metrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	total := len(m.sessions)
	online := 0
	for k := range m.sessions {
		if m.online[k] {
			online++
		}
	}

	// 拷贝 map 避免逃逸引用
	restarts := make(map[string]int, len(m.restartCount))
	for k, v := range m.restartCount {
		restarts[k] = v
	}
	last := make(map[string]time.Time, len(m.lastRestartAt))
	for k, v := range m.lastRestartAt {
		last[k] = v
	}

	sum := 0
	for _, v := range restarts {
		sum += v
	}

	return Metrics{
		TotalSessions: total,
		OnlineCount:   online,
		RestartTotal:  sum,
		RestartsByKey: restarts,
		LastRestartAt: last,
	}
}

// Health 依据在线率与重启次数给出粗略健康状态
func (m *manager) Health() HealthStatus {
	mt := m.Metrics()
	if mt.TotalSessions == 0 {
		return HealthStatus{OK: true, Reason: "no sessions"}
	}
	offline := mt.TotalSessions - mt.OnlineCount
	if float64(mt.OnlineCount)/float64(mt.TotalSessions) < 0.5 {
		return HealthStatus{OK: false, Reason: fmt.Sprintf("online %d/%d, too many offline: %d", mt.OnlineCount, mt.TotalSessions, offline)}
	}
	if mt.RestartTotal > mt.TotalSessions*5 { // 任意粗略阈值
		return HealthStatus{OK: true, Reason: fmt.Sprintf("high restarts total=%d", mt.RestartTotal)}
	}
	return HealthStatus{OK: true, Reason: "ok"}
}

func (m *manager) Get(userID, houseGID int) (*utilsplaza.Session, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.sessions[userKey(userID, houseGID)]
	return s, ok
}

func (m *manager) GetAnyByHouse(houseGID int) (*utilsplaza.Session, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	suffix := fmt.Sprintf(":%d", houseGID)
	for k, s := range m.sessions {
		if strings.HasSuffix(k, suffix) && s != nil {
			return s, true
		}
	}
	return nil, false
}

func (m *manager) ListTables(userID, houseGID int) ([]*utilsplaza.TableInfo, error) {
	s, ok := m.Get(userID, houseGID)
	if !ok || s == nil {
		return nil, fmt.Errorf("session not found for user %d in house %d", userID, houseGID)
	}
	return s.ListTables(), nil
}

func (m *manager) GetByUser(userID int) map[int]*utilsplaza.Session {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make(map[int]*utilsplaza.Session)
	prefix := fmt.Sprintf("%d:", userID)
	for k, s := range m.sessions {
		// k = "userID:houseGID"
		if len(k) > len(prefix) && k[:len(prefix)] == prefix {
			// 解析 houseGID
			var uid, hg int
			if _, err := fmt.Sscanf(k, "%d:%d", &uid, &hg); err == nil && uid == userID {
				out[hg] = s
			}
		}
	}
	return out
}

// 在线状态
func (m *manager) IsOnline(userID, houseGID int) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.online[userKey(userID, houseGID)]
}

func (m *manager) WaitOnline(ctx context.Context, userID, houseGID int) error {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("wait online timeout: user=%d house=%d", userID, houseGID)
		case <-ticker.C:
			if m.IsOnline(userID, houseGID) {
				return nil
			}
		}
	}
}

// ====== 封装操作（要求 houseGID）======

func (m *manager) ForbidMembers(userID, houseGID int, key string, members []int, forbid bool) error {
	s, ok := m.Get(userID, houseGID)
	if !ok {
		return fmt.Errorf("session not found for user %d in house %d", userID, houseGID)
	}
	// 推送到游戏服务器并更新当前 Session 的缓存
	s.ForbidMembers(key, members, forbid)

	// 同步禁用状态到该 house 的所有其他 Session
	m.mu.RLock()
	suffix := fmt.Sprintf(":%d", houseGID)
	for k, sess := range m.sessions {
		if strings.HasSuffix(k, suffix) && sess != nil && sess != s {
			// 只更新缓存，不重复发送命令
			for _, memberID := range members {
				sess.SetForbiddenStatus(memberID, forbid)
			}
		}
	}
	m.mu.RUnlock()

	return nil
}

func (m *manager) GetGroupMembers(userID, houseGID int) error {
	s, ok := m.Get(userID, houseGID)
	if !ok {
		return fmt.Errorf("session not found for user %d in house %d", userID, houseGID)
	}
	s.GetGroupMembers()
	return nil
}

func (m *manager) DismissTable(userID, houseGID int, kindID, mappedNum int) error {
	s, ok := m.Get(userID, houseGID)
	if !ok {
		return fmt.Errorf("session not found for user %d in house %d", userID, houseGID)
	}
	s.DismissTable(kindID, mappedNum)
	return nil
}

func (m *manager) KickMember(userID, houseGID, memberID int) error {
	s, ok := m.Get(userID, houseGID)
	if !ok || s == nil {
		return fmt.Errorf("session not found for user %d in house %d", userID, houseGID)
	}
	s.KickOffMember(houseGID, memberID)
	return nil
}

func (m *manager) StopUser(userID, houseGID int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := userKey(userID, houseGID)
	if s, ok := m.sessions[key]; ok && s != nil {
		s.Shutdown()
		delete(m.sessions, key)
	}
	delete(m.online, key) // 清理在线标志
}

func (m *manager) StopUserAll(userID int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	prefix := fmt.Sprintf("%d:", userID)
	for k, s := range m.sessions {
		if len(k) > len(prefix) && k[:len(prefix)] == prefix {
			if s != nil {
				s.Shutdown()
			}
			delete(m.sessions, k)
		}
	}
	// 清理在线表
	for k := range m.online {
		if strings.HasPrefix(k, prefix) {
			delete(m.online, k)
		}
	}
}

func (m *manager) StopAll() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, s := range m.sessions {
		if s != nil {
			s.Shutdown()
		}
		delete(m.sessions, k)
	}
	m.online = make(map[string]bool)
}

// 只做登录探测，不建连接
func (m *manager) ProbeLogin(ctx context.Context, mode consts.GameLoginMode, identifier, pwdMD5 string) error {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}
	id := strings.TrimSpace(identifier)
	pass := strings.ToUpper(strings.TrimSpace(pwdMD5))

	switch mode {
	case consts.GameLoginModeAccount:
		if _, err := utilsplaza.GetUserInfoByAccountCtx(ctx, m.cfg.Server82, id, pass); err != nil {
			return fmt.Errorf("login verify failed on 82 (account): %w", err)
		}
		return nil
	case consts.GameLoginModeMobile:
		if _, err := utilsplaza.GetUserInfoByMobileCtx(ctx, m.cfg.Server82, id, pass); err != nil {
			return fmt.Errorf("login verify failed on 82 (mobile): %w", err)
		}
		return nil
	default:
		return fmt.Errorf("invalid login mode: %v", mode)
	}
}

func (m *manager) ProbeLoginWithInfo(ctx context.Context, mode consts.GameLoginMode, identifier, pwdMD5 string) (*gamevo.UserLogonInfo, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}
	id := strings.TrimSpace(identifier)
	pass := strings.ToUpper(strings.TrimSpace(pwdMD5))
	switch mode {
	case consts.GameLoginModeAccount:
		return utilsplaza.GetUserInfoByAccountCtx(ctx, m.cfg.Server82, id, pass)
	case consts.GameLoginModeMobile:
		return utilsplaza.GetUserInfoByMobileCtx(ctx, m.cfg.Server82, id, pass)
	default:
		return nil, fmt.Errorf("invalid login mode: %v", mode)
	}
}

func (m *manager) ListHousesByLogin(ctx context.Context, mode consts.GameLoginMode, identifier, pwdMD5 string) ([]int, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 8*time.Second)
		defer cancel()
	}
	id := strings.TrimSpace(identifier)
	pass := strings.ToUpper(strings.TrimSpace(pwdMD5))

	// 先探测获取游戏端 UserID
	info, err := m.ProbeLoginWithInfo(ctx, mode, id, pass)
	if err != nil {
		return nil, err
	}

	cfg := utilsplaza.SessionConfig{
		Server82:      m.cfg.Server82,
		Server87Host:  m.cfg.Server87Host,
		KeepAlive:     m.cfg.KeepAlive,
		AutoReconnect: false,
		Identifier:    id,
		UserPwdMD5:    pass,
		UserID:        int(info.UserID),
		HouseGID:      0,
		Handler:       nil,
		LoginMode:     mode,
	}
	s, err := utilsplaza.NewSessionWithConfig(cfg)
	if err != nil {
		return nil, err
	}
	defer s.Shutdown()

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			hs := s.ListHouses()
			if len(hs) == 0 {
				return nil, fmt.Errorf("no houses received")
			}
			return hs, nil
		case <-ticker.C:
			if arr := s.ListHouses(); len(arr) > 0 {
				return arr, nil
			}
		}
	}
}

// SetReconnectFailedCallback 设置重连失败回调
func (m *manager) SetReconnectFailedCallback(callback func(houseGID int, retryCount int) int32) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onReconnectFailedCallback = callback
}
