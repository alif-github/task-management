package domain

type GetListParameter struct {
	Page   int64  `json:"page"`
	Limit  int64  `json:"limit"`
	Filter string `json:"filter"`
	Order  string `json:"order"`
}

type ContextModel struct {
	UserLoginID int64 `json:"user_login_id"`
	LimitedID   int64 `json:"limited_id"`
}
