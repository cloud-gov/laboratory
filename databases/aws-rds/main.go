package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "gopkg.in/goracle.v2"
)

func main() {
	appEnv, err := cfenv.Current()
	if err != nil {
		panic("not in cloud foundry")
	}

	svcs := appEnv.Services
	dbType := os.Getenv("DB_TYPE")
	svcName := os.Getenv("SERVICE_NAME")
	svc, err := svcs.WithName(svcName)

	switch dbType {
	case "postgres":
		openAndTest("postgres", svc.Credentials["uri"].(string))
	case "mysql":
		openAndTest("mysql", fmtMysql(svc))
	case "oracle":
		openAndTest("goracle", svc.Credentials["uri"].(string))
	}
	fmt.Printf("%#v", svcs)
}

// set up the mysql connection strings.
func fmtMysql(svc *cfenv.Service) string {
	cfg := mysql.NewConfig()
	var ok bool

	if cfg.User, ok = svc.CredentialString("username"); !ok {
		panic("cannot parse username in mysql config")
	}
	if cfg.Passwd, ok = svc.CredentialString("password"); !ok {
		panic("cannot parse password in mysql config")
	}
	if cfg.DBName, ok = svc.CredentialString("db_name"); !ok {
		panic("cannot parse db_name in mysql config")
	}

	cfg.Net = "tcp"
	cfg.Addr = fmt.Sprintf("%s", svc.Credentials["host"].(string))

	return cfg.FormatDSN()
}

func openAndTest(dbType string, dsn interface{}) {
	db, err := sql.Open(dbType, fmt.Sprintf("%v", dsn))
	if err != nil {
		panic(err)
	}

	// create the table.
	if _, err := db.Exec("create table smoke (id integer, name text)"); err != nil {
		panic(err)
	}

	// insert into the table.
	if _, err := db.Exec("insert into smoke values (?, 'smoke')", 1); err != nil {
		panic(err)
	}

	f := os.Getenv("ENABLE_FUNCTIONS")
	yes, err := strconv.ParseBool(f)
	if yes {
		// test a function.
		if _, err := db.Exec("create function hello(id INT) returns CHAR(50) return 'foobar'", 1); err != nil {
			panic(err)
		}
	}

	// cleanup
	if _, err := db.Exec("drop table smoke", 1); err != nil {
		panic(err)
	}
}
