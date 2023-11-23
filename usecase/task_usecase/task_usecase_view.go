package task_usecase

import (
	"context"
	"github.com/alif-github/task-management/app/serverconfig"
	"github.com/alif-github/task-management/domain"
	"github.com/alif-github/task-management/repository/task_repository/model"
	"github.com/alif-github/task-management/util"
)

func (t taskUsecase) GetByID(_ context.Context, context domain.ContextModel, id int64) (output domain.ViewTaskResponse, err util.ErrorModel) {
	var (
		fileName = "task_usecase_view.go"
		funcName = "GetByID"
		viewDB   model.TaskModel
		db       = serverconfig.ServerAttribute.DBConnection
	)

	viewDB, err = t.taskRepo.GetByID(nil, db, id)
	if err.Err != nil {
		return
	}

	if viewDB.ID.Int64 < 1 {
		err = util.GenerateNotFoundError(fileName, funcName)
		return
	}

	if context.LimitedID > 0 {
		if viewDB.UserID.Int64 != context.LimitedID {
			err = util.GenerateForbiddenError(fileName, funcName)
			return
		}
	}

	output = domain.ViewTaskResponse{
		ID:          viewDB.ID.Int64,
		Title:       viewDB.Title.String,
		Description: viewDB.Description.String,
		Status:      viewDB.Status.String,
		DueDate:     viewDB.DueDate.Time,
		User:        viewDB.User.String,
		CreatedUser: viewDB.Creator.String,
		CreatedAt:   viewDB.CreatedAt.Time,
		UpdatedUser: viewDB.Editor.String,
		UpdatedAt:   viewDB.UpdatedAt.Time,
	}

	err = util.GenerateNonError()
	return
}
