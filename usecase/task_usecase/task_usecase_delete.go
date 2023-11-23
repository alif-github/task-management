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

func (t taskUsecase) Delete(_ context.Context, context domain.ContextModel, task *domain.TaskRequest) (err util.ErrorModel) {
	var (
		fileName  = "task_usecase_delete.go"
		funcName  = "Delete"
		db        = serverconfig.ServerAttribute.DBConnection
		timeNow   = time.Now()
		taskModel model.TaskModel
		taskOnDB  model.TaskModel
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

	taskModel = model.TaskModel{
		ID:        sql.NullInt64{Int64: task.ID},
		UpdatedBy: sql.NullInt64{Int64: context.UserLoginID},
		UpdatedAt: sql.NullTime{Time: timeNow},
	}

	taskOnDB, err = t.taskRepo.GetByID(nil, db, taskModel.ID.Int64)
	if err.Err != nil {
		return
	}

	if taskOnDB.ID.Int64 < 1 {
		err = util.GenerateNotFoundError(fileName, funcName)
		return
	}

	if context.LimitedID > 0 {
		if taskOnDB.UserID.Int64 != context.LimitedID {
			err = util.GenerateForbiddenError(fileName, funcName)
			return
		}
	}

	err = t.taskRepo.Delete(nil, tx, taskModel)
	if err.Err != nil {
		return
	}

	err = util.GenerateNonError()
	return
}
