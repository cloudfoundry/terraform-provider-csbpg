package csbpg

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/lib/pq"
)

func relaxPublicSchemaRestrictions(tx *sql.Tx, cf connectionFactory) error {
	log.Println("[DEBUG] make admin user owner of the public schema")
	if _, err := tx.Exec(fmt.Sprintf("ALTER SCHEMA public OWNER TO %s", pq.QuoteIdentifier(cf.username))); err != nil {
		return fmt.Errorf("make schema public be owned by public admin user: %s", err)
	}
	log.Println("[DEBUG] granting permission on schema public to all users (required since postgres 15)")
	if _, err := tx.Exec(fmt.Sprintf("GRANT ALL ON SCHEMA PUBLIC TO PUBLIC")); err != nil {
		return fmt.Errorf("granting all privileges on schema public to all users: %s", err)
	}

	return nil
}
