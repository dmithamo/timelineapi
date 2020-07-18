// package dbservice avails mySQL DB functionality
package dbservice

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql" // will run pkg's init() func
)

func init() {
	// initialize schemas
	initializeSchemas()
}

// ConnectDB establishes a connection to the db
func ConnectDB(dsn *string) (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%v?parseTime=true", *dsn))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// CreateTables creates all tables in the db
func CreateTables(db *sql.DB) error {
	createTableHelper := func(tableName string) error {
		stmt, err := db.Prepare(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v %v", tableName, schemas[tableName]))
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec()
		if err != nil {
			return err
		}
		return nil
	}

	// createTables in order - because references
	err := createTableHelper("users")
	if err != nil {
		return err
	}

	err = createTableHelper("actions")
	if err != nil {
		return err
	}

	err = createTableHelper("outputs")
	if err != nil {
		return err
	}

	return nil
}

// DropTables drops the tables in the db
func DropTables(db *sql.DB) error {
	tables := []string{}
	for tableName := range schemas {
		tables = append(tables, tableName)
	}

	stmt, err := db.Prepare(fmt.Sprintf("DROP TABLE IF EXISTS %v", strings.Join(tables, ", ")))
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}
