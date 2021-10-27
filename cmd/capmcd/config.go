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

package main

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

//type PowerCtlMap map[string]PowerCtl

// Config holds the node rules, power controls, and system parameter information
// read from the config.toml file.
type Config struct {
	NodeRules     PowerOpRules        `toml:"NodeRules"`
	PowerControls map[string]PowerCtl `toml:"PowerControls"`
	SystemParams  SystemParameters    `toml:"SystemParameters"`
	CapmcConf     CapmcConfiguration  `toml:"CapmcConfiguration"`
}

// PowerCtl holds the list of blocked roles, component sequences, and reset
// types
type PowerCtl struct {
	BlockRole []string `toml:"BlockRole"`
	CompSeq   []string `toml:"ComponentSequence"`
	ResetType []string `toml:"ResetType"`
}

var config = Config{
	defaultNodeRules,
	defaultPowerControl,
	defaultSystemParameters,
	defaultCapmcConfiguration,
}

const (
	// ConfigPath is the path for the CAPMC service configuration file(s)
	ConfigPath string = "/usr/local/etc/capmc-service"
	// ConfigFile is the filename for the CAPMC service configuration file
	ConfigFile string = "config.toml"
	// TODO Add a development ConfigPath controlled by an environment
	// variable, allowing for non-install/non-container development.
)

func loadConfigFile(file string, config *Config) bool {

	// Attempt to read from file
	log.Printf("Info: Reading config from: %s\n", file)
	md, err := toml.DecodeFile(file, &config)
	if err != nil {
		switch {
		case os.IsNotExist(err):
			log.Printf("Warning: %s: %s", file, os.ErrNotExist)
		case os.IsPermission(err):
			log.Printf("Warning: %s: %s", file, os.ErrPermission)
		default:
			log.Printf("Warning: %s", err)
		}

		return false
	}

	// Check for any undecoded keys from the configuration file
	undecodedKeys := md.Undecoded()
	if len(undecodedKeys) > 0 {
		log.Printf("Info: %s: unexpected configuration keys", file)
		if svc.debug {
			log.Printf("DEBUG: Undecoded: keys: %q", undecodedKeys)
		}
	}

	return true
}

func loadConfig(configFile string) *Config {

	var ok bool
	if len(configFile) > 0 {
		ok = loadConfigFile(configFile, &config)
		if ok {
			log.Printf("Info: %s: config values override defaults",
				configFile)
		}
	}

	if !ok {
		log.Printf("Info: using internal default config values")
	}

	return &config
}
