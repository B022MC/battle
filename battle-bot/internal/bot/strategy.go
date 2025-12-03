package bot

import (
	"log"
	"math/rand"
)

// Strategy 打牌策略接口
type Strategy interface {
	// DecidePlay 决定出牌
	// cards: 当前手牌
	// gameState: 当前游戏状态
	// 返回: 要出的牌
	DecidePlay(cards []string, gameState *GameState) []string
}

// RandomStrategy 随机策略
type RandomStrategy struct{}

func (s *RandomStrategy) DecidePlay(cards []string, gameState *GameState) []string {
	if len(cards) == 0 {
		return nil
	}

	// 随机选择1-3张牌出
	count := 1 + rand.Intn(3)
	if count > len(cards) {
		count = len(cards)
	}

	// 随机打乱
	shuffled := make([]string, len(cards))
	copy(shuffled, cards)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	result := shuffled[:count]
	log.Printf("随机策略: 从%d张牌中选择%d张", len(cards), count)
	return result
}

// ConservativeStrategy 保守策略
type ConservativeStrategy struct{}

func (s *ConservativeStrategy) DecidePlay(cards []string, gameState *GameState) []string {
	if len(cards) == 0 {
		return nil
	}

	// 优先出小牌
	// TODO: 实现牌型排序和分析
	log.Println("保守策略: 优先出小牌")
	return []string{cards[0]}
}

// AggressiveStrategy 激进策略
type AggressiveStrategy struct{}

func (s *AggressiveStrategy) DecidePlay(cards []string, gameState *GameState) []string {
	if len(cards) == 0 {
		return nil
	}

	// 优先出大牌
	// TODO: 实现牌型排序和分析
	log.Println("激进策略: 优先出大牌")
	return []string{cards[len(cards)-1]}
}

// NewStrategy 根据配置创建策略
func NewStrategy(strategyName string) Strategy {
	switch strategyName {
	case "random":
		return &RandomStrategy{}
	case "conservative":
		return &ConservativeStrategy{}
	case "aggressive":
		return &AggressiveStrategy{}
	default:
		log.Printf("未知策略: %s, 使用默认随机策略", strategyName)
		return &RandomStrategy{}
	}
}

// CardAnalyzer 牌型分析器
type CardAnalyzer struct{}

// AnalyzeCards 分析手牌
func (a *CardAnalyzer) AnalyzeCards(cards []string) *CardAnalysis {
	// TODO: 实现牌型分析
	// - 单张、对子、三张、炸弹
	// - 顺子识别
	// - 牌力评估
	return &CardAnalysis{
		TotalCount: len(cards),
	}
}

// CardAnalysis 牌型分析结果
type CardAnalysis struct {
	TotalCount int
	Singles    []string   // 单张
	Pairs      [][]string // 对子
	Triples    [][]string // 三张
	Bombs      [][]string // 炸弹
	Sequences  [][]string // 顺子
	Power      int        // 牌力评估
}
