package task_usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/alif-github/task-management/app/serverconfig"
	"github.com/alif-github/task-management/domain"
	"github.com/alif-github/task-management/repository/task_repository/model"
	"github.com/alif-github/task-management/repository/util"
	util2 "github.com/alif-github/task-management/util"
	"strings"
)

func (t taskUsecase) Fetch(_ context.Context, context domain.ContextModel, task *domain.TaskRequest) (output []domain.ListTaskResponse, err util2.ErrorModel) {
	var (
		fileName      = "task_usecase_list.go"
		funcName      = "Fetch"
		db            = serverconfig.ServerAttribute.DBConnection
		searchByParam util.GetListParameterModel
		listFilter    []util.ListFilter
		listOrder     []util.ListOrder
		listDataDB    []model.TaskModel
	)

	filter := strings.Split(task.Filter, ",")
	for _, itemFilter := range filter {
		filterComponent := strings.Split(itemFilter, " ")
		if len(filterComponent) != 3 {
			err = util2.GenerateUnknownServerError(fileName, funcName, errors.New("wrong format"))
			return
		}

		if filterComponent[1] != "eq" && filterComponent[1] != "like" {
			err = util2.GenerateUnknownServerError(fileName, funcName, errors.New("wrong format must eq or like"))
			return
		}

		listFilter = append(listFilter, util.ListFilter{
			Key:      sql.NullString{String: filterComponent[0]},
			Operator: sql.NullString{String: filterComponent[1]},
			Value:    sql.NullString{String: filterComponent[2]},
		})
	}

	order := strings.Split(task.Order, ",")
	for _, itemOrder := range order {
		listOrder = append(listOrder, util.ListOrder{Order: sql.NullString{String: itemOrder}})
	}

	if context.LimitedID > 0 {
		listFilter = append(listFilter, util.ListFilter{
			Key:      sql.NullString{String: "user_id"},
			Operator: sql.NullString{String: "eq"},
			Value:    sql.NullString{String: fmt.Sprintf(`%d`, context.LimitedID)},
		})
	}

	searchByParam = util.GetListParameterModel{
		Page:   sql.NullInt64{Int64: task.Page},
		Limit:  sql.NullInt64{Int64: task.Limit},
		Filter: listFilter,
		Order:  listOrder,
	}

	listDataDB, err = t.taskRepo.Fetch(nil, db, searchByParam)
	if err.Err != nil {
		return
	}

	if len(listDataDB) > 0 {
		for _, itemListDataDB := range listDataDB {
			output = append(output, domain.ListTaskResponse{
				ID:      itemListDataDB.ID.Int64,
				Title:   itemListDataDB.Title.String,
				Status:  itemListDataDB.Status.String,
				DueDate: itemListDataDB.DueDate.Time,
				User:    itemListDataDB.User.String,
			})
		}
	}

	err = util2.GenerateNonError()
	return
}
