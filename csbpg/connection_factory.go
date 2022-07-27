package csbpg

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

type connectionFactory struct {
	host          string
	port          int
	username      string
	password      string
	database      string
	dataOwnerRole string
	sslClientCert *clientCertificateConfig
	sslRootCert   string
	sslMode       string
}

func (c connectionFactory) ConnectAsAdmin() (*sql.DB, error) {
	return c.connect(c.uri())
}

func (c connectionFactory) ConnectAsUser(bindingUser string, bindingUserPassword string) (*sql.DB, error) {
	return c.connect(c.uriWithCreds(bindingUser, bindingUserPassword))
}

func (c connectionFactory) connect(uri string) (*sql.DB, error) {
	db, err := sql.Open("postgres", uri)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL %q: %w", c.uriRedacted(), err)
	}

	return db, nil
}
func (c connectionFactory) uriWithCreds(username, password string) string {
	fields := map[string]string{
		"host":     c.host,
		"port":     fmt.Sprintf("%d", c.port),
		"user":     username,
		"password": password,
		"database": c.database,
		"sslmode":  c.sslMode,
	}

	if c.sslClientCert != nil {
		fields["sslinline"] = "true"
		fields["sslcert"] = fmt.Sprintf("'%s'", c.sslClientCert.Certificate)
		fields["sslkey"] = fmt.Sprintf("'%s'", c.sslClientCert.Key)
		fields["sslrootcert"] = fmt.Sprintf("'%s'", c.sslRootCert)
	}

	var s strings.Builder
	for k, v := range fields {
		s.WriteString(k)
		s.WriteRune('=')
		s.WriteString(v)
		s.WriteRune(' ')
	}
	return s.String()
}

func (c connectionFactory) uri() string {
	return c.uriWithCreds(c.username, c.password)
}

func (c connectionFactory) uriRedacted() string {
	return strings.ReplaceAll(c.uri(), c.password, "REDACTED")
}

type clientCertificateConfig struct {
	Certificate string
	Key         string
}
