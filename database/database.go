package database

/*-----------------------*/

func New(DatabaseType DBType, params ...string) *Database {
	var newDB Database

	newDB.Type = DatabaseType

	switch DatabaseType {
	case Postgres:
		if len(params) == 0 {
			return nil
		}
		newDB.SQLConn = NewPostgresDB(params[0])
		newDB.PostgresInit()
	}

	return &newDB
}

/*-----------------------*/

func (db *Database) Close() {
	switch db.Type {
	case Postgres:
		db.PostgresClose()
	}
}

/*-----------------------*/

func (db *Database) Insert(table string, fields RowType, tags ...map[string]string) (ExecResult, error) {

	switch db.Type {
	case Postgres:
		return db.PostgresInsert(table, fields)
	}

	return ExecResult{}, nil //TODO: provide a useful error here
}

/*-----------------------*/

func (db *Database) Update(table string, fields RowType, conditions RowType) (ExecResult, error) {

	switch db.Type {
	case Postgres:
		return db.PostgresUpdate(table, fields, conditions)
	}

	return ExecResult{}, nil //TODO: provide a useful error here
}

/*-----------------------*/

func (db *Database) Delete(table string, conditions RowType) (ExecResult, error) {

	switch db.Type {
	case Postgres:
		return db.PostgresDelete(table, conditions)
	}

	return ExecResult{}, nil //TODO: provide a useful error here
}

/*-----------------------*/

func (db *Database) Load(table string, searchOnFields RowType) (QueryResult, error) {

	switch db.Type {
	case Postgres:
		return db.PostgresLoad(table, searchOnFields)
	}

	return QueryResult{}, nil //TODO: provide a useful error here

}

/*-----------------------*/

func (db *Database) Query(query string, params QueryParams) (QueryResult, error) {

	switch db.Type {
	case Postgres:
		return db.PostgresQuery(query, params)
	}

	return QueryResult{}, nil //TODO: provide a useful error here

}

/*-----------------------*/

func (db *Database) Exec(query string, params QueryParams) (ExecResult, error) {

	switch db.Type {
	case Postgres:
		return db.PostgresExec(query, params)
	}

	return ExecResult{}, nil //TODO: provide a useful error here

}

/*-----------------------*/
