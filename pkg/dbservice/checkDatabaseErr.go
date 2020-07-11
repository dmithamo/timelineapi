package dbservice

import (
	"fmt"
	"strings"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
)

// CheckDatabaseErr returns `better` err messages for db errors
func CheckDatabaseErr(err error, columns ...string) error {
	// verify that err is driver specific err
	driverErr, isDriverErr := err.(*mysql.MySQLError)
	switch {
	case !isDriverErr:
		return err

	// check which field reported a duplication err, if err is duplication err
	case driverErr.Number == mysqlerr.ER_DUP_ENTRY:
		return checkEachColumnHelper(driverErr.Message, "is already in use", columns)
	}

	return nil
}

// checkEachColumnHelper loops thru the list of columns to find the
// specific column reporting a db err
func checkEachColumnHelper(dbErrMessage string, semanticErrMesage string, columns []string) error {
	var specificErrMessage error
	for column := range columns {
		if strings.Contains(dbErrMessage, columns[column]) {
			specificErrMessage = fmt.Errorf("%v %v", columns[column], semanticErrMesage)
		}
	}

	return specificErrMessage
}
