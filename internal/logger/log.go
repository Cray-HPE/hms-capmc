// Copyright 2019 Cray Inc. All Rights Reserved.
//
// Except as permitted by contract or express written permission of Cray Inc.,
// no part of this work or its content may be modified, used, reproduced or
// disclosed in any form. Modifications made without express permission of
// Cray Inc. may damage the system the software is installed within, may
// disqualify the user from receiving support from Cray Inc. under support or
// maintenance contracts, or require additional support services outside the
// scope of those contracts to repair the software or system.

package logger

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)
//SetupLogging - setups the Global Logrus logging
func SetupLogging() {
	logLevel := os.Getenv("LOG_LEVEL")
	logLevel = strings.ToUpper(logLevel)

	log.SetFormatter(&myFormatter{log.TextFormatter{
		FullTimestamp:          true,
		TimestampFormat:        "2006/01/02 15:04:05",
		ForceColors:            true,
		DisableLevelTruncation: true,
	},
	})
	log.SetReportCaller(true)

	switch logLevel {

	case "TRACE":
		log.SetLevel(log.TraceLevel)
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "WARN":
		log.SetLevel(log.WarnLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	case "FATAL":
		log.SetLevel(log.FatalLevel)
	case "PANIC":
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.ErrorLevel)

	}

	log.WithFields(log.Fields{"LogLevel": logLevel}).Info("Initializing Logging")
}

type myFormatter struct {
	log.TextFormatter
}

//This function is VERY close to the current CAPMC output.  Ive added a data={} on the end if it has elements.
// I have a fencepost problem (dangling ',') but not sure that matters
func (f *myFormatter) Format(entry *log.Entry) ([]byte, error) {

	path := entry.Caller.File
	file := filepath.Base(path)
	var str string
	if len(entry.Data) > 0 {
		str += "Data={"
		for k, v := range entry.Data {
			str += fmt.Sprintf("%s:%s, ", k, v)
		}
		str += "}"
	}
	re := regexp.MustCompile(`\r?\n`)
	str = re.ReplaceAllString(str, " ")
	str = strings.Replace(str, "\t", " ", 0)

	return []byte(fmt.Sprintf("%s %s:%d: %s: %s %s\n", entry.Time.Format(f.TimestampFormat), file, entry.Caller.Line, strings.Title(entry.Level.String()), entry.Message, str)), nil
}
