package dsn_test

import (
	"net/url"
	"testing"

	"github.com/sekthor/dsn"
)

func TestConfig_FromTemplate(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		tmplStr string
		want    string
		wantErr bool
		dc      dsn.Config
	}{
		{
			name: "postgres template",
			dc: dsn.Config{
				Database: "db",
				User:     "user",
				Password: "password",
				Host:     "host",
			},
			tmplStr: "postgresql://{{.User}}:{{.Password}}@{{.Host}}/{{.Database}}",
			wantErr: false,
			want:    "postgresql://user:password@host/db",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := tt.dc.FromTemplate(tt.tmplStr)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("FromTemplate() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("FromTemplate() succeeded unexpectedly")
			}
			if tt.want != got {
				t.Errorf("FromTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_Postgresql(t *testing.T) {
	options := url.Values{}
	options.Add("sslmode", "disable")

	tests := []struct {
		name string // description of this test case
		want string
		conf dsn.Config
	}{
		{
			name: "full",
			want: "postgresql://user:password@host:5432/db?sslmode=disable",
			conf: dsn.Config{
				Password: "password",
				Host:     "host",
				Port:     5432,
				User:     "user",
				Database: "db",
				Options:  options,
			},
		},
		{
			name: "full, no options",
			want: "postgresql://user:password@host:5432/db",
			conf: dsn.Config{
				Password: "password",
				Host:     "host",
				Port:     5432,
				User:     "user",
				Database: "db",
			},
		},
		{
			name: "user, password, host, port, no db",
			want: "postgresql://user:password@host:5432",
			conf: dsn.Config{
				Password: "password",
				Host:     "host",
				Port:     5432,
				User:     "user",
			},
		},
		{
			name: "user, host, port, db, no password",
			want: "postgresql://user@host:5432/db",
			conf: dsn.Config{
				Host:     "host",
				Port:     5432,
				User:     "user",
				Database: "db",
			},
		},
		{
			name: "user, password, host, db, no port",
			want: "postgresql://user:password@host/db",
			conf: dsn.Config{
				Password: "password",
				Host:     "host",
				User:     "user",
				Database: "db",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := tt.conf.Postgresql()
			if got != tt.want {
				t.Errorf("Postresql() = %v, want %v", got, tt.want)
			}
		})
	}
}
