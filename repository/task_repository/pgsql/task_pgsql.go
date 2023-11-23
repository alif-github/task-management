package pgsql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/alif-github/task-management/domain"
	"github.com/alif-github/task-management/repository/task_repository/model"
	util_repo "github.com/alif-github/task-management/repository/util"
	"github.com/alif-github/task-management/util"
)

type pgsqlTaskRepository struct {
	FileName  string
	TableName string
	Conn      *sql.DB
}

func NewPgsqlTaskRepository(conn *sql.DB) domain.TaskRepository {
	return &pgsqlTaskRepository{
		FileName:  "task_pgsql.go",
		TableName: "task",
		Conn:      conn,
	}
}

func (p pgsqlTaskRepository) Add(_ context.Context, tx *sql.Tx, task *model.TaskModel) util.ErrorModel {
	var (
		funcName = "Add"
		query    string
		param    []interface{}
	)

	query = fmt.Sprintf(`
		INSERT INTO %s 
			(
			 title, description, due_date, 
			 status, user_id, created_by, 
			 created_at, updated_by, updated_at
			) 
		VALUES 
		($1, $2, $3, $4, $5, $6, $7, $8, $9) 
		RETURNING id`,
		p.TableName)

	if task.Title.Valid {
		param = append(param, task.Title.String)
	} else {
		param = append(param, nil)
	}

	if task.Description.Valid {
		param = append(param, task.Description.String)
	} else {
		param = append(param, nil)
	}

	if task.DueDate.Valid {
		param = append(param, task.DueDate.Time)
	} else {
		param = append(param, nil)
	}

	if task.Status.Valid {
		param = append(param, task.Status.String)
	} else {
		param = append(param, nil)
	}

	if task.UserID.Valid {
		param = append(param, task.UserID.Int64)
	} else {
		param = append(param, nil)
	}

	param = append(param,
		task.CreatedBy.Int64, task.CreatedAt.Time, task.UpdatedBy.Int64,
		task.UpdatedAt.Time)

	result := tx.QueryRow(query, param...)
	errs := result.Scan(&task.ID)
	if errs != nil && errs.Error() != sql.ErrNoRows.Error() {
		return util.GenerateInternalDBServerError(p.FileName, funcName, errs)
	}

	return util.GenerateNonError()
}

func (p pgsqlTaskRepository) Update(_ context.Context, tx *sql.Tx, task model.TaskModel) util.ErrorModel {
	var (
		funcName = "Update"
		query    string
	)

	query = fmt.Sprintf(`
		UPDATE %s 
		SET status = $1, updated_by = $2, updated_at = $3 
		WHERE id = $4 AND deleted = FALSE `,
		p.TableName)

	params := []interface{}{
		task.Status.String, task.UpdatedBy.Int64, task.UpdatedAt.Time,
		task.ID.Int64}

	return util_repo.UpdateRow(nil, tx, query, params, p.FileName, funcName)
}

func (p pgsqlTaskRepository) Delete(_ context.Context, tx *sql.Tx, task model.TaskModel) util.ErrorModel {
	var (
		funcName = "Delete"
		query    string
	)

	query = fmt.Sprintf(`
		UPDATE %s 
		SET deleted = FALSE, updated_by = $1, updated_at = $2 
		WHERE id = $3 AND deleted = FALSE `,
		p.TableName)

	params := []interface{}{task.UpdatedBy.Int64, task.UpdatedAt.Time, task.ID.Int64}
	return util_repo.UpdateRow(nil, tx, query, params, p.FileName, funcName)
}

func (p pgsqlTaskRepository) Fetch(_ context.Context, db *sql.DB, param util_repo.GetListParameterModel) (output []model.TaskModel, err util.ErrorModel) {
	var (
		funcName = "Fetch"
		query    string
		params   []interface{}
	)

	query = fmt.Sprintf(`
		SELECT 
		    t.id, t.title, t.due_date, 
		    t.status, CONCAT(u.first_name, ' ', u.last_name) as person 
		FROM %s t 
		INNER JOIN "user" u ON t.user_id = u.id 
		WHERE t.deleted = FALSE `,
		p.TableName)

	if len(param.Filter) > 0 {
		for i, itemFilter := range param.Filter {
			switch itemFilter.Key.String {
			case "title":
				param.Filter[i].Key.String = "t.title"
				if param.Filter[i].Operator.String == "like" {
					query += " AND LOWER(t.title) like '%" + param.Filter[i].Value.String + "%'"
				} else if param.Filter[i].Operator.String == "eq" {
					query += " AND t.title = '" + param.Filter[i].Value.String + "'"
				} else {
					err = util.GenerateInternalDBServerError(p.FileName, funcName, errors.New("title only eq and like"))
					return
				}
			case "status":
				param.Filter[i].Key.String = "t.status"
				if param.Filter[i].Operator.String == "eq" {
					query += " AND t.status = '" + param.Filter[i].Value.String + "'"
				} else {
					err = util.GenerateInternalDBServerError(p.FileName, funcName, errors.New("status only eq"))
					return
				}
			case "user_id":
				param.Filter[i].Key.String = "t.user_id"
				if param.Filter[i].Operator.String == "eq" {
					query += " AND t.user_id = " + param.Filter[i].Value.String
				} else {
					err = util.GenerateInternalDBServerError(p.FileName, funcName, errors.New("user_id only eq"))
					return
				}
			default:
			}
		}
	}

	if len(param.Order) > 0 {
		query += " ORDER BY "
		for i, itemOrder := range param.Order {
			query += itemOrder.Order.String
			if len(param.Order)-(i+1) > 0 {
				query += ", "
			}
		}
	}

	query += fmt.Sprintf(`LIMIT $1 OFFSET $2`)
	params = []interface{}{param.Limit.Int64, util_repo.CountOffset(int(param.Page.Int64), int(param.Limit.Int64))}

	rows, errs := db.Query(query, params...)
	if errs != nil && errs.Error() != sql.ErrNoRows.Error() {
		err = util.GenerateInternalDBServerError(p.FileName, funcName, errs)
		return
	}

	var result []interface{}
	result, err = util_repo.GetRows(rows, func(rws *sql.Rows) (interface{}, error) {
		var temp model.TaskModel
		errs1 := rws.Scan(
			&temp.ID, &temp.Title, &temp.DueDate,
			&temp.Status, &temp.User)
		return temp, errs1
	})

	if len(result) > 0 {
		for _, itemResult := range result {
			output = append(output, itemResult.(model.TaskModel))
		}
	}

	err = util.GenerateNonError()
	return
}

func (p pgsqlTaskRepository) GetByID(_ context.Context, db *sql.DB, id int64) (output model.TaskModel, err util.ErrorModel) {
	var (
		funcName = "GetByID"
		query    string
	)

	query = fmt.Sprintf(`
		SELECT 
			t.id, t.title, t.description, 
			t.due_date, t.status, CONCAT(ut.first_name, ' ', ut.last_name) as person, 
			CONCAT(uc.first_name, ' ', uc.last_name) as creator, t.created_at, 
			CONCAT(up.first_name, ' ', up.last_name) as editor, t.updated_at, 
			t.user_id
		FROM %s t 
		INNER JOIN "user" ut ON t.user_id = ut.id 
		INNER JOIN "user" uc ON t.user_id = uc.id 
		INNER JOIN "user" up ON t.user_id = up.id 
		WHERE 
		    t.id = $1 AND t.deleted = FALSE `,
		p.TableName)

	params := []interface{}{id}
	results := db.QueryRow(query, params...)
	dbErr := results.Scan(
		&output.ID, &output.Title, &output.Description,
		&output.DueDate, &output.Status, &output.User,
		&output.Creator, &output.CreatedAt, &output.Editor,
		&output.UpdatedAt, &output.UserID)

	if dbErr != nil && dbErr.Error() != sql.ErrNoRows.Error() {
		err = util.GenerateInternalDBServerError(p.FileName, funcName, dbErr)
		return
	}

	err = util.GenerateNonError()
	return
}
