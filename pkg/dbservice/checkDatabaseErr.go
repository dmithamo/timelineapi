package dbservice

import (
	"fmt"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
)

// CheckDatabaseErr returns `better` err messages for db errors
func CheckDatabaseErr(err error, uniqueColum ...string) error {
	if err == nil {
		return nil
	}

	// db driver errs that we are currently checking
	driverErrs := map[uint16]error{
		mysqlerr.ER_DUP_ENTRY: fmt.Errorf("%v is already in use", uniqueColum),
	}

	// verify that err is driver specific err
	driverErr, isDriverErr := err.(*mysql.MySQLError)
	if !isDriverErr {
		return err
	}

	formattedErr, exists := driverErrs[driverErr.Number]
	if !exists {
		return err
	}

	return formattedErr
}
