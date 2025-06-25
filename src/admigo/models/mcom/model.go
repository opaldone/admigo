package mcom

import (
	"database/sql"
	"fmt"

	"admigo/applog"
	"admigo/config"

	_ "github.com/lib/pq"
)

var Db *sql.DB

func LoadDriver() {
	var err error

	c := config.Env(false)

	connect := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Db.Host,
		c.Db.User,
		c.Db.Password,
		c.Db.Dbname,
		c.Db.Sslmode,
	)

	Db, err = sql.Open(c.Db.Driver, connect)
	if err != nil {
		applog.Danger("Error in loadDriver", err)
	}
}
