// package dbservice avails mySQL DB functionality
package dbservice

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // will run pkg's init() func
)

func init() {
	// initialize schemas
	initializeSchemas()
}

// ConnectDB establishes a connection to the db
func ConnectDB(dsn *string) (*sql.DB, error) {
	if dsn == nil {
		return nil, errors.New("conn err[dsn]: invalid dsn")
	}

	conn, err := sql.Open("mysql", fmt.Sprintf("%v?parseTime=true", *dsn))
	if err != nil {
		return nil, fmt.Errorf("sql [open]: %v", err)
	}

	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("sql [ping]: %v", err)
	}

	return conn, nil
}

// CreateTable creates a table in the db with the given SQL
// Use tableName="ALL" to create all tables
// NOTE: Will drop table `tableName` if exists
func CreateTable(conn *sql.DB, tableName string) error {
	createTableHelper := func(name, schema string) error {
		// drop table first
		err := DropTable(conn, name)
		if err != nil {
			return err
		}
		_, err = conn.Exec("CREATE TABLE IF NOT EXISTS ? (?)", name, schema)
		return err
	}

	if tableName == "ALL" {
		for name, schema := range schemas {
			err := createTableHelper(name, schema)
			if err != nil {
				return fmt.Errorf("sql [create]: err creating table %v", name)
			}
		}
	}

	// else - if tableName != 'ALL', create specific table
	tableSchema, ok := schemas[tableName]
	if !ok {
		return fmt.Errorf("sql [create]: no schema found for tableName %v", tableName)
	}

	err := createTableHelper(tableName, tableSchema)
	if err != nil {
		return fmt.Errorf("sql [create]: %v", err)
	}

	return nil
}

// DropTable drops the table in the db with the given tableName
// Use tableName="ALL" to drop all tables
func DropTable(conn *sql.DB, tableName string) error {
	dropTableHelper := func(name string) error {
		_, err := conn.Exec("DROP TABLE IF EXISTS ?", name)
		return err
	}

	if tableName == "ALL" {
		for name := range schemas {
			err := dropTableHelper(name)
			if err != nil {
				return fmt.Errorf("sql [drop]: %v", err)
			}
		}
	}

	// else - if tableName != 'ALL', drop the specific table
	err := dropTableHelper(tableName)
	return fmt.Errorf("sql [drop]: %v", err)
}
