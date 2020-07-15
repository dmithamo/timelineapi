package dbservice

// schemas is the structure of the tables
var schemas map[string]string

// initializeSchemas assigns value to the schemas variable above
func initializeSchemas() {
	schemas = map[string]string{
		"users": `
			(
				userID BINARY(16) PRIMARY KEY,
				username VARCHAR(100) UNIQUE NOT NULL,
				password  VARCHAR(100) NOT NULL,
				createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
			)
		`,

		"actions": `
			(
				actionID VARCHAR(100) UNIQUE NOT NULL,
				title VARCHAR(50) UNIQUE NOT NULL,
				description TEXT NOT NULL,
				createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
			)
		`,

		"outputs": `
			(
				outputID VARCHAR(100) UNIQUE NOT NULL,
				title VARCHAR(50) UNIQUE NOT NULL,
				description TEXT NOT NULL,
				createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
			)
		`,
	}
}
