package database

import "database/sql"

type DBType int

const (
	Postgres DBType = iota
	// MySQL
	// MongoDB
)

type Database struct {
	Type    DBType // influx, mysql, sqlite,...
	SQLConn *sql.DB
	// MySQLConn ...
}

type ExecResult struct {
	RowsAffected int64
	LastInsertId int64
}

type RowType map[string]interface{}
type QueryResult []RowType
type QueryParams []interface{}
