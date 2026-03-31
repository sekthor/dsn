package dsn

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"strings"
	"text/template"
)

type Config struct {
	User         string
	Password     string
	PasswordFile string
	Database     string
	Host         string
	Port         int
	Options      url.Values
}

// creates a postgresql connection string with format
// postgresql://[userspec@][hostspec][/dbname][?paramspec]
//
// Examples (taken from postgresql docs):
// postgresql://
// postgresql://localhost
// postgresql://localhost:5433
// postgresql://localhost/mydb
// postgresql://user@localhost
// postgresql://user:secret@localhost
// postgresql://other@localhost/otherdb?connect_timeout=10&application_name=myapp
// postgresql://host1:123,host2:456/somedb?target_session_attrs=any&application_name=myapp
func (c Config) Postgresql() string {
	c.Init()
	if c.Host == "" {
		return ""
	}

	dsn := c.Host
	if c.Port != 0 {
		dsn = fmt.Sprintf("%s:%d", dsn, c.Port)
	}

	if c.User != "" {
		user := c.User
		if c.Password != "" {
			user = fmt.Sprintf("%s:%s", user, c.Password)
		}
		dsn = fmt.Sprintf("%s@%s", user, dsn)
	}

	if c.Database != "" {
		dsn = fmt.Sprintf("%s/%s", dsn, c.Database)
	}

	if len(c.Options) > 0 {
		dsn = fmt.Sprintf("%s?%s", dsn, c.Options.Encode())
	}

	return fmt.Sprintf("postgresql://%s", dsn)
}

// creates a postgres connection string of key-value format
func (c Config) PostgresqlKV() string {
	c.Init()
	dsnElements := []string{}
	if c.Host != "" {
		dsnElements = append(dsnElements, fmt.Sprintf("host=%s", c.Host))
	}
	if c.User != "" {
		dsnElements = append(dsnElements, fmt.Sprintf("user=%s", c.User))
	}
	if c.Password != "" {
		dsnElements = append(dsnElements, fmt.Sprintf("password=%s", c.Password))
	}
	if c.Port != 0 {
		dsnElements = append(dsnElements, fmt.Sprintf("port=%d", c.Port))
	}
	if c.Database != "" {
		dsnElements = append(dsnElements, fmt.Sprintf("dbname=%s", c.Database))
	}
	for k, v := range c.Options {
		var value string
		// TODO: don't just ignore subsequent users
		if len(v) > 0 {
			value = v[0]
		}
		dsnElements = append(dsnElements, fmt.Sprintf("%s=%s", k, value))
	}
	return strings.Join(dsnElements, " ")
}

func (c Config) FromTemplate(tmplStr string) (string, error) {
	c.Init()
	var buffer bytes.Buffer
	tmpl, err := template.New("dsn").Parse(tmplStr)
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(&buffer, c)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func (c *Config) Init() {
	if c.Password != "" {
		return
	}

	if c.PasswordFile != "" {
		data, err := os.ReadFile(c.PasswordFile)
		if err != nil {
			return
		}
		c.Password = string(data)
	}
}
