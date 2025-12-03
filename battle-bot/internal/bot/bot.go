package bot

import (
	"context"
	"crypto/md5"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"battle-bot/internal/config"
	"battle-bot/internal/plaza"
	"battle-bot/internal/plaza/game"
)

// Bot 游戏机器人
type Bot struct {
	cfg     *config.Config
	session *plaza.Session

	mu          sync.RWMutex
	isPlaying   bool
	currentGame *GameState
	gamesPlayed int

	stopChan chan struct{}
}

// GameState 当前游戏状态
type GameState struct {
	TableID   int
	PlayerPos int
	Cards     []string
	GamePhase string
}

// NewBot 创建新的机器人
func NewBot(cfg *config.Config) (*Bot, error) {
	bot := &Bot{
		cfg:      cfg,
		stopChan: make(chan struct{}),
	}

	return bot, nil
}

// Start 启动机器人
func (b *Bot) Start(ctx context.Context) error {
	// MD5密码
	pwdMD5 := strings.ToUpper(fmt.Sprintf("%x", md5.Sum([]byte(b.cfg.Account.Password))))

	// 如果game_user_id为0，先获取用户ID
	gameUserID := b.cfg.Game.GameUserID
	if gameUserID == 0 {
		log.Println("正在获取游戏用户ID...")

		var userInfo *game.UserLogonInfo
		var err error

		if b.cfg.Account.LoginMode == "mobile" {
			userInfo, err = plaza.GetUserInfoByMobileCtx(ctx, b.cfg.Plaza.Server82, b.cfg.Account.Username, pwdMD5)
		} else {
			userInfo, err = plaza.GetUserInfoByAccountCtx(ctx, b.cfg.Plaza.Server82, b.cfg.Account.Username, pwdMD5)
		}

		if err != nil {
			return fmt.Errorf("获取用户ID失败: %w", err)
		}

		gameUserID = int(userInfo.UserID)
		b.cfg.Game.GameUserID = gameUserID
		log.Printf("✅ 获取到游戏用户ID: %d", gameUserID)
	}

	// 创建Plaza会话配置
	sessionCfg := plaza.SessionConfig{
		Server82:      b.cfg.Plaza.Server82,
		Server87Host:  b.cfg.Plaza.Server87Host,
		KeepAlive:     time.Duration(b.cfg.Plaza.KeepAlive) * time.Second,
		AutoReconnect: b.cfg.Plaza.AutoReconnect,

		Identifier: b.cfg.Account.Username,
		UserPwdMD5: pwdMD5,
		UserID:     gameUserID,
		HouseGID:   b.cfg.Game.HouseGID,

		Handler: b, // Bot实现IPlazaHandler接口
	}

	// 连接并登录
	session, err := plaza.NewSessionWithConfig(sessionCfg)
	if err != nil {
		return fmt.Errorf("创建会话失败: %w", err)
	}

	b.session = session

	// 启动自动游戏协程
	go b.autoPlayLoop(ctx)

	return nil
}

// Stop 停止机器人
func (b *Bot) Stop(ctx context.Context) error {
	close(b.stopChan)

	if b.session != nil {
		b.session.Shutdown()
	}

	return nil
}

// autoPlayLoop 自动游戏循环
func (b *Bot) autoPlayLoop(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-b.stopChan:
			return
		case <-ticker.C:
			b.checkAndPlay()
		}
	}
}

// checkAndPlay 检查并执行游戏动作
func (b *Bot) checkAndPlay() {
	b.mu.RLock()
	defer b.mu.RUnlock()

	// 检查是否在活跃时间
	now := time.Now()
	hour := now.Hour()
	if hour < b.cfg.Bot.ActiveHoursStart || hour >= b.cfg.Bot.ActiveHoursEnd {
		return
	}

	// 检查今日游戏次数
	if b.gamesPlayed >= b.cfg.Bot.MaxGamesPerDay {
		return
	}

	// 自动加入桌台
	if b.cfg.Bot.AutoJoinTable && !b.isPlaying {
		// 不再打印，避免日志刷屏
		// log.Println("尝试查找并加入桌台...")
	}

	// 如果启用自动打牌且正在游戏中
	if b.cfg.Bot.AutoPlay && b.isPlaying {
		// TODO: 实现自动打牌逻辑
		log.Println("分析牌局并出牌...")
	}
}

// 实现IPlazaHandler接口的方法

func (b *Bot) OnSessionRestarted(session *plaza.Session) {
	log.Println("会话已重启")
	b.session = session
}

func (b *Bot) OnMemberListUpdated(members []*plaza.GroupMember) {
	log.Printf("成员列表更新: %d个成员", len(members))
}

func (b *Bot) OnMemberInserted(member *plaza.MemberInserted) {
	log.Printf("新成员加入: %+v", member)
}

func (b *Bot) OnMemberDeleted(member *plaza.MemberDeleted) {
	log.Printf("成员离开: %+v", member)
}

func (b *Bot) OnMemberRightUpdated(key string, memberID int, success bool) {
	log.Printf("成员权限更新: key=%s, memberID=%d, success=%v", key, memberID, success)
}

func (b *Bot) OnLoginDone(success bool) {
	if success {
		log.Println("✅ 登录成功！")
		log.Println("正在获取房间列表...")
		// Session会自动调用GetGroupMembers，触发服务器推送房间列表
	} else {
		log.Println("❌ 登录失败！")
	}
}

func (b *Bot) OnRoomListUpdated(tables []*plaza.TableInfo) {
	log.Printf("房间列表更新: %d个房间", len(tables))

	// 打印桌台列表
	log.Printf("📋 当前桌台列表：")
	for i, table := range tables {
		log.Printf("  [%d] MappedNum=%d, TableID=%d, KindID=%d",
			i+1, table.MappedNum, table.TableID, table.KindID)
	}

	// 尝试自动坐下（实验性）
	if b.cfg.Bot.AutoJoinTable && !b.isPlaying && len(tables) > 0 {
		// 选择第一个可用的桌台
		table := tables[0]
		log.Printf("🎯 尝试坐下到桌台: MappedNum=%d, TableID=%d, KindID=%d",
			table.MappedNum, table.TableID, table.KindID)

		go func() {
			time.Sleep(1 * time.Second) // 等待1秒后尝试
			b.tryJoinTable(table)
		}()
	}

	log.Printf("💡 机器人已就绪，等待坐下事件...")
	log.Printf("💡 如果自动坐下失败，请通过游戏客户端登录账号 %s 并手动坐下", b.cfg.Account.Username)
}

// tryJoinTable 尝试加入桌台
func (b *Bot) tryJoinTable(table *plaza.TableInfo) {
	if b.session == nil {
		return
	}

	b.mu.Lock()
	if b.isPlaying {
		b.mu.Unlock()
		return // 已经在游戏中
	}
	b.mu.Unlock()

	// 生成MD5密码
	pwdMD5 := strings.ToUpper(fmt.Sprintf("%x", md5.Sum([]byte(b.cfg.Account.Password))))

	// 简单的坐下命令
	chairID := 0
	cmd := plaza.CmdUserSitDown(table.MappedNum, chairID, pwdMD5)

	if err := b.session.SendCommand(cmd); err != nil {
		log.Printf("❌ 发送坐下命令失败: %v", err)
		return
	}

	log.Printf("📤 已发送坐下命令: MappedNum=%d (TableID=%d), Chair=%d", table.MappedNum, table.TableID, chairID)
	log.Printf("⏳ 等待服务器响应...")
}

func (b *Bot) OnUserSitDown(sitdown *plaza.UserSitDown) {
	log.Printf("玩家坐下: UserID=%d, MappedNum=%d, Chair=%d", sitdown.UserID, sitdown.MappedNum, sitdown.ChairID)

	// 如果是自己坐下，更新游戏状态
	if b.cfg.Game.GameUserID == int(sitdown.UserID) {
		b.mu.Lock()
		b.isPlaying = true
		b.currentGame = &GameState{
			TableID:   int(sitdown.MappedNum),
			PlayerPos: int(sitdown.ChairID),
			GamePhase: "waiting",
		}
		b.mu.Unlock()
		log.Printf("✅ 成功坐下！Table=%d, Chair=%d", sitdown.MappedNum, sitdown.ChairID)

		// 延迟发送准备命令，避免过快操作被服务器踢出
		go func() {
			time.Sleep(2 * time.Second)
			// 再次检查是否还在游戏中
			b.mu.Lock()
			playing := b.isPlaying
			b.mu.Unlock()

			if playing {
				readyCmd := plaza.CmdUserReady()
				if err := b.session.SendCommand(readyCmd); err != nil {
					log.Printf("❌ 发送准备命令失败: %v", err)
				} else {
					log.Printf("✅ 已发送准备命令")
				}
			}
		}()
	}
}

func (b *Bot) OnUserStandUp(standup *plaza.UserStandUp) {
	log.Printf("玩家站起: UserID=%d, MappedNum=%d, Chair=%d", standup.UserID, standup.MappedNum, standup.ChairID)

	// 如果是自己站起，更新游戏状态
	if b.cfg.Game.GameUserID == int(standup.UserID) {
		b.mu.Lock()
		b.isPlaying = false
		b.currentGame = nil
		b.mu.Unlock()
		log.Printf("已站起离开桌台")
	}
}

func (b *Bot) OnTableRenew(item *plaza.TableRenew) {
	log.Printf("桌台续约: %d -> %d", item.MappedNum, item.NewMappedNum)
}

func (b *Bot) OnDismissTable(table int) {
	log.Printf("桌台解散: %d", table)

	b.mu.Lock()
	if b.currentGame != nil && b.currentGame.TableID == table {
		b.isPlaying = false
		b.currentGame = nil
		b.gamesPlayed++
	}
	b.mu.Unlock()
}

func (b *Bot) OnAppliesForHouse(applyInfos []*plaza.ApplyInfo) {
	log.Printf("收到%d个申请", len(applyInfos))
}

func (b *Bot) OnReconnectFailed(houseGID int, retryCount int) {
	log.Printf("重连失败: HouseGID=%d, 重试次数=%d", houseGID, retryCount)
}

// 辅助方法

// randomDelay 随机延迟（模拟真实玩家）
func (b *Bot) randomDelay() {
	min := b.cfg.Bot.PlayDelayMin
	max := b.cfg.Bot.PlayDelayMax
	delay := min + rand.Intn(max-min)
	time.Sleep(time.Duration(delay) * time.Millisecond)
}
