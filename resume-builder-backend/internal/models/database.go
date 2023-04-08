// internal/models/database.go

package models

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
)

func Migrate(db *sql.DB) error {
	// Read schema.sql file
	content, err := ioutil.ReadFile("./schema.sql")
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	// Execute schema.sql
	if _, err := db.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	log.Println("Database schema migrated successfully")
	return nil
}
