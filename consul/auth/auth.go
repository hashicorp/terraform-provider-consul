// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package auth

import (
	"errors"
	"fmt"
	"os"
	"sync"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type (
	loginSchemaFunc   func(string) *schema.Schema
	getSchemaResource func(string) *schema.Resource
	validateFunc      func(data *schema.ResourceData, params map[string]interface{}) error
	authLoginFunc     func(*schema.ResourceData) (AuthLogin, error)
)

// authLoginEntry is the tuple of authLoginFunc, schemaFunc.
type authLoginEntry struct {
	field      string
	loginFunc  authLoginFunc
	schemaFunc loginSchemaFunc
}

// AuthLogin returns a new AuthLogin instance from provided schema.ResourceData.
func (a *authLoginEntry) AuthLogin(r *schema.ResourceData) (AuthLogin, error) {
	return a.loginFunc(r)
}

// LoginSchema returns the AuthLogin's schema.Schema.
func (a *authLoginEntry) LoginSchema() *schema.Schema {
	return a.schemaFunc(a.Field())
}

// Field returns the entry's top level schema field name. E.g. auth_login_aws.
func (a *authLoginEntry) Field() string {
	return a.field
}

// authLoginRegistry provides the storage for authLoginEntry, mapped to the
// entry's field name.
type authLoginRegistry struct {
	m sync.Map
}

// Register field for loginFunc and schemaFunc. A field can only be registered
// once.
func (r *authLoginRegistry) Register(field string, loginFunc authLoginFunc, schemaFunc loginSchemaFunc) error {
	e := &authLoginEntry{
		field:      field,
		loginFunc:  loginFunc,
		schemaFunc: schemaFunc,
	}

	_, loaded := r.m.LoadOrStore(field, e)
	if loaded {
		return fmt.Errorf("auth login field %s is already registered", field)
	}
	return nil
}

// Get the authLoginEntry for field.
func (r *authLoginRegistry) Get(field string) (*authLoginEntry, error) {
	v, ok := r.m.Load(field)
	if !ok {
		return nil, fmt.Errorf("auth login function not registered for %s", field)
	}
	if entry, ok := v.(*authLoginEntry); ok {
		return entry, nil
	} else {
		return nil, fmt.Errorf("invalid type %T store in registry", v)
	}
}

// Fields returns the names of all registered AuthLogin's
func (r *authLoginRegistry) Fields() []string {
	var keys []string
	r.m.Range(func(key, _ interface{}) bool {
		keys = append(keys, key.(string))
		return true
	})

	return keys
}

// Values returns a slice of all registered authLoginEntry(s).
func (r *authLoginRegistry) Values() []*authLoginEntry {
	var result []*authLoginEntry
	r.m.Range(func(key, value interface{}) bool {
		result = append(result, value.(*authLoginEntry))
		return true
	})

	return result
}

// AuthLoginFields supported by the provider.
var (
	authLoginInitCheckError = errors.New("auth login not initialized")

	globalAuthLoginRegistry = &authLoginRegistry{}
)

// AuthLogin interface defines the methods that all auth login implementations must provide.
type AuthLogin interface {
	Init(*schema.ResourceData, string) (AuthLogin, error)
	AuthMethodName() string
	Login(*consulapi.Client) (string, error)
	Namespace() (string, bool)
	Partition() (string, bool)
	Params() map[string]interface{}
}

// AuthLoginCommon providing common methods for other AuthLogin* implementations.
type AuthLoginCommon struct {
	authField   string
	params      map[string]interface{}
	initialized bool
}

func (l *AuthLoginCommon) Params() map[string]interface{} {
	return l.params
}

func (l *AuthLoginCommon) Init(d *schema.ResourceData, authField string, validators ...validateFunc) error {
	l.authField = authField
	params, err := l.init(d)
	if err != nil {
		return err
	}

	for _, vf := range validators {
		if err := vf(d, params); err != nil {
			return err
		}
	}

	l.params = params

	return l.validate()
}

func (l *AuthLoginCommon) Namespace() (string, bool) {
	if l.params != nil {
		if ns, ok := l.params["namespace"]; ok && ns.(string) != "" {
			return ns.(string), true
		}
	}
	return "", false
}

func (l *AuthLoginCommon) Partition() (string, bool) {
	if l.params != nil {
		if part, ok := l.params["partition"]; ok && part.(string) != "" {
			return part.(string), true
		}
	}
	return "", false
}

func (l *AuthLoginCommon) AuthMethodName() string {
	return ""
}

func (l *AuthLoginCommon) copyParams(includes ...string) (map[string]interface{}, error) {
	if err := l.validate(); err != nil {
		return nil, err
	}

	params := make(map[string]interface{}, len(l.params))
	if len(includes) == 0 {
		for k, v := range l.params {
			params[k] = v
		}
	} else {
		var missing []string
		for _, k := range includes {
			v, ok := l.params[k]
			if !ok {
				missing = append(missing, k)
				continue
			}
			params[k] = v
		}
		if len(missing) > 0 {
			return nil, fmt.Errorf("missing params %v", missing)
		}
	}

	return params, nil
}

func (l *AuthLoginCommon) copyParamsExcluding(excludes ...string) (map[string]interface{}, error) {
	params, err := l.copyParams()
	if err != nil {
		return nil, err
	}
	for _, k := range excludes {
		delete(params, k)
	}

	return params, nil
}

func (l *AuthLoginCommon) login(client *consulapi.Client, authMethodName string, bearerToken string, meta map[string]string) (string, error) {
	token, _, err := client.ACL().Login(&consulapi.ACLLoginParams{
		AuthMethod:  authMethodName,
		BearerToken: bearerToken,
		Meta:        meta,
	}, nil)
	if err != nil {
		return "", err
	}

	return token.SecretID, nil
}

func (l *AuthLoginCommon) init(d *schema.ResourceData) (map[string]interface{}, error) {
	if l.initialized {
		return nil, fmt.Errorf("auth login already initialized")
	}

	v, ok := d.GetOk(l.authField)
	if !ok {
		return nil, fmt.Errorf("resource data missing field %q", l.authField)
	}

	config := v.([]interface{})
	if len(config) != 1 {
		// this should never happen
		return nil, fmt.Errorf("empty config for %q", l.authField)
	}

	var params map[string]interface{}
	v = config[0]
	if v == nil {
		params = make(map[string]interface{})
	} else {
		params = v.(map[string]interface{})
	}

	l.initialized = true

	return params, nil
}

type authDefault struct {
	field string

	// envVars will override defaultVal.
	// If there are multiple entries in the slice, we use the first value we
	// find that is set in the environment.
	envVars []string
	// defaultVal is the fallback if an env var is not set
	defaultVal string
}

type authDefaults []authDefault

func (l *AuthLoginCommon) setDefaultFields(d *schema.ResourceData, defaults authDefaults, params map[string]interface{}) error {
	for _, f := range defaults {
		if _, ok := l.getOk(d, f.field); !ok {
			// if field is unset in the config, check env
			params[f.field] = f.defaultVal
			for _, envVar := range f.envVars {
				val := os.Getenv(envVar)
				if val != "" {
					params[f.field] = val
					// found a value, no need to check other options
					break
				}
			}
		}
	}

	return nil
}

func (l *AuthLoginCommon) checkRequiredFields(d *schema.ResourceData, params map[string]interface{}, required ...string) error {
	var missing []string
	for _, f := range required {
		if data, ok := l.getOk(d, f); !ok {
			// if the field was unset in the config
			if params[f] == data {
				missing = append(missing, f)
			}
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("required fields are unset: %v", missing)
	}

	return nil
}

func (l *AuthLoginCommon) checkFieldsOneOf(d *schema.ResourceData, fields ...string) error {
	if len(fields) == 0 {
		return nil
	}

	for _, f := range fields {
		if _, ok := l.getOk(d, f); ok {
			return nil
		}
	}

	return fmt.Errorf(
		"at least one field must be set: %v", fields)
}

func (l *AuthLoginCommon) getOk(d *schema.ResourceData, field string) (interface{}, bool) {
	return d.GetOk(l.fieldPath(d, field))
}

func (l *AuthLoginCommon) get(d *schema.ResourceData, field string) interface{} {
	return d.Get(l.fieldPath(d, field))
}

func (l *AuthLoginCommon) fieldPath(_ *schema.ResourceData, field string) string {
	return fmt.Sprintf("%s.0.%s", l.authField, field)
}

func (l *AuthLoginCommon) validate() error {
	if !l.initialized {
		return authLoginInitCheckError
	}

	return nil
}

// GetAuthLogin returns the configured AuthLogin instance from the provider configuration.
func GetAuthLogin(r *schema.ResourceData) (AuthLogin, error) {
	for _, authField := range globalAuthLoginRegistry.Fields() {
		_, ok := r.GetOk(authField)
		if !ok {
			continue
		}

		entry, err := globalAuthLoginRegistry.Get(authField)
		if err != nil {
			return nil, err
		}

		return entry.AuthLogin(r)
	}

	return nil, nil
}

// mustAddLoginSchema adds common login schema fields to the auth resource.
func mustAddLoginSchema(r *schema.Resource, _ string) *schema.Resource {
	m := map[string]*schema.Schema{
		"namespace": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The Consul namespace to authenticate to.",
		},
		"partition": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The Consul admin partition to authenticate to.",
		},
	}

	for k, v := range m {
		if _, ok := r.Schema[k]; ok {
			panic(fmt.Sprintf("cannot add schema field %q, already exists in the Schema map", k))
		}

		r.Schema[k] = v
	}

	return r
}

func getLoginSchema(authField, description string, resourceFunc getSchemaResource) *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeList,
		Optional:      true,
		MaxItems:      1,
		Description:   description,
		Elem:          resourceFunc(authField),
		ConflictsWith: calculateConflictsWith(authField, globalAuthLoginRegistry.Fields()),
	}
}

// calculateConflictsWith returns a list of fields that conflict with the given field.
func calculateConflictsWith(field string, allFields []string) []string {
	var conflicts []string
	for _, f := range allFields {
		if f != field {
			conflicts = append(conflicts, f)
		}
	}
	return conflicts
}

// MustAddAuthLoginSchema adds all supported auth login type schema.Schema to
// a schema map.
func MustAddAuthLoginSchema(s map[string]*schema.Schema) {
	for _, v := range globalAuthLoginRegistry.Values() {
		mustAddSchema(v.Field(), v.LoginSchema(), s)
	}
}

func mustAddSchema(field string, s *schema.Schema, m map[string]*schema.Schema) {
	if _, ok := m[field]; ok {
		panic(fmt.Sprintf("schema field %q already exists", field))
	}
	m[field] = s
}
