// Package csbpg is a Terraform provider specialised for CSB PostgreSQL bindings
package csbpg

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	dataOwnerRoleKey = "data_owner_role"
	databaseKey      = "database"
	passwordKey      = "password"
	usernameKey      = "username"
	portKey          = "port"
	hostKey          = "host"
	sslModeKey       = "sslmode"
	clientCertKey    = "clientcert"
	sslRootCertKey   = "sslrootcert"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			portKey: {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IsPortNumber,
			},
		},
		ConfigureContextFunc: providerConfigure,
		ResourcesMap: map[string]*schema.Resource{
			"csbpg_binding_user": resourceBindingUser(),
		},
	}
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
	var diags diag.Diagnostics

	factory := connectionFactory{
		host:          d.Get(hostKey).(string),
		port:          d.Get(portKey).(int),
		username:      d.Get(usernameKey).(string),
		password:      d.Get(passwordKey).(string),
		database:      d.Get(databaseKey).(string),
		dataOwnerRole: d.Get(dataOwnerRoleKey).(string),
		sslMode:       d.Get(sslModeKey).(string),
		sslRootCert:   d.Get(sslRootCertKey).(string),
	}

	if value, ok := d.GetOk(clientCertKey); ok {
		if spec, ok := value.([]any)[0].(map[string]any); ok {
			factory.sslClientCert = &clientCertificateConfig{
				Certificate: spec["cert"].(string),
				Key:         spec["key"].(string),
			}
		}
	}

	return factory, diags
}
