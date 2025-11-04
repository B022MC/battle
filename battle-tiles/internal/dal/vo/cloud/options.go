package cloud

type Options []Option
type Option struct {
	Platform string `json:"cloud"`
	Label    string `json:"label"`
}
