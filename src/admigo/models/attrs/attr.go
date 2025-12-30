package attrs

type Attr struct {
	Al  string `json:"al"`
	Nm  string `json:"nm"`
	Tp  string `json:"tp,omitempty"`
	So  int    `json:"so"`
	Val string `json:"val"`
}
