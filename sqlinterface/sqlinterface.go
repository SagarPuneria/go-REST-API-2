package SqlInterface

import (
	"database/sql"
	"sync"

	// Register mysql driver
	_ "github.com/go-sql-driver/mysql"
)

type MySqldb struct {
	db        *sql.DB
	tx        *sql.Tx
	dbRWMutex sync.RWMutex
}

func CreateDataBase(DNS string, quries ...string) (*MySqldb, error) {
	var dbConn = new(MySqldb)
	dbConn.dbRWMutex.Lock()
	defer dbConn.dbRWMutex.Unlock()
	var err error
	dbConn.db, err = sql.Open("mysql", DNS)
	if err != nil {
		return nil, err
	}
	for _, query := range quries {
		_, err := dbConn.db.Exec(query)
		if err != nil {
			dbConn.Close()
			return dbConn, err
		}
	}
	return dbConn, nil
}

func (DBObject *MySqldb) Close() {
	DBObject.db.Close()
}

func (DBObject *MySqldb) BeginTx() error {
	var err error
	DBObject.tx, err = DBObject.db.Begin()
	if err != nil {
		return err
	}
	return nil
}

func (DBObject *MySqldb) ExecuteQuery(strQuery string) error {
	_, execErr := DBObject.tx.Exec(strQuery)
	if execErr != nil {
		if rollbackErr := DBObject.tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
	}
	return nil
}

func (DBObject *MySqldb) CommitTx() error {
	DBObject.dbRWMutex.Lock()
	defer DBObject.dbRWMutex.Unlock()
	if err := DBObject.tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (DBObject *MySqldb) SelectQuery(strQuery string) (*sql.Rows, error) {
	DBObject.dbRWMutex.Lock()
	defer DBObject.dbRWMutex.Unlock()

	rows, err := DBObject.db.Query(strQuery)

	return rows, err
}
