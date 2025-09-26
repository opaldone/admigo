package mcom

import (
	"database/sql"
	"fmt"

	"admigo/applog"
	"admigo/config"

	_ "github.com/lib/pq"
)

var Dbc *sql.DB

func LoadDriver() {
	var err error

	c := config.Env(false)

	connect := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DB.Host,
		c.DB.User,
		c.DB.Password,
		c.DB.Dbname,
		c.DB.Sslmode,
	)

	Dbc, err = sql.Open(c.DB.Driver, connect)
	if err != nil {
		applog.Danger("Error in loadDriver", err)
	}
}
