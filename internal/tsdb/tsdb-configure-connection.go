// Copyright 2019 Cray Inc. All Rights Reserved.
//
// Except as permitted by contract or express written permission of Cray Inc.,
// no part of this work or its content may be modified, used, reproduced or
// disclosed in any form. Modifications made without express permission of
// Cray Inc. may damage the system the software is installed within, may
// disqualify the user from receiving support from Cray Inc. under support or
// maintenance contracts, or require additional support services outside the
// scope of those contracts to repair the software or system.

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

