package dbservice

// schemas is the structure of the tables
var schemas map[string]string

func initializeSchemas() {
	schemas = map[string]string{
		"users":    "",
		"projects": "",
		"tasks":    "",
	}
}
