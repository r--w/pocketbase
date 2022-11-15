package pocketbase

import (
	"testing"

	"github.com/r--w/pocketbase/migrations"
	"github.com/stretchr/testify/assert"
)

const (
	defaultURL = "http://127.0.0.1:8090"
)

func TestAuthorizeAnonymous(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Empty credentials",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient(defaultURL)
			err := c.Authorize()
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestListAccess(t *testing.T) {
	type auth struct {
		email    string
		password string
	}
	tests := []struct {
		name       string
		admin      auth
		user       auth
		collection string
		wantResult bool
		wantErr    bool
	}{
		{
			name:       "With admin credentials - posts_admin",
			admin:      auth{email: migrations.AdminEmailPassword, password: migrations.AdminEmailPassword},
			collection: migrations.PostsAdmin,
			wantResult: true,
			wantErr:    false,
		},
		{
			name:       "Without credentials - posts_admin",
			collection: migrations.PostsAdmin,
			wantErr:    true,
		},
		{
			name:       "Without credentials - posts_public",
			collection: migrations.PostsPublic,
			wantResult: true,
			wantErr:    false,
		},
		{
			// For access rule @request.auth.id != ""
			// no error is returned, but empty result
			name:       "Without credentials - posts_user",
			collection: migrations.PostsUser,
			wantResult: false,
			wantErr:    false,
		},
		{
			name:       "With user credentials - posts_user",
			user:       auth{email: migrations.UserEmailPassword, password: migrations.UserEmailPassword},
			collection: migrations.PostsUser,
			wantResult: true,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient(defaultURL)
			if tt.admin.email != "" {
				c = NewClient(defaultURL, WithAdminEmailPassword(tt.admin.email, tt.admin.password))
			} else if tt.user.email != "" {
				c = NewClient(defaultURL, WithUserEmailPassword(tt.user.email, tt.user.password))
			}
			r, err := c.List(tt.collection, Params{})
			assert.Equal(t, tt.wantErr, err != nil, err)
			assert.Equal(t, tt.wantResult, r.TotalItems > 0)
		})
	}
}

func TestAuthorizeEmailPassword(t *testing.T) {
	type args struct {
		email    string
		password string
	}
	tests := []struct {
		name    string
		admin   args
		user    args
		wantErr bool
	}{
		{
			name:    "Valid credentials admin",
			admin:   args{email: migrations.AdminEmailPassword, password: migrations.AdminEmailPassword},
			wantErr: false,
		},
		{
			name:    "Invalid credentials admin",
			admin:   args{email: "invalid_" + migrations.AdminEmailPassword, password: "no_admin@admin.com"},
			wantErr: true,
		},
		{
			name:    "Valid credentials user",
			user:    args{email: migrations.UserEmailPassword, password: migrations.UserEmailPassword},
			wantErr: false,
		},
		{
			name:    "Invalid credentials user",
			user:    args{email: "invalid_" + migrations.UserEmailPassword, password: migrations.UserEmailPassword},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient(defaultURL)
			if tt.admin.email != "" {
				c = NewClient(defaultURL, WithAdminEmailPassword(tt.admin.email, tt.admin.password))
			} else if tt.user.email != "" {
				c = NewClient(defaultURL, WithUserEmailPassword(tt.user.email, tt.user.password))
			}
			err := c.Authorize()
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestClient_List(t *testing.T) {
	defaultClient := NewClient(defaultURL)

	tests := []struct {
		name       string
		client     *Client
		collection string
		params     Params
		wantResult bool
		wantErr    bool
	}{
		{
			name:       "List with no params",
			client:     defaultClient,
			collection: migrations.PostsPublic,
			wantErr:    false,
			wantResult: true,
		},
		{
			name:       "List no results - query",
			client:     defaultClient,
			collection: migrations.PostsPublic,
			params: Params{
				Filters: "field='some_random_value'",
			},
			wantErr:    false,
			wantResult: false,
		},
		{
			name:       "List no results - invalid query",
			client:     defaultClient,
			collection: migrations.PostsPublic,
			params: Params{
				Filters: "field~~~some_random_value'",
			},
			wantErr:    true,
			wantResult: false,
		},
		{
			name:       "List invalid collection",
			client:     defaultClient,
			collection: "invalid_collection",
			wantErr:    true,
			wantResult: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.List(tt.collection, tt.params)
			assert.Equal(t, tt.wantErr, err != nil, err)
			assert.Equal(t, tt.wantResult, got.TotalItems > 0)
		})
	}
}
