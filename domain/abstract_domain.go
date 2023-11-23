package domain

type GetListParameter struct {
	Page   int64  `json:"page"`
	Limit  int64  `json:"limit"`
	Filter string `json:"filter"`
	Order  string `json:"order"`
}
