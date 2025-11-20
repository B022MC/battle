package plaza

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"testing"
)

// TestAccountLoginAndVerifyGameID 测试账号登录并验证GameID
func TestAccountLoginAndVerifyGameID(t *testing.T) {
	// 测试账号信息
	account := "1106162940"
	password := "1475369"

	// MD5加密密码（大写）
	hash := md5.Sum([]byte(password))
	pwdMD5 := hex.EncodeToString(hash[:])
	pwdMD5Upper := fmt.Sprintf("%X", hash[:])

	fmt.Printf("\n=== 账号登录测试 ===\n")
	fmt.Printf("账号: %s\n", account)
	fmt.Printf("密码: %s\n", password)
	fmt.Printf("MD5(小写): %s\n", pwdMD5)
	fmt.Printf("MD5(大写): %s\n", pwdMD5Upper)
	fmt.Printf("\n")

	// 游戏服务器地址（从配置文件获取）
	server82 := "androidsc.foxuc.com:8200"

	ctx := context.Background()

	// 调用登录接口
	userInfo, err := GetUserInfoByAccountCtx(ctx, server82, account, pwdMD5Upper)
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}

	fmt.Printf("=== 登录成功 ===\n")
	fmt.Printf("UserID: %d\n", userInfo.UserID)
	fmt.Printf("GameID: %d\n", userInfo.GameID)
	fmt.Printf("\n")

	// 验证关系
	fmt.Printf("=== 字段关系验证 ===\n")
	fmt.Printf("account字段值 (输入): %s\n", account)
	fmt.Printf("GameID (登录返回): %d\n", userInfo.GameID)
	fmt.Printf("GameID字符串形式: %d\n", userInfo.GameID)
	fmt.Printf("\n")

	// 关键确认
	accountAsInt := account
	gameIDStr := fmt.Sprintf("%d", userInfo.GameID)

	fmt.Printf("=== 结论 ===\n")
	fmt.Printf("game_account.account 应该存储: %s\n", account)
	fmt.Printf("game_account.game_player_id 应该存储: %s\n", gameIDStr)

	if accountAsInt == gameIDStr {
		fmt.Printf("✅ account == GameID (都是 %s)\n", account)
		fmt.Printf("说明: Plaza API返回的GameID就是账号本身\n")
	} else {
		fmt.Printf("❌ account (%s) != GameID (%s)\n", account, gameIDStr)
		fmt.Printf("说明: Plaza API返回的GameID是另一个值\n")
	}
	fmt.Printf("\n")

	fmt.Printf("=== 战绩同步时的映射关系 ===\n")
	fmt.Printf("Plaza API战绩中的 player.UserGameID = %d\n", userInfo.GameID)
	fmt.Printf("应该通过什么字段查询 game_account 表?\n")
	if accountAsInt == gameIDStr {
		fmt.Printf("  → 通过 game_account.account = '%s' 查询\n", account)
		fmt.Printf("  → 或通过 game_account.game_player_id = '%s' 查询\n", gameIDStr)
		fmt.Printf("  → (两个字段值相同)\n")
	} else {
		fmt.Printf("  → 通过 game_account.game_player_id = '%s' 查询\n", gameIDStr)
		fmt.Printf("  → (account字段存储的是登录账号，game_player_id存储的是游戏玩家ID)\n")
	}
	fmt.Printf("\n")
}
