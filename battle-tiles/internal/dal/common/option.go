package common

type Options []*Option

type Option struct {
	Value int32  `json:"value"`
	Label string `json:"label"`
}
type OptionString struct {
	Value string `json:"value"`
	Label string `json:"label"`
}
