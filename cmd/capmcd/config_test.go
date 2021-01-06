// Copyright 2019 Cray Inc. All Rights Reserved.
//
// This file contains the unit tests for config
//

package main

import (
	"testing"
)

const configFile string = "testconfigfile.toml"

func TestLoadConfigFile(t *testing.T) {

	var (
		config Config
		ok     bool
	)

	ok = loadConfigFile("", &config)
	if ok {
		t.Error("TestLoadConfigFile Test Case No File: FAIL: Expected false but got true")
	}

	ok = loadConfigFile("no-such-file.toml", &config)
	if ok {
		t.Error("TestLoadConfigFile Test Case Non-existent File: FAIL: Expected false but got true")
	}

	ok = loadConfigFile(configFile, &config)
	if !ok {
		t.Error("TestLoadConfigFile Test Case Config File: FAIL: Expected true but got false")
	}
}

// Yes this is silly as loadConfig can't fail
func TestLoadConfig(t *testing.T) {

	loadConfig("")
	loadConfig("no-such-config.toml")
	loadConfig(configFile)
}
