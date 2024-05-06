package mysqlqgen

import (
	"fmt"
	"time"

	"github.com/lib/pq"
	"golang.org/x/exp/slog"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	sqlxtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/jmoiron/sqlx"
)

func Init(host string, port string, uname string, pwd string, dbname string, dbparams string, servicename string) *sqlx.DB {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Errors")
			fmt.Println("Recovered from panic:", r)
		}
	}()

	// Register PostgreSQL driver for tracing
	sqltrace.Register("postgres", &pq.Driver{}, sqltrace.WithServiceName(servicename))

	// Construct the connection string
	dsnString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s", uname, pwd, host, port, dbname, dbparams)

	// Connect to PostgreSQL database with tracing
	db, err := sqlxtrace.Connect("postgres", dsnString)
	if err != nil {
		msg := fmt.Sprintf("Cannot connect to PostgreSQL: %s, %v", dsnString, err)
		slog.Error(msg)
		panic(msg)
	}

	// Set database connection pool settings
	db.SetMaxOpenConns(300)
	db.SetMaxIdleConns(300)
	db.SetConnMaxLifetime(3 * time.Minute)

	return db
}
