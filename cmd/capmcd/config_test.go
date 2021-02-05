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
