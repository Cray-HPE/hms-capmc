// MIT License
//
// (C) Copyright [2019, 2021] Hewlett Packard Enterprise Development LP
//
// Permission is hereby granted, free of charge, to any person obtaining a
// copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
// THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
// OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.
//

package tsdb

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/lib/pq" //needed for DB stuff
	log "github.com/sirupsen/logrus"
)

// DB - the Database connection
var DB *sql.DB

// mutex - the mutex for creating the singleton
var mutex = &sync.Mutex{}

// DatabaseInfo - the struct used to configure the DB connection
type DatabaseInfo struct {
	Hostname string
	Port     int
	Username string
	Password string
	Database string
}

// NewPostgresDatabase - this is a singleton that will create a connection to the database
func NewPostgresDatabase(database *DatabaseInfo) (err error) {
	mutex.Lock()

	if database.Hostname == "" {
		database.Hostname = "craysma-postgres-cluster.sma.svc.cluster.local"
	}
	if database.Port == 0 {
		database.Port = 5432
	}
	if database.Username == "" {
		database.Username = "pmdbuser"
	}
	if database.Database == "" {
		database.Database = "pmdb"
	}

	connStr := fmt.Sprintf("sslmode=disable user=%s dbname=%s host=%s port=%d", database.Username, database.Database,
		database.Hostname, database.Port)
	if database.Password != "" {
		connStr += " password=" + database.Password
	}
	// log.Info(connStr)
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Error(err)
	}
	log.Info(connStr)

	// This needs to be less than the current global max_connections
	// via a cnf file or the default of 151.  The sql package default is
	// unlimited and this causes the server-side limit to get overwhelmed.
	// This should be kept in sync with the configured value in the
	// postgres docker container.  Note that it is a global value, however,
	// and we should leave a little slack in any case (setting it to 100
	// exacly causes an occasional failure, even with no other processes
	// connecting).
	DB.SetMaxOpenConns(70)

	// Workaround for HMS-1080 (likely applicable in CAPMC too), one of
	// these, so long as a minute is less than wait_timeout.
	DB.SetConnMaxLifetime(time.Minute)

	err = DB.Ping()
	if err != nil {
		log.Error(err)
	}

	mutex.Unlock()
	return err
}
