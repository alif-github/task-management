package task_usecase

import (
	"context"
	"database/sql"
	"github.com/alif-github/task-management/app/serverconfig"
	"github.com/alif-github/task-management/domain"
	"github.com/alif-github/task-management/repository/task_repository/model"
	"github.com/alif-github/task-management/util"
	"time"
)

func (t taskUsecase) Add(_ context.Context, context domain.ContextModel, task *domain.TaskRequest) (err util.ErrorModel) {
	var (
		fileName  = "task_usecase_add.go"
		funcName  = "Add"
		db        = serverconfig.ServerAttribute.DBConnection
		timeNow   = time.Now()
		taskModel model.TaskModel
		dueDate   time.Time
	)

	tx, errs := db.Begin()
	if errs != nil {
		err = util.GenerateUnknownServerError(fileName, funcName, errs)
		return
	}

	defer func() {
		if err.Err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	dueDate, errs = time.Parse("2006-01-02T15:04:05Z", task.DueDateStr)
	if errs != nil {
		err = util.GenerateUnknownServerError(fileName, funcName, errs)
		return
	}

	taskModel = model.TaskModel{
		Title:       sql.NullString{String: task.Title},
		Description: sql.NullString{String: task.Description},
		DueDate:     sql.NullTime{Time: dueDate},
		Status:      sql.NullString{String: task.Status},
		UserID:      sql.NullInt64{Int64: task.UserID},
		CreatedBy:   sql.NullInt64{Int64: context.UserLoginID},
		CreatedAt:   sql.NullTime{Time: timeNow},
		UpdatedBy:   sql.NullInt64{Int64: context.UserLoginID},
		UpdatedAt:   sql.NullTime{Time: timeNow},
	}

	err = t.taskRepo.Add(nil, tx, &taskModel)
	if err.Err != nil {
		return
	}

	err = util.GenerateNonError()
	return
}
