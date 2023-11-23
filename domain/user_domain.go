package domain

import (
	"context"
	"database/sql"
	"github.com/alif-github/task-management/util"
	"time"
)

type User struct {
	ID           int64  `json:"id"`
	FirstName    string `json:"first_name" binding:"required,min=1,max=50"`
	LastName     string `json:"last_name" binding:"max=100"`
	Username     string `json:"username" binding:"required,min=4,max=20"`
	Password     string `json:"password" binding:"required,min=8,max=20"`
	Email        string `json:"email" binding:"required,email"`
	RoleID       int64  `json:"role_id" binding:"required,gt=0"`
	CreatedBy    int64  `json:"created_by"`
	CreatedAtStr string `json:"created_at"`
	UpdatedBy    int64  `json:"updated_by"`
	UpdatedAtStr string `json:"updated_at"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Deleted      bool
}

type Login struct {
	Username string `json:"username" binding:"required,min=4,max=20"`
	Password string `json:"password" binding:"required,min=8,max=20"`
}

type UserUsecase interface {
	Store(ctx context.Context, users *User) util.ErrorModel
	Update(ctx context.Context, users *User) util.ErrorModel
	Delete(ctx context.Context, users *User) util.ErrorModel
	Fetch(ctx context.Context, page, limit int) ([]User, util.ErrorModel)
	GetByID(ctx context.Context, id int64) (User, util.ErrorModel)
	GetLogin(ctx context.Context, users *User) (User, util.ErrorModel)
}

type UserRepository interface {
	Store(ctx context.Context, tx *sql.Tx, users *User) util.ErrorModel
	Update(ctx context.Context, tx *sql.Tx, users *User) util.ErrorModel
	Delete(ctx context.Context, tx *sql.Tx, users *User) util.ErrorModel
	Fetch(ctx context.Context, db *sql.DB, page, limit int) ([]User, util.ErrorModel)
	GetByID(ctx context.Context, db *sql.DB, id int64) (User, util.ErrorModel)
	GetByName(ctx context.Context, db *sql.DB, name string) (User, util.ErrorModel)
	GetLogin(ctx context.Context, db *sql.DB, users *User) (User, util.ErrorModel)
}
