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

type Response struct {
	RequestID string      `json:"request_id"`
	Status    bool        `json:"status"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
}

type Pagination struct {
	Page   int    `form:"page" binding:"required,gte=1"`
	Limit  int    `form:"limit" binding:"required,gte=1"`
	Filter string `form:"filter"`
	Order  string `form:"order"`
}
