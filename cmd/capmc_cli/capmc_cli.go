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
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"stash.us.cray.com/HMS/hms-capmc/internal/capmc"
	"strings"
)

type osConfig struct {
	serviceUrl string `json:"os_service_url"`
	cert       string `json:"os_cert"`
	key        string `json:"os_key"`
	cacert     string `json:"os_cacert"`
}

var (
	configFile = "/etc/opt/cray/capmc/capmc.json"
	// TODO: wrong defaults
	osconfig = osConfig{
		serviceUrl: "https://localhost:8443",
		cert:       "capmc-client.pem",
		key:        "capmc-client.key",
		cacert:     "capmc-cacert.pem",
	}
	serviceName string
)

// Util functions
//-------------------------------------------

func getConfig() {
	raw, err := ioutil.ReadFile(configFile)
	if err == nil {
		err = json.Unmarshal(raw, &osconfig)
		if err != nil {
			panic(err) // TODO
		}
	}

	if val := os.Getenv("OS_SERVICE_URL"); val != "" {
		osconfig.serviceUrl = val
	}
	if val := os.Getenv("OS_CERT"); val != "" {
		osconfig.cert = val
	}
	if val := os.Getenv("OS_KEY"); val != "" {
		osconfig.key = val
	}
	if val := os.Getenv("OS_CACERT"); val != "" {
		osconfig.cacert = val
	}
}

// API functions
//-------------------------------------------

func getPowerCap(args []string) {
	flagset := flag.FlagSet{}
	// TODO: short version of these opts also needed
	nidsFlag := flagset.String("nids", "", "TODO")
	flagset.Parse(args)

	nidlist, err1 := capmc.ParseNidlist(*nidsFlag)
	if err1 != nil {
		panic(err1)
	}
	request := capmc.NidlistRequest{Nids: nidlist}

	data, err2 := submitRequest("get_power_cap", request)
	if err2 != nil {
		panic(err2)
	}
	fmt.Print(data)
}

func getPowerCapCapabilities(args []string) {
	flagset := flag.FlagSet{}
	// TODO: short version of these opts also needed
	nidsFlag := flagset.String("nids", "", "TODO")
	flagset.Parse(args)

	nidlist, err1 := capmc.ParseNidlist(*nidsFlag)
	if err1 != nil {
		panic(err1)
	}
	request := capmc.NidlistRequest{Nids: nidlist}

	data, err2 := submitRequest("get_power_cap_capabilities", request)
	if err2 != nil {
		panic(err2)
	}
	fmt.Print(data)
}

func nodeRules(args []string) {
	data, err := submitRequest("get_node_rules", struct{}{})
	if err != nil {
		panic(err) // TODO
	}
	fmt.Print(data)
}

func nodeStatus(args []string) {
	flagset := flag.FlagSet{}
	// TODO: short version of these opts also needed
	nidsFlag := flagset.String("nids", "", "TODO")
	filterFlag := flagset.String("filter", "", "TODO")
	flagset.Parse(args)

	nidlist, err1 := capmc.ParseNidlist(*nidsFlag)
	if err1 != nil {
		panic(err1)
	}

	request := capmc.NodeStatusRequest{
		Filter: *filterFlag,
		Nids:   nidlist}

	data, err2 := submitRequest("get_node_status", request)
	if err2 != nil {
		panic(err2) // TODO
	}
	fmt.Print(data)
}

// TODO: interface{} ?
func submitRequest(path string, request interface{}) (string, error) {
	b, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}
	return submitRequestRaw(path, b)
}

func submitRequestRaw(path string, requestJson []byte) (string, error) {
	cert, err := tls.LoadX509KeyPair(osconfig.cert, osconfig.key)
	if err != nil {
		panic(err)
	}
	caCert, err := ioutil.ReadFile(osconfig.cacert)
	if err != nil {
		panic(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true, // TODO: remove
	}

	tr := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: tr}
	req,rerr := http.NewRequest("POST",osconfig.serviceUrl+"/capmc/"+path,bytes.NewBuffer(requestJson))
	if (rerr != nil) {
		panic(rerr)
	}
	base.SetHTTPUserAgent(req,serviceName)
	req.Header.Add("Content-Type","text/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	// TODO: close body?
	// TODO: better error handling, and handling of I/O etc...
	if resp.StatusCode == http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return bodyString, nil
	}

	fmt.Println(resp)
	return "", errors.New("TODO: something went wrong")
}

func jsonInput(args []string) {
	flagset := flag.FlagSet{}
	// TODO: short version of these flags
	resourceFlag := flagset.String("resource", "", "TODO")
	flagset.Parse(args)

	if *resourceFlag == "" {
		panic("--resource must be specified")
	}

	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	var response string
	response, err = submitRequestRaw(*resourceFlag, input)
	if err != nil {
		panic(err)
	}
	fmt.Print(response)
}

func main() {
	getConfig()

	appletMap := map[string]func([]string){
		"get_power_cap":              getPowerCap,
		"get_power_cap_capabilities": getPowerCapCapabilities,
		"node_rules":                 nodeRules,
		"node_status":                nodeStatus,
		"json":                       jsonInput,
	}

	if len(os.Args) <= 1 {
		fmt.Println("TODO: show help")
		os.Exit(1)
	}
	serviceName,err := base.GetServiceInstanceName()
	if (err != nil) {
		serviceName = "CAPMC_CLI"
	}
	applet, ok := appletMap[os.Args[1]]
	if !ok {
		fmt.Fprintln(os.Stderr, "Error: Unknown applet:", os.Args[1])
		fmt.Fprintln(os.Stderr, "The available applets are:")
		keys := []string{}
		for key, _ := range appletMap {
			keys = append(keys, key)
		}
		fmt.Fprintln(os.Stderr, strings.Join(keys, ","))
		fmt.Fprintln(os.Stderr, "Type'", os.Args[0], "help' for usage information")
		os.Exit(130)
	}

	applet(os.Args[2:])
}
