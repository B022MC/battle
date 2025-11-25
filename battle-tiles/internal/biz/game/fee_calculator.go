package game

import (
	"encoding/json"
	"fmt"
	"sort"

	gameVO "battle-tiles/internal/dal/vo/game"
)

// FeeRule 运费规则
type FeeRule struct {
	Threshold int `json:"threshold"` // 阈值（分数）
	Fee       int `json:"fee"`       // 费用（单位：分）
	Kind      int `json:"kind"`      // 玩法类型ID（可选，0表示所有玩法）
	Base      int `json:"base"`      // 底分（可选，0表示所有底分）
}

// FeesConfig 运费配置
type FeesConfig struct {
	Rules []FeeRule `json:"rules"`
}

// ParseFeesJSON 解析运费配置JSON
func ParseFeesJSON(feesJSON string) (*FeesConfig, error) {
	if feesJSON == "" {
		return &FeesConfig{Rules: []FeeRule{}}, nil
	}

	var config FeesConfig
	if err := json.Unmarshal([]byte(feesJSON), &config); err != nil {
		return nil, fmt.Errorf("解析运费配置失败: %w", err)
	}

	return &config, nil
}

// CalculateFee 计算运费
// 根据配置规则和战绩信息计算应收取的运费
// 匹配优先级：
//  1. 先查找全局规则 (kind=0, base=0)
//  2. 如果没有，再查找特定游戏类型和底分的规则
func CalculateFee(feesJSON string, battle *gameVO.BattleInfo) int32 {
	config, err := ParseFeesJSON(feesJSON)
	if err != nil || len(config.Rules) == 0 {
		return 0 // 无规则或解析失败，不收费
	}

	// 计算最高分
	maxScore := 0
	for _, p := range battle.Players {
		if p.Score > maxScore {
			maxScore = p.Score
		}
	}

	// 匹配逻辑与 passing-dragonfly 对齐：
	// 1. 优先查找全局通用规则 (kind=0 && base=0)
	for _, rule := range config.Rules {
		if rule.Kind == 0 && rule.Base == 0 {
			// 全局规则：检查分数是否达到阈值
			if maxScore >= rule.Threshold {
				return int32(rule.Fee)
			}
		}
	}

	// 2. 如果没有全局规则，查找特定游戏类型和底分的规则
	for _, rule := range config.Rules {
		// 跳过全局规则（已经处理过）
		if rule.Kind == 0 && rule.Base == 0 {
			continue
		}

		// 匹配游戏类型
		kindMatches := (rule.Kind == 0 || rule.Kind == battle.KindID)
		// 匹配底分
		baseMatches := (rule.Base == 0 || rule.Base == battle.BaseScore)

		// 同时匹配游戏类型和底分，且分数达到阈值
		if kindMatches && baseMatches && maxScore >= rule.Threshold {
			return int32(rule.Fee)
		}
	}

	return 0 // 未匹配到任何规则
}

// Winner 赢家信息
type Winner struct {
	UserGameID int
	Score      int
	GroupID    int32
}

// FindWinners 查找赢家
// 返回所有最高分的玩家（可能有多个平分的情况）
// 优化：一次遍历完成，避免重复扫描
func FindWinners(players []*gameVO.BattleSettle, playerGroups map[int]int32) []Winner {
	if len(players) == 0 {
		return nil
	}

	// 一次遍历找出最高分和所有赢家
	maxScore := players[0].Score
	winners := make([]Winner, 0, 4) // 预分配容量，通常不会超过4个赢家

	for _, p := range players {
		if p.Score > maxScore {
			// 发现更高分，清空之前的赢家，重新开始
			maxScore = p.Score
			winners = winners[:0]
		}

		if p.Score == maxScore {
			groupID := int32(0)
			if gid, ok := playerGroups[p.UserGameID]; ok {
				groupID = gid
			}
			winners = append(winners, Winner{
				UserGameID: p.UserGameID,
				Score:      p.Score,
				GroupID:    groupID,
			})
		}
	}

	return winners
}

// GroupInfo 圈子信息
type GroupInfo struct {
	GroupID   int32
	PlayerIDs []int
	IsWinner  bool
	PlayerFee int32 // 每个玩家应付的费用
	TotalFee  int32 // 圈子总费用
}

// CalculateFeeDistribution 计算费用分配
// 返回每个圈子的费用信息
// 参数：
//   - players: 参与战绩的所有玩家
//   - playerGroups: 玩家ID到圈子ID的映射
//   - totalFee: 本局总运费
//   - shareFee: 是否分运费模式
//
// 返回：圈子ID到费用信息的映射
func CalculateFeeDistribution(
	players []*gameVO.BattleSettle,
	playerGroups map[int]int32,
	totalFee int32,
	shareFee bool,
) map[int32]*GroupInfo {
	// 1. 找出赢家
	winners := FindWinners(players, playerGroups)
	if len(winners) == 0 {
		return nil
	}

	// 2. 统计各圈子的玩家
	groups := make(map[int32]*GroupInfo)
	winnerGroups := make(map[int32]bool)

	for _, p := range players {
		groupID := int32(0)
		if gid, ok := playerGroups[p.UserGameID]; ok {
			groupID = gid
		}

		if groups[groupID] == nil {
			groups[groupID] = &GroupInfo{
				GroupID:   groupID,
				PlayerIDs: []int{},
			}
		}
		groups[groupID].PlayerIDs = append(groups[groupID].PlayerIDs, p.UserGameID)
	}

	// 标记赢家圈子
	for _, w := range winners {
		winnerGroups[w.GroupID] = true
		if groups[w.GroupID] != nil {
			groups[w.GroupID].IsWinner = true
		}
	}

	// 3. 计算费用分配
	if shareFee {
		// 分运费模式：与 passing-dragonfly 对齐
		// 逻辑：总运费由所有圈子平分，但赢家圈需要收到差价补偿
		numGroups := len(groups)
		numWinners := len(winners)
		if numGroups == 0 || numWinners == 0 {
			return groups
		}

		// 零除保护
		sharedFee := int32(0)
		winnerFee := int32(0)
		if numGroups > 0 {
			sharedFee = totalFee / int32(numGroups) // 每个圈子应分摊的金额
		}
		if numWinners > 0 {
			winnerFee = totalFee / int32(numWinners) // 赢家本应支付的金额
		}
		feePayoff := winnerFee - sharedFee // 赢家圈应收到的差价补偿

		for _, g := range groups {
			if g.IsWinner {
				// 赢家圈：支付差价（实际支付少于应付）
				g.TotalFee = feePayoff
			} else {
				// 其他圈：支付平摊金额
				g.TotalFee = sharedFee
			}
			// 圈子内平均分摊到每个玩家
			if len(g.PlayerIDs) > 0 {
				g.PlayerFee = g.TotalFee / int32(len(g.PlayerIDs))
			}
		}
	} else {
		// 不分运：赢家承担全部
		numWinners := len(winners)
		if numWinners == 0 {
			return groups
		}

		// 零除保护
		winnerFee := int32(0)
		if numWinners > 0 {
			winnerFee = totalFee / int32(numWinners)
		}

		for _, w := range winners {
			if g := groups[w.GroupID]; g != nil {
				g.TotalFee += winnerFee
				// 赢家圈子内平均分摊
				if len(g.PlayerIDs) > 0 {
					g.PlayerFee = g.TotalFee / int32(len(g.PlayerIDs))
				}
			}
		}
	}

	return groups
}

// FeeSettlement 费用结算信息
type FeeSettlement struct {
	GroupID  int32
	Amount   int32 // 结算金额（正数=支出，负数=收入）
	IsPayoff bool  // 是否是结转金额
}

// CalculateFeeSettlements 计算费用结算（用于分运费模式）
func CalculateFeeSettlements(
	players []*gameVO.BattleSettle,
	playerGroups map[int]int32,
	totalFee int32,
	shareFee bool,
) []FeeSettlement {
	settlements := make([]FeeSettlement, 0)

	if !shareFee {
		// 不分运模式：只记录赢家支付的费用，无需结转
		return settlements
	}

	// 分运费模式
	winners := FindWinners(players, playerGroups)
	if len(winners) == 0 {
		return settlements
	}

	// 统计圈子
	groups := make(map[int32]bool)
	winnerGroups := make(map[int32]bool)
	for _, p := range players {
		groupID := int32(0)
		if gid, ok := playerGroups[p.UserGameID]; ok {
			groupID = gid
		}
		groups[groupID] = true
	}
	for _, w := range winners {
		winnerGroups[w.GroupID] = true
	}

	// 计算分摊和结转
	numGroups := len(groups)
	numWinners := len(winners)
	if numGroups == 0 || numWinners == 0 {
		return settlements
	}

	// 零除保护
	sharedFee := int32(0)
	winnerFee := int32(0)
	if numGroups > 0 {
		sharedFee = totalFee / int32(numGroups) // 每个圈子应分摊
	}
	if numWinners > 0 {
		winnerFee = totalFee / int32(numWinners) // 赢家本应支付
	}
	feePayoff := winnerFee - sharedFee // 应退给赢家圈子的金额

	// 为每个圈子计算结转
	for groupID := range groups {
		if winnerGroups[groupID] {
			// 赢家圈子：收到补偿（负数表示收入）
			if feePayoff > 0 {
				settlements = append(settlements, FeeSettlement{
					GroupID:  groupID,
					Amount:   -feePayoff,
					IsPayoff: true,
				})
			}
		} else {
			// 其他圈子：支付补偿给赢家
			if sharedFee > 0 {
				settlements = append(settlements, FeeSettlement{
					GroupID:  groupID,
					Amount:   sharedFee,
					IsPayoff: true,
				})
			}
		}
	}

	// 按圈子ID排序，保持一致性
	sort.Slice(settlements, func(i, j int) bool {
		return settlements[i].GroupID < settlements[j].GroupID
	})

	return settlements
}
