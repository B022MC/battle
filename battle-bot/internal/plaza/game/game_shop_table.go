package game

type TableInfoVO struct {
	TableID   int `json:"table_id"`
	MappedNum int `json:"mapped_num"`
	GroupID   int `json:"group_id"`
	KindID    int `json:"kind_id"`
	BaseScore int `json:"base_score"`
}
