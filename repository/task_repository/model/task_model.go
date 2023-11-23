package model

import "database/sql"

type TaskModel struct {
	ID          sql.NullInt64
	Title       sql.NullString
	Description sql.NullString
	DueDate     sql.NullTime
	Status      sql.NullString
	UserID      sql.NullInt64
	User        sql.NullString
	CreatedBy   sql.NullInt64
	Creator     sql.NullString
	CreatedAt   sql.NullTime
	UpdatedBy   sql.NullInt64
	Editor      sql.NullString
	UpdatedAt   sql.NullTime
}
