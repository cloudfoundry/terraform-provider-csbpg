package csbpg

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/lib/pq"
)

func createDataOwnerRole(tx *sql.Tx, cf connectionFactory) error {
	exists, err := roleExists(tx, cf.dataOwnerRole)
	if err != nil {
		return err
	}

	if !exists {
		log.Println("[DEBUG] data owner role does not exist - creating")
		_, err = tx.Exec(fmt.Sprintf("CREATE ROLE %s WITH NOLOGIN", pq.QuoteIdentifier(cf.dataOwnerRole)))

		if err != nil {
			return err
		}
	}

	log.Println("[DEBUG] granting data owner role")
	_, err = tx.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s", pq.QuoteIdentifier(cf.database), pq.QuoteIdentifier(cf.dataOwnerRole)))
	if err != nil {
		return err
	}

	_, err = tx.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO %s", pq.QuoteIdentifier(cf.dataOwnerRole)))
	return err
}
