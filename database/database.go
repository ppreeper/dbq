package database

import (
	"fmt"

	_ "github.com/denisenkom/go-mssqldb" //mssql driver
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //postgresql driver
)

// Database struct contains sql pointer
type Database struct {
	Name     string `json:"name"`
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	uri      string
	*sqlx.DB
}

// OpenDatabase open database
func OpenDatabase(db Database) (*Database, error) {
	var err error
	db.GenURI()
	db.DB, err = sqlx.Open(db.Driver, db.uri)
	if err != nil {
		fmt.Printf("Open sql (%v): %v", db.uri, err)
		panic(err)
	}
	if err = db.Ping(); err != nil {
		fmt.Printf("Ping sql: %v", err)
		panic(err)
	}
	return &db, err
}

// GenURI generate db uri string
func (db *Database) GenURI() {
	// fmt.Println(db.Driver)
	if db.Driver == "postgres" {
		if db.Port == "" {
			db.uri = "postgres://" + db.Username + ":" + db.Password + "@" + db.Host + ":5432/" + db.Database + "?sslmode=disable"
		} else {
			db.uri = "postgres://" + db.Username + ":" + db.Password + "@" + db.Host + ":" + db.Port + "/" + db.Database + "?sslmode=disable"
		}
	}
	if db.Driver == "mssql" {
		db.uri = "server=" + db.Host + ";user id=" + db.Username + ";password=" + db.Password + ";database=" + db.Database + ";encrypt=disable;connection timeout=7200;keepAlive=30"
	}
}
