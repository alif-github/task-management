package pgsql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/alif-github/task-management/domain"
	"github.com/alif-github/task-management/repository/task_repository/model"
	helper "github.com/alif-github/task-management/repository/util"
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

	return helper.UpdateRow(nil, tx, query, params, p.FileName, funcName)
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
	return helper.UpdateRow(nil, tx, query, params, p.FileName, funcName)
}

func (p pgsqlTaskRepository) Fetch(_ context.Context, db *sql.DB, page, limit int) ([]model.TaskModel, util.ErrorModel) {
	//TODO implement me
	panic("implement me")
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
			CONCAT(up.first_name, ' ', up.last_name) as editor, t.updated_at 
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
		&output.UpdatedAt)

	if dbErr != nil && dbErr.Error() != sql.ErrNoRows.Error() {
		err = util.GenerateInternalDBServerError(p.FileName, funcName, dbErr)
		return
	}

	err = util.GenerateNonError()
	return
}
