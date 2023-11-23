package domain

import (
	"context"
	"database/sql"
	"github.com/alif-github/task-management/repository/task_repository/model"
	param "github.com/alif-github/task-management/repository/util"
	"github.com/alif-github/task-management/util"
	"time"
)

type TaskRequest struct {
	ID           int64  `json:"id"`
	Title        string `json:"title" binding:"required,min=1,max=255"`
	Description  string `json:"description"`
	DueDateStr   string `json:"due_date" binding:"required"`
	Status       string `json:"status" binding:"required,oneof=New InProgress Pending Done"`
	UserID       int64  `json:"user_id" binding:"required,gt=0"`
	CreatedBy    int64  `json:"created_by"`
	CreatedAtStr string `json:"created_at"`
	UpdatedBy    int64  `json:"updated_by"`
	UpdatedAtStr string `json:"updated_at"`
	DueDate      time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type ListTaskResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	DueDate     time.Time `json:"due_date"`
	User        string    `json:"user"`
}

type ViewTaskResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	DueDate     time.Time `json:"due_date"`
	User        string    `json:"user"`
	CreatedUser string    `json:"created_user"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedUser string    `json:"updated_user"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TaskUsecase interface {
	Add(ctx context.Context, task *TaskRequest) util.ErrorModel
	Update(ctx context.Context, task *TaskRequest) util.ErrorModel
	Delete(ctx context.Context, task *TaskRequest) util.ErrorModel
	Fetch(ctx context.Context, task *TaskRequest) ([]ListTaskResponse, util.ErrorModel)
	GetByID(ctx context.Context, id int64) (ViewTaskResponse, util.ErrorModel)
}

type TaskRepository interface {
	Add(ctx context.Context, tx *sql.Tx, task *model.TaskModel) util.ErrorModel
	Update(ctx context.Context, tx *sql.Tx, task model.TaskModel) util.ErrorModel
	Delete(ctx context.Context, tx *sql.Tx, task model.TaskModel) util.ErrorModel
	Fetch(ctx context.Context, db *sql.DB, param param.GetListParameterModel) ([]model.TaskModel, util.ErrorModel)
	GetByID(ctx context.Context, db *sql.DB, id int64) (model.TaskModel, util.ErrorModel)
}
