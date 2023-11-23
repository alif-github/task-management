package pgsql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/alif-github/task-management/domain"
	util_repo "github.com/alif-github/task-management/repository/util"
	"github.com/alif-github/task-management/util"
)

type pgsqlUserRepository struct {
	FileName  string
	TableName string
	Conn      *sql.DB
}

func NewPgsqlUserRepository(conn *sql.DB) domain.UserRepository {
	return &pgsqlUserRepository{
		FileName:  "user_pgsql.go",
		TableName: "user",
		Conn:      conn,
	}
}

func (p pgsqlUserRepository) Store(_ context.Context, tx *sql.Tx, users *domain.User) util.ErrorModel {
	funcName := "Store"
	query := fmt.Sprintf(`
		INSERT INTO %s 
		(
		 username, password, first_name, 
		 last_name, email, created_at, 
		 updated_at, created_by, updated_by, 
		 role_id
		) 
		VALUES 
		(
		    $1, $2, $3, $4, $5, $6, 
		    $7, $8, $9, $10
		) returning id `,
		p.TableName)

	params := []interface{}{
		users.Username, users.Password, users.FirstName,
		users.LastName, users.Email, users.CreatedAt,
		users.UpdatedAt, users.CreatedBy, users.UpdatedBy,
		users.RoleID}
	result := tx.QueryRow(query, params...)
	errs := result.Scan(&users.ID)

	if errs != nil && errs.Error() != sql.ErrNoRows.Error() {
		return util.GenerateInternalDBServerError(p.FileName, funcName, errs)
	}

	return util.GenerateNonError()
}

func (p pgsqlUserRepository) Update(_ context.Context, tx *sql.Tx, users *domain.User) util.ErrorModel {
	funcName := "Update"
	query := fmt.Sprintf(`
		UPDATE %s SET 
		username = $1, password = $2, first_name = $3, 
		last_name = $4, email = $5, role_id = $6,
		updated_by = $7, updated_at = $8 
		WHERE id = $9 AND deleted = FALSE `,
		p.TableName)

	params := []interface{}{
		users.Username, users.Password, users.FirstName,
		users.LastName, users.Email, users.RoleID,
		users.UpdatedBy, users.UpdatedAt,
		users.ID}

	return util_repo.UpdateRow(nil, tx, query, params, p.FileName, funcName)
}

func (p pgsqlUserRepository) Delete(ctx context.Context, tx *sql.Tx, users *domain.User) util.ErrorModel {
	funcName := "Delete"
	query := fmt.Sprintf(`
		UPDATE %s SET deleted = TRUE, updated_at = $1, updated_by = $2 WHERE id = $3 `,
		p.TableName)

	params := []interface{}{users.UpdatedAt, users.UpdatedBy, users.ID}
	return util_repo.UpdateRow(ctx, tx, query, params, p.FileName, funcName)
}

func (p pgsqlUserRepository) Fetch(_ context.Context, db *sql.DB, page, limit int) (output []domain.User, err util.ErrorModel) {
	funcName := "Fetch"
	query := fmt.Sprintf(`
		SELECT 
		    id, first_name, last_name, 
		    username, email 
		FROM %s 
		ORDER BY id 
		LIMIT $1 OFFSET $2 `,
		p.TableName)

	params := []interface{}{limit, util_repo.CountOffset(page, limit)}
	results, errs := db.Query(query, params...)
	if errs != nil {
		err = util.GenerateNonError()
		return
	}

	var outputTemp []interface{}
	outputTemp, err = util_repo.GetRows(results, func(rows *sql.Rows) (interface{}, error) {
		var temp domain.User
		errors := rows.Scan(
			&temp.ID, &temp.FirstName, &temp.LastName,
			&temp.Username, &temp.Email)
		return temp, errors
	})
	if err.Err != nil {
		return
	}

	if outputTemp == nil {
		err = util.GenerateNotFoundError(p.FileName, funcName)
		return
	}

	for _, outputTempItem := range outputTemp {
		repo := outputTempItem.(domain.User)
		output = append(output, repo)
	}

	err = util.GenerateNonError()
	return
}

func (p pgsqlUserRepository) GetByID(_ context.Context, db *sql.DB, id int64) (output domain.User, err util.ErrorModel) {
	funcName := "GetByID"
	query := fmt.Sprintf(` 
		SELECT 
		    id, first_name, last_name, 
		    username, email, created_at, 
		    updated_at, deleted 
		FROM %s 
		WHERE id = $1 AND deleted = FALSE `,
		p.TableName)

	params := []interface{}{id}
	results := db.QueryRow(query, params...)
	dbErr := results.Scan(
		&output.ID, &output.FirstName, &output.LastName,
		&output.Username, &output.Email,
		&output.CreatedAt, &output.UpdatedAt,
		&output.Deleted)
	if dbErr != nil && dbErr.Error() != sql.ErrNoRows.Error() {
		err = util.GenerateInternalDBServerError(p.FileName, funcName, dbErr)
		return
	}

	err = util.GenerateNonError()
	return
}

func (p pgsqlUserRepository) GetByName(_ context.Context, db *sql.DB, name string) (output domain.User, err util.ErrorModel) {
	funcName := "GetByName"
	query := fmt.Sprintf(` 
		SELECT id, first_name, last_name, username, email 
		FROM %s 
		WHERE CONCAT(first_name, ' ', last_name) = $1 AND deleted = FALSE `,
		p.TableName)

	params := []interface{}{name}
	results := db.QueryRow(query, params...)
	dbErr := results.Scan(
		&output.ID, &output.FirstName, &output.LastName,
		&output.Username, &output.Email)

	if dbErr != nil && dbErr.Error() != sql.ErrNoRows.Error() {
		err = util.GenerateInternalDBServerError(p.FileName, funcName, dbErr)
		return
	}

	err = util.GenerateNonError()
	return
}

func (p pgsqlUserRepository) GetLogin(_ context.Context, db *sql.DB, users *domain.User) (output domain.User, err util.ErrorModel) {
	funcName := "GetLogin"
	query := fmt.Sprintf(` 
		SELECT u.id, u.password, r.role_name
		FROM %s u 
		INNER JOIN role r ON u.role_id = r.id
		WHERE u.username = $1 AND u.deleted = FALSE AND r.deleted = FALSE `,
		p.TableName)

	params := []interface{}{users.Username}
	results := db.QueryRow(query, params...)
	dbErr := results.Scan(&output.ID, &output.Password, &output.Role)
	if dbErr != nil && dbErr.Error() != sql.ErrNoRows.Error() {
		err = util.GenerateInternalDBServerError(p.FileName, funcName, dbErr)
		return
	}

	err = util.GenerateNonError()
	return
}
