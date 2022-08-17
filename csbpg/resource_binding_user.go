package csbpg

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lib/pq"
)

const (
	bindingUsernameKey       = "username"
	bindingPasswordKey       = "password"
	legacyBrokerBindingGroup = "binding_group"
)

var (
	createBindingMutex sync.Mutex
	deleteBindingMutex sync.Mutex
)

func resourceBindingUser() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			bindingUsernameKey: {
				Type:     schema.TypeString,
				Required: true,
			},
			bindingPasswordKey: {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},
		CreateContext: resourceBindingUserCreate,
		ReadContext:   resourceBindingUserRead,
		UpdateContext: resourceBindingUserUpdate,
		DeleteContext: resourceBindingUserDelete,
		Description:   "TODO",
		UseJSONNumber: true,
	}
}

func resourceBindingUserCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	createBindingMutex.Lock()
	defer createBindingMutex.Unlock()

	log.Println("[DEBUG] ENTRY resourceBindingUserCreate()")
	defer log.Println("[DEBUG] EXIT resourceBindingUserCreate()")

	username := d.Get(bindingUsernameKey).(string)
	password := d.Get(bindingPasswordKey).(string)

	cf := m.(connectionFactory)

	db, err := cf.ConnectAsAdmin()
	if err != nil {
		return diag.FromErr(err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	log.Println("[DEBUG] connected")

	err = createDataOwnerRole(tx, cf)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] create binding user")
	userPresent, err := roleExists(tx, username)
	if err != nil {
		return diag.FromErr(err)
	}

	if userPresent {
		statements := []string{
			fmt.Sprintf("GRANT %s TO %s", pq.QuoteIdentifier(cf.dataOwnerRole), pq.QuoteIdentifier(username)),
		}
		legacyBrokerBindingGroupPresent, err := roleExists(tx, legacyBrokerBindingGroup)
		if err != nil {
			return diag.FromErr(err)
		}
		if legacyBrokerBindingGroupPresent {
			for _, obj := range []string{"TABLES", "SEQUENCES", "FUNCTIONS"} {
				statements = append(statements, fmt.Sprintf("ALTER DEFAULT PRIVILEGES FOR ROLE %s REVOKE ALL ON %s FROM %s", pq.QuoteIdentifier(username), obj, legacyBrokerBindingGroup))
			}
		}
		for _, statement := range statements {
			_, err := tx.Exec(statement)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	} else {
		_, err = tx.Exec(fmt.Sprintf("CREATE ROLE %s WITH LOGIN PASSWORD %s INHERIT IN ROLE %s", pq.QuoteIdentifier(username), safeQuote(password), pq.QuoteIdentifier(cf.dataOwnerRole)))
	}
	if err != nil {
		return diag.FromErr(err)
	}

	err = tx.Commit()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] setting ID %s\n", username)
	d.SetId(username)

	return nil
}

func resourceBindingUserRead(_ context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	log.Println("[DEBUG] ENTRY resourceBindingUserRead()")
	defer log.Println("[DEBUG] EXIT resourceBindingUserRead()")

	username := d.Get(bindingUsernameKey).(string)

	cf := m.(connectionFactory)

	db, err := cf.ConnectAsAdmin()
	if err != nil {
		return diag.FromErr(err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	log.Println("[DEBUG] connected")

	rows, err := db.Query(fmt.Sprintf("SELECT FROM pg_catalog.pg_roles WHERE rolname = '%s'", username))
	if err != nil {
		return diag.FromErr(err)
	}

	if !rows.Next() {
		d.SetId("")
		return nil
	}

	d.SetId(username)

	return nil
}

func resourceBindingUserUpdate(_ context.Context, _ *schema.ResourceData, _ any) diag.Diagnostics {
	return diag.FromErr(fmt.Errorf("update lifecycle not implemented"))
}

func resourceBindingUserDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	log.Println("[DEBUG] ENTRY resourceBindingUserDelete()")
	defer log.Println("[DEBUG] EXIT resourceBindingUserDelete()")

	deleteBindingMutex.Lock()
	defer deleteBindingMutex.Unlock()

	bindingUser := d.Get(bindingUsernameKey).(string)
	bindingUserPassword := d.Get(bindingPasswordKey).(string)

	cf := m.(connectionFactory)

	userDb, err := cf.ConnectAsUser(bindingUser, bindingUserPassword)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = userDb.ExecContext(ctx, fmt.Sprintf("GRANT %s TO %s", pq.QuoteIdentifier(bindingUser), pq.QuoteIdentifier(cf.username)))
	if err != nil {
		return diag.FromErr(err)
	}

	db, err := cf.ConnectAsAdmin()
	if err != nil {
		return diag.FromErr(err)
	}

	defer func(connection *sql.DB) {
		_ = connection.Close()
	}(db)

	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return diag.FromErr(err)
	}
	defer func(transaction *sql.Tx) {
		_ = transaction.Rollback()
	}(tx)

	log.Println("[DEBUG] dropping binding user")
	statements := []string{
		fmt.Sprintf("SET ROLE %s", pq.QuoteIdentifier(bindingUser)),
		fmt.Sprintf("REASSIGN OWNED BY CURRENT_USER TO %s", pq.QuoteIdentifier(cf.dataOwnerRole)),
		fmt.Sprintf("SET ROLE %s", pq.QuoteIdentifier(cf.username)),
		fmt.Sprintf("REVOKE ALL PRIVILEGES ON DATABASE %s FROM %s CASCADE;", pq.QuoteIdentifier(cf.database), pq.QuoteIdentifier(bindingUser)),
		fmt.Sprintf("DROP ROLE %s", pq.QuoteIdentifier(bindingUser)),
	}
	for _, statement := range statements {
		_, err = tx.Exec(statement)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func safeQuote(s string) string {
	return fmt.Sprintf("'%s'", strings.ReplaceAll(strings.ReplaceAll(s, `\`, `\\`), `'`, `\\`))
}

func roleExists(tx *sql.Tx, name string) (bool, error) {
	log.Println("[DEBUG] ENTRY roleExists()")
	defer log.Println("[DEBUG] EXIT roleExists()")

	rows, err := tx.Query(fmt.Sprintf("SELECT FROM pg_catalog.pg_roles WHERE rolname = '%s'", name))
	if err != nil {
		return false, fmt.Errorf("error finding role %q: %w", name, err)
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	return rows.Next(), nil
}
