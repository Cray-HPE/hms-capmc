/*
 * MIT License
 *
 * (C) Copyright [2019-2021] Hewlett Packard Enterprise Development LP
 *
 * Permission is hereby granted, free of charge, to any person obtaining a
 * copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the
 * Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included
 * in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 * THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
 * OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
 * ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
 * OTHER DEALINGS IN THE SOFTWARE.
 */

package tsdb

import (
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	"sync"
)


// dbmutex - Helps create the singleton
var dbmutex = &sync.Mutex{}

// IsConnected - Protects the ConfigureDataImplementation from switching DB connections
var isConnected bool

//ConnectionType is used to setup the type of connection
type ConnectionType int

//ImplementedConnection is the Type of ConnectionType implemented
var ImplementedConnection ConnectionType

const (
	//UNKNOWN - unknown -> used to tell the system need to pick between one of the following
	UNKNOWN ConnectionType = -1
	//DUMMY - dummy implementation
	DUMMY ConnectionType = 0
	//POSTGRES - the default connection
	POSTGRES ConnectionType = 1
)

// ConfigureDataImplementation - The default will be the DUMMY;
// but the Dockerfile sets the appropriate defaults so that Postrgres is used.
func ConfigureDataImplementation(conType ConnectionType ) (err error) {
	dbmutex.Lock()
	if !isConnected  {
		isConnected = true

		var desiredImplementation string
		if conType == UNKNOWN {
			desiredImplementation = os.Getenv("DATA_IMPLEMENTATION")
			desiredImplementation = strings.ToUpper(desiredImplementation)
		} else if conType == POSTGRES {
			desiredImplementation = "POSTGRES"
		} else {
			desiredImplementation = "DUMMY"
		}

		log.WithField("DATA_IMPLEMENTATION", os.Getenv("DATA_IMPLEMENTATION")).Info("Configuring DATA_IMPLEMENTATION")

		if desiredImplementation  == "DUMMY" {
			ImplementedConnection = DUMMY
			TSDBContext = DummyDB{}
			log.Info("Using InMemoryDB")
		} else {
		// desiredImplementation  == "POSTGRES"
			ImplementedConnection = POSTGRES
			TSDBContext = PostgresqlDB{}
			log.Info("Using PostgresqlDB")

			databaseInfo := DatabaseInfo{}
			databaseInfo.Hostname = os.Getenv("DB_HOSTNAME")
			databaseInfo.Port, _ = strconv.Atoi(os.Getenv("DB_PORT"))
			err = NewPostgresDatabase(&databaseInfo)
			if err!= nil {
				isConnected = false
			} else {
				log.Info("Connected to TSDB database")
			}

		}
	}
	dbmutex.Unlock()
	return err
}

