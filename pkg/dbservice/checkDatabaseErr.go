package dbservice

import (
	"fmt"
	"strings"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
)

// CheckDatabaseErr returns `better` err messages for db errors
func CheckDatabaseErr(err error, uniqueColum ...string) error {
	// verify that err is driver specific err
	driverErr, isDriverErr := err.(*mysql.MySQLError)
	switch {
	case !isDriverErr:
		return err

	// check which field reported a duplication err, if err is duplication err
	case driverErr.Number == mysqlerr.ER_DUP_ENTRY:
		return fmt.Errorf("%v is already in use", strings.Join(uniqueColum, ", "))

	default:
		return nil
	}
}
