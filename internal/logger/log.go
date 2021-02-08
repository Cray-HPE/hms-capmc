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
