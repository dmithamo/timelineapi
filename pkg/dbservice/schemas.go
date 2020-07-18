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
				actionID BINARY(16) PRIMARY KEY,
				title VARCHAR(50) UNIQUE NOT NULL,
				description TEXT NOT NULL,
				isArchived BOOLEAN DEFAULT FALSE,
				createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				userID BINARY(16) NOT NULL,
				FOREIGN KEY (userID)
					REFERENCES users(userID)
					ON DELETE CASCADE
			)
		`,

		"outputs": `
			(
				outputID BINARY(16) PRIMARY KEY,
				title VARCHAR(50) UNIQUE NOT NULL,
				description TEXT NOT NULL,
				isArchived BOOLEAN DEFAULT FALSE,
				createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				actionID BINARY(16),
				FOREIGN KEY (actionID)
					REFERENCES actions(actionID)
					ON DELETE CASCADE
			)
		`,
	}
}
