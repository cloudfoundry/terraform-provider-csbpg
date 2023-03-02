package csbpg

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/lib/pq"
)

func createDataOwnerRole(tx *sql.Tx, cf connectionFactory) error {
	log.Println("[DEBUG] ENTRY createDataOwnerRole()")
	defer log.Println("[DEBUG] EXIT createDataOwnerRole()")

	exists, err := roleExists(tx, cf.dataOwnerRole)
	if err != nil {
		return fmt.Errorf("checking whether dataowner exists: %w", err)
	}

	if !exists {
		log.Println("[DEBUG] data owner role does not exist - creating")
		if _, err := tx.Exec(fmt.Sprintf("CREATE ROLE %s WITH NOLOGIN", pq.QuoteIdentifier(cf.dataOwnerRole))); err != nil {
			return fmt.Errorf("creating dataowner role: %w", err)
		}
	}

	log.Println("[DEBUG] granting data owner role")
	if _, err := tx.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s", pq.QuoteIdentifier(cf.database), pq.QuoteIdentifier(cf.dataOwnerRole))); err != nil {
		return fmt.Errorf("granting database privilege to datawoner role: %w", err)
	}

	if _, err := tx.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO %s", pq.QuoteIdentifier(cf.dataOwnerRole))); err != nil {
		return fmt.Errorf("granting table privilege to datawoner role: %w", err)
	}

	return nil
}
