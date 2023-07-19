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
		Description:   "Represents a CloudFoundry binding in PostgreSQL",
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
	err := sqlUserCreate(ctx, username, password, m)
	if err != nil {
		return err
	}
	d.SetId(username)
	return nil
}

func sqlUserCreate(ctx context.Context, username, password string, m any) diag.Diagnostics {
	cf := m.(connectionFactory)

	db, err := cf.ConnectAsAdmin()
	if err != nil {
		return diag.Errorf("connecting as admin: %s", err)
	}
	defer func() {
		_ = db.Close()
	}()

	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return diag.Errorf("starting transaction: %s", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	log.Println("[DEBUG] connected")
	if err := grantAllPrivilegesToPublicSchema(tx, cf); err != nil {
		return diag.FromErr(err)
	}

	userPresent, err := roleExists(tx, username)
	if err != nil {
		return diag.Errorf("checking whether binding user exists: %s", err)
	}

	if userPresent {
		// The following instruction ensures admin has access and permissions over any objects created by the legacy user
		// We need to do this before executing the createDataOwnerRole because there are some instructions in that function
		// which can fail if there are tables in public schema for which the admin user doesn't have elevated permissions
		if _, err = tx.Exec(fmt.Sprintf("GRANT %s TO %s", pq.QuoteIdentifier(username), pq.QuoteIdentifier(cf.username))); err != nil {
			return diag.Errorf("grant admin the right to impersonate legecy role and manipulate its objects: %s", err)
		}
	}

	if err := createDataOwnerRole(tx, cf); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] create binding user")

	if userPresent {
		statements := []string{
			fmt.Sprintf("GRANT %s TO %s", pq.QuoteIdentifier(cf.dataOwnerRole), pq.QuoteIdentifier(username)),
		}
		legacyBrokerBindingGroupPresent, err := roleExists(tx, legacyBrokerBindingGroup)
		if err != nil {
			return diag.Errorf("checking whether legacy binding group exists: %s", err)
		}
		if legacyBrokerBindingGroupPresent {
			for _, obj := range []string{"TABLES", "SEQUENCES", "FUNCTIONS"} {
				statements = append(statements, fmt.Sprintf("ALTER DEFAULT PRIVILEGES FOR ROLE %s REVOKE ALL ON %s FROM %s", pq.QuoteIdentifier(username), obj, legacyBrokerBindingGroup))
			}
		}
		for _, statement := range statements {
			if _, err := tx.Exec(statement); err != nil {
				return diag.Errorf("running statement %q: %s", statement, err)
			}
		}
	} else {
		if _, err := tx.Exec(fmt.Sprintf("CREATE ROLE %s WITH LOGIN PASSWORD %s INHERIT IN ROLE %s", pq.QuoteIdentifier(username), safeQuote(password), pq.QuoteIdentifier(cf.dataOwnerRole))); err != nil {
			return diag.Errorf("creating binding role: %s", err)
		}
		if _, err = tx.Exec(fmt.Sprintf("GRANT %s TO %s", pq.QuoteIdentifier(username), pq.QuoteIdentifier(cf.username))); err != nil {
			return diag.Errorf("grant admin the right to impersonate new role and manipulate its objects: %s", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return diag.Errorf("committing transaction: %s", err)
	}

	log.Printf("[DEBUG] setting ID %s\n", username)

	return nil
}

func resourceBindingUserRead(_ context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	log.Println("[DEBUG] ENTRY resourceBindingUserRead()")
	defer log.Println("[DEBUG] EXIT resourceBindingUserRead()")

	username := d.Get(bindingUsernameKey).(string)

	cf := m.(connectionFactory)

	db, err := cf.ConnectAsAdmin()
	if err != nil {
		return diag.Errorf("connecting as admin: %s", err)
	}
	defer func() {
		_ = db.Close()
	}()
	log.Println("[DEBUG] connected")

	exists, err := roleExists(db, username)
	switch {
	case err != nil:
		return diag.Errorf("querying for existing role: %s", err)
	case exists:
		d.SetId(username)
	default:
		d.SetId("")
	}

	return nil
}

func resourceBindingUserUpdate(_ context.Context, _ *schema.ResourceData, _ any) diag.Diagnostics {
	log.Println("[DEBUG] ENTRY resourceBindingUserUpdate()")
	defer log.Println("[DEBUG] EXIT resourceBindingUserUpdate()")
	return diag.Errorf("update lifecycle not implemented")
}

func resourceBindingUserDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	log.Println("[DEBUG] ENTRY resourceBindingUserDelete()")
	defer log.Println("[DEBUG] EXIT resourceBindingUserDelete()")

	deleteBindingMutex.Lock()
	defer deleteBindingMutex.Unlock()

	bindingUser := d.Get(bindingUsernameKey).(string)
	bindingUserPassword := d.Get(bindingPasswordKey).(string)
	err := sqlUserDelete(ctx, bindingUser, bindingUserPassword, m)
	if err != nil {
		return err
	}
	return nil
}

func sqlUserDelete(ctx context.Context, bindingUser, bindingUserPassword string, m any) diag.Diagnostics {
	cf := m.(connectionFactory)

	userDb, err := cf.ConnectAsUser(bindingUser, bindingUserPassword)
	if err != nil {
		return diag.Errorf("connecting as binding user: %s", err)
	}
	defer func() {
		_ = userDb.Close()
	}()

	if _, err := userDb.ExecContext(ctx, fmt.Sprintf("GRANT %s TO %s", pq.QuoteIdentifier(bindingUser), pq.QuoteIdentifier(cf.username))); err != nil {
		return diag.Errorf("granting admin user access to binding user: %s", err)
	}

	db, err := cf.ConnectAsAdmin()
	if err != nil {
		return diag.Errorf("connecting as admin: %s", err)
	}
	defer func() {
		_ = db.Close()
	}()

	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return diag.Errorf("starting transaction: %s", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	log.Println("[DEBUG] dropping binding user")
	statements := []string{
		fmt.Sprintf("SET ROLE %s", pq.QuoteIdentifier(bindingUser)),
		fmt.Sprintf("REASSIGN OWNED BY CURRENT_USER TO %s", pq.QuoteIdentifier(cf.dataOwnerRole)),
		fmt.Sprintf("SET ROLE %s", pq.QuoteIdentifier(cf.username)),
		fmt.Sprintf("REVOKE ALL PRIVILEGES ON DATABASE %s FROM %s CASCADE;", pq.QuoteIdentifier(cf.database), pq.QuoteIdentifier(bindingUser)),
		fmt.Sprintf("DROP ROLE %s", pq.QuoteIdentifier(bindingUser)),
	}
	for _, statement := range statements {
		if _, err = tx.Exec(statement); err != nil {
			return diag.Errorf("running statement %q: %s", statement, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return diag.Errorf("committing transaction: %s", err)
	}

	return nil
}

func safeQuote(s string) string {
	return fmt.Sprintf("'%s'", strings.ReplaceAll(strings.ReplaceAll(s, `\`, `\\`), `'`, `\\`))
}

type querier interface {
	Query(query string, args ...any) (*sql.Rows, error)
}

func roleExists(q querier, name string) (bool, error) {
	log.Println("[DEBUG] ENTRY roleExists()")
	defer log.Println("[DEBUG] EXIT roleExists()")

	rows, err := q.Query(fmt.Sprintf("SELECT FROM pg_catalog.pg_roles WHERE rolname = '%s'", name))
	if err != nil {
		return false, fmt.Errorf("error finding role %q: %w", name, err)
	}
	defer func() {
		_ = rows.Close()
	}()

	return rows.Next(), nil
}
