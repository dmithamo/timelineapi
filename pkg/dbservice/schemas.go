package dbservice

// schemas is the structure of the tables
var schemas map[string]string

func initializeSchemas() {
	schemas = map[string]string{
		"users": `
			(
				userID VARCHAR(50) UNIQUE NOT NULL,
				username VARCHAR(100) UNIQUE NOT NULL,
				password  VARCHAR(100) NOT NULL,
				createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
			)
		`,

		"projects": `
			(
				projectID VARCHAR(50) UNIQUE NOT NULL,
				title VARCHAR(100) UNIQUE NOT NULL,
				description TEXT NOT NULL,
				createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
			)
		`,

		"tasks": `
			(
				taskID VARCHAR(50) UNIQUE NOT NULL,
				title VARCHAR(100) UNIQUE NOT NULL,
				description TEXT NOT NULL,
				createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
			)
		`,
	}
}
