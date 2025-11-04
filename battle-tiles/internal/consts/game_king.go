package consts

const (
	GameKindDingErHong  = 60
	GameKindHongErShi   = 110
	GameKindDuanGouQia3 = 57
	GameKindPaoDeKuai2  = 115
	GameKindDouShiSi    = 95
	GameKindDuanGouQia2 = 114
	GameKindHongZhong2  = 150
)

func GetKindName(id int) string {
	switch id {
	case GameKindDingErHong:
		return "丁二红"
	case GameKindHongErShi:
		return "红二十"
	case GameKindDuanGouQia3:
		return "三人断勾卡"
	case GameKindPaoDeKuai2:
		return "跑得快"
	case GameKindDouShiSi:
		return "斗十四"
	case GameKindDuanGouQia2:
		return "断勾卡"
	case GameKindHongZhong2:
		return "红中"
	default:
		return "unknown"
	}
}

func GetKindID(name string) int {
	switch name {
	case "丁二红":
		return GameKindDingErHong
	case "红二十":
		return GameKindHongErShi
	case "三人断勾卡":
		return GameKindDuanGouQia3
	case "跑得快":
		return GameKindPaoDeKuai2
	case "斗十四":
		return GameKindDouShiSi
	case "断勾卡":
		return GameKindDuanGouQia2
	case "红中":
		return GameKindHongZhong2
	default:
		return 0
	}
}
