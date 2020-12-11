// Copyright 2017-2020 Hewlett Packard Enterprise Development LP

// Emergency Power Off API Test

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"stash.us.cray.com/HMS/hms-capmc/internal/capmc"
	compcreds "stash.us.cray.com/HMS/hms-compcredentials"
	sstorage "stash.us.cray.com/HMS/hms-securestorage"
)

func TestCapmcDvalidateXnamesForEPO(t *testing.T) {
	var svc CapmcD
	expectedError := errors.New("Invalid component")
	type args struct {
		xnames []string
	}
	tests := []struct {
		name   string
		d      *CapmcD
		xnames []string
		err    error
	}{
		{"s0", &svc, []string{"s0"}, nil},
		{"all", &svc, []string{"all"}, nil},
		{"cabinet", &svc, []string{"x0"}, nil},
		{"chassis", &svc, []string{"x0c0"}, nil},
		{"cabinetpdu", &svc, []string{"x0m0p0"}, nil},
		{"bad", &svc, []string{"bad"}, expectedError},
		{"routermodule", &svc, []string{"x0c0r0"}, expectedError},
		{"hsnboard", &svc, []string{"x0c0r0e0"}, expectedError},
		{"computemodule", &svc, []string{"x0c0s0"}, expectedError},
		{"node", &svc, []string{"x0c0s0b0n0"}, expectedError},
		{"cabinetbmc", &svc, []string{"x0b0"}, expectedError},
		{"cabinetpducontroller", &svc, []string{"x0m0"}, expectedError},
		{"cabinetpdupowerconnector", &svc, []string{"x0m0p0v0"}, expectedError},
		{"cabinetpduoutlet", &svc, []string{"x0m0p0j0"}, expectedError},
		{"chassisbmc", &svc, []string{"x0c0b0"}, expectedError},
		{"nodebmc", &svc, []string{"x0c0s0b0"}, expectedError},
		{"routerbmc", &svc, []string{"x0c0r0b0"}, expectedError},
		{"partition", &svc, []string{"p0.0"}, expectedError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasError := false
			err := tt.d.validateXnamesForEPO(tt.xnames)
			if err == nil {
				if tt.err != nil {
					hasError = true
				}
			} else {
				if tt.err == nil || err.Error() != tt.err.Error() {
					hasError = true
				}
			}
			if hasError {
				t.Errorf("validateXnamesForEPO() error = %v expected %v", err, tt.err)
			}
		})
	}
}

const testQueryMountainComponentState = `
{
  "Components": [
    {
      "ID": "x1000c3",
      "Type": "Chassis",
      "State": "On",
      "Flag": "OK",
      "Enabled": true,
      "NetType": "Sling",
      "Arch": "X86"
    },
    {
      "ID": "x1000c5",
      "Type": "Chassis",
      "State": "On",
      "Flag": "OK",
      "Enabled": true,
      "NetType": "Sling",
      "Arch": "X86"
    },
    {
      "ID": "x1000c0",
      "Type": "Chassis",
      "State": "On",
      "Flag": "OK",
      "Enabled": true,
      "NetType": "Sling",
      "Arch": "X86"
    },
    {
      "ID": "x1000c6",
      "Type": "Chassis",
      "State": "On",
      "Flag": "OK",
      "Enabled": true,
      "NetType": "Sling",
      "Arch": "X86"
    },
    {
      "ID": "x1000c1",
      "Type": "Chassis",
      "State": "On",
      "Flag": "OK",
      "Enabled": true,
      "NetType": "Sling",
      "Arch": "X86"
    },
    {
      "ID": "x1000c4",
      "Type": "Chassis",
      "State": "On",
      "Flag": "OK",
      "Enabled": true,
      "NetType": "Sling",
      "Arch": "X86"
    },
    {
      "ID": "x1000c7",
      "Type": "Chassis",
      "State": "On",
      "Flag": "OK",
      "Enabled": true,
      "NetType": "Sling",
      "Arch": "X86"
    },
    {
      "ID": "x1000c2",
      "Type": "Chassis",
      "State": "On",
      "Flag": "OK",
      "Enabled": true,
      "NetType": "Sling",
      "Arch": "X86"
    }
  ]
}
`
const testQueryMountainComponentEndpoints = `
{
  "ComponentEndpoints": [
    {
      "ID": "x1000c3",
      "Type": "Chassis",
      "RedfishType": "Chassis",
      "RedfishSubtype": "Enclosure",
      "OdataID": "/redfish/v1/Chassis/Enclosure",
      "RedfishEndpointID": "x1000c3b0",
      "RedfishEndpointFQDN": "x1000c3b0",
      "RedfishURL": "x1000c3b0/redfish/v1/Chassis/Enclosure",
      "ComponentEndpointType": "ComponentEndpointChassis",
      "RedfishChassisInfo": {
        "Name": "Enclosure",
        "Actions": {
          "#Chassis.Reset": {
            "ResetType@Redfish.AllowableValues": [
              "ForceOff",
              "On",
              "Off"
            ],
            "@Redfish.ActionInfo": "",
            "target": "/redfish/v1/Chassis/Enclosure/Actions/Chassis.Reset"
		  },
		  "Oem": {
			"#Chassis.EmergencyPower": {
			  "ResetType@Redfish.AllowableValues": [
				"ForceOff"
			  ],
			  "target": "/redfish/v1/Chassis/Enclosure/Actions/Oem/Chassis.EmergencyPower",
			  "@Redfish.ActionInfo": ""
			}
		  }
        }
      }
    },
    {
      "ID": "x1000c5",
      "Type": "Chassis",
      "RedfishType": "Chassis",
      "RedfishSubtype": "Enclosure",
      "OdataID": "/redfish/v1/Chassis/Enclosure",
      "RedfishEndpointID": "x1000c5b0",
      "RedfishEndpointFQDN": "x1000c5b0",
      "RedfishURL": "x1000c5b0/redfish/v1/Chassis/Enclosure",
      "ComponentEndpointType": "ComponentEndpointChassis",
      "RedfishChassisInfo": {
        "Name": "Enclosure",
        "Actions": {
          "#Chassis.Reset": {
            "ResetType@Redfish.AllowableValues": [
              "On",
              "ForceOff",
              "Off"
            ],
            "@Redfish.ActionInfo": "",
            "target": "/redfish/v1/Chassis/Enclosure/Actions/Chassis.Reset"
		  },
		  "Oem": {
			"#Chassis.EmergencyPower": {
			  "ResetType@Redfish.AllowableValues": [
				"ForceOff"
			  ],
			  "target": "/redfish/v1/Chassis/Enclosure/Actions/Oem/Chassis.EmergencyPower",
			  "@Redfish.ActionInfo": ""
			}
		  }
        }
      }
    },
    {
      "ID": "x1000c7",
      "Type": "Chassis",
      "RedfishType": "Chassis",
      "RedfishSubtype": "Enclosure",
      "OdataID": "/redfish/v1/Chassis/Enclosure",
      "RedfishEndpointID": "x1000c7b0",
      "RedfishEndpointFQDN": "x1000c7b0",
      "RedfishURL": "x1000c7b0/redfish/v1/Chassis/Enclosure",
      "ComponentEndpointType": "ComponentEndpointChassis",
      "RedfishChassisInfo": {
        "Name": "Enclosure",
        "Actions": {
          "#Chassis.Reset": {
            "ResetType@Redfish.AllowableValues": [
              "On",
              "ForceOff",
              "Off"
            ],
            "@Redfish.ActionInfo": "",
            "target": "/redfish/v1/Chassis/Enclosure/Actions/Chassis.Reset"
          },
		  "Oem": {
			"#Chassis.EmergencyPower": {
			  "ResetType@Redfish.AllowableValues": [
				"ForceOff"
			  ],
			  "target": "/redfish/v1/Chassis/Enclosure/Actions/Oem/Chassis.EmergencyPower",
			  "@Redfish.ActionInfo": ""
			}
		  }
		}
      }
    },
    {
      "ID": "x1000c2",
      "Type": "Chassis",
      "RedfishType": "Chassis",
      "RedfishSubtype": "Enclosure",
      "OdataID": "/redfish/v1/Chassis/Enclosure",
      "RedfishEndpointID": "x1000c2b0",
      "RedfishEndpointFQDN": "x1000c2b0",
      "RedfishURL": "x1000c2b0/redfish/v1/Chassis/Enclosure",
      "ComponentEndpointType": "ComponentEndpointChassis",
      "RedfishChassisInfo": {
        "Name": "Enclosure",
        "Actions": {
          "#Chassis.Reset": {
            "ResetType@Redfish.AllowableValues": [
              "ForceOff",
              "Off",
              "On"
            ],
            "@Redfish.ActionInfo": "",
            "target": "/redfish/v1/Chassis/Enclosure/Actions/Chassis.Reset"
		  },
		  "Oem": {
			"#Chassis.EmergencyPower": {
			  "ResetType@Redfish.AllowableValues": [
				"ForceOff"
			  ],
			  "target": "/redfish/v1/Chassis/Enclosure/Actions/Oem/Chassis.EmergencyPower",
			  "@Redfish.ActionInfo": ""
			}
		  }
        }
      }
    },
    {
      "ID": "x1000c0",
      "Type": "Chassis",
      "RedfishType": "Chassis",
      "RedfishSubtype": "Enclosure",
      "OdataID": "/redfish/v1/Chassis/Enclosure",
      "RedfishEndpointID": "x1000c0b0",
      "RedfishEndpointFQDN": "x1000c0b0",
      "RedfishURL": "x1000c0b0/redfish/v1/Chassis/Enclosure",
      "ComponentEndpointType": "ComponentEndpointChassis",
      "RedfishChassisInfo": {
        "Name": "Enclosure",
        "Actions": {
          "#Chassis.Reset": {
            "ResetType@Redfish.AllowableValues": [
              "On",
              "ForceOff",
              "Off"
            ],
            "@Redfish.ActionInfo": "",
            "target": "/redfish/v1/Chassis/Enclosure/Actions/Chassis.Reset"
		  },
		  "Oem": {
			"#Chassis.EmergencyPower": {
			  "ResetType@Redfish.AllowableValues": [
				"ForceOff"
			  ],
			  "target": "/redfish/v1/Chassis/Enclosure/Actions/Oem/Chassis.EmergencyPower",
			  "@Redfish.ActionInfo": ""
			}
		  }
        }
      }
    },
    {
      "ID": "x1000c6",
      "Type": "Chassis",
      "RedfishType": "Chassis",
      "RedfishSubtype": "Enclosure",
      "OdataID": "/redfish/v1/Chassis/Enclosure",
      "RedfishEndpointID": "x1000c6b0",
      "RedfishEndpointFQDN": "x1000c6b0",
      "RedfishURL": "x1000c6b0/redfish/v1/Chassis/Enclosure",
      "ComponentEndpointType": "ComponentEndpointChassis",
      "RedfishChassisInfo": {
        "Name": "Enclosure",
        "Actions": {
          "#Chassis.Reset": {
            "ResetType@Redfish.AllowableValues": [
              "ForceOff",
              "Off",
              "On"
            ],
            "@Redfish.ActionInfo": "",
            "target": "/redfish/v1/Chassis/Enclosure/Actions/Chassis.Reset"
		  },
		  "Oem": {
			"#Chassis.EmergencyPower": {
			  "ResetType@Redfish.AllowableValues": [
				"ForceOff"
			  ],
			  "target": "/redfish/v1/Chassis/Enclosure/Actions/Oem/Chassis.EmergencyPower",
			  "@Redfish.ActionInfo": ""
			}
		  }
        }
      }
    },
    {
      "ID": "x1000c1",
      "Type": "Chassis",
      "RedfishType": "Chassis",
      "RedfishSubtype": "Enclosure",
      "OdataID": "/redfish/v1/Chassis/Enclosure",
      "RedfishEndpointID": "x1000c1b0",
      "RedfishEndpointFQDN": "x1000c1b0",
      "RedfishURL": "x1000c1b0/redfish/v1/Chassis/Enclosure",
      "ComponentEndpointType": "ComponentEndpointChassis",
      "RedfishChassisInfo": {
        "Name": "Enclosure",
        "Actions": {
          "#Chassis.Reset": {
            "ResetType@Redfish.AllowableValues": [
              "On",
              "ForceOff",
              "Off"
            ],
            "@Redfish.ActionInfo": "",
            "target": "/redfish/v1/Chassis/Enclosure/Actions/Chassis.Reset"
		  },
		  "Oem": {
			"#Chassis.EmergencyPower": {
			  "ResetType@Redfish.AllowableValues": [
				"ForceOff"
			  ],
			  "target": "/redfish/v1/Chassis/Enclosure/Actions/Oem/Chassis.EmergencyPower",
			  "@Redfish.ActionInfo": ""
			}
		  }
        }
      }
    },
    {
      "ID": "x1000c4",
      "Type": "Chassis",
      "RedfishType": "Chassis",
      "RedfishSubtype": "Enclosure",
      "OdataID": "/redfish/v1/Chassis/Enclosure",
      "RedfishEndpointID": "x1000c4b0",
      "RedfishEndpointFQDN": "x1000c4b0",
      "RedfishURL": "x1000c4b0/redfish/v1/Chassis/Enclosure",
      "ComponentEndpointType": "ComponentEndpointChassis",
      "RedfishChassisInfo": {
        "Name": "Enclosure",
        "Actions": {
          "#Chassis.Reset": {
            "ResetType@Redfish.AllowableValues": [
              "ForceOff",
              "Off",
              "On"
            ],
            "@Redfish.ActionInfo": "",
            "target": "/redfish/v1/Chassis/Enclosure/Actions/Chassis.Reset"
		  },
		  "Oem": {
			"#Chassis.EmergencyPower": {
			  "ResetType@Redfish.AllowableValues": [
				"ForceOff"
			  ],
			  "target": "/redfish/v1/Chassis/Enclosure/Actions/Oem/Chassis.EmergencyPower",
			  "@Redfish.ActionInfo": ""
			}
		  }
        }
      }
    }
  ]
}
`

const testQueryMountainRedfishEndpoints = `
{
  "RedfishEndpoints": [
    {
      "ID": "x1000c0b0",
      "Type": "ChassisBMC",
      "Hostname": "",
      "Domain": "",
      "FQDN": "x1000c0b0",
      "Enabled": true,
      "User": "root",
      "Password": "********",
      "RediscoverOnUpdate": false,
      "DiscoveryInfo": {
        "LastDiscoveryAttempt": "2019-08-17T00:37:45.038887Z",
        "LastDiscoveryStatus": "DiscoverOK",
        "RedfishVersion": "1.2.0"
      }
    },
    {
      "ID": "x1000c4b0",
      "Type": "ChassisBMC",
      "Hostname": "",
      "Domain": "",
      "FQDN": "x1000c4b0",
      "Enabled": true,
      "User": "root",
      "Password": "********",
      "RediscoverOnUpdate": false,
      "DiscoveryInfo": {
        "LastDiscoveryAttempt": "2019-08-17T00:37:45.041232Z",
        "LastDiscoveryStatus": "DiscoverOK",
        "RedfishVersion": "1.2.0"
      }
    },
    {
      "ID": "x1000c7b0",
      "Type": "ChassisBMC",
      "Hostname": "",
      "Domain": "",
      "FQDN": "x1000c7b0",
      "Enabled": true,
      "User": "root",
      "Password": "********",
      "RediscoverOnUpdate": false,
      "DiscoveryInfo": {
        "LastDiscoveryAttempt": "2019-08-17T00:37:45.005758Z",
        "LastDiscoveryStatus": "DiscoverOK",
        "RedfishVersion": "1.2.0"
      }
    },
    {
      "ID": "x1000c1b0",
      "Type": "ChassisBMC",
      "Hostname": "",
      "Domain": "",
      "FQDN": "x1000c1b0",
      "Enabled": true,
      "User": "root",
      "Password": "********",
      "RediscoverOnUpdate": false,
      "DiscoveryInfo": {
        "LastDiscoveryAttempt": "2019-08-17T00:37:45.039136Z",
        "LastDiscoveryStatus": "DiscoverOK",
        "RedfishVersion": "1.2.0"
      }
    },
    {
      "ID": "x1000c2b0",
      "Type": "ChassisBMC",
      "Hostname": "",
      "Domain": "",
      "FQDN": "x1000c2b0",
      "Enabled": true,
      "User": "root",
      "Password": "********",
      "RediscoverOnUpdate": false,
      "DiscoveryInfo": {
        "LastDiscoveryAttempt": "2019-08-17T00:37:45.031925Z",
        "LastDiscoveryStatus": "DiscoverOK",
        "RedfishVersion": "1.2.0"
      }
    },
    {
      "ID": "x1000c6b0",
      "Type": "ChassisBMC",
      "Hostname": "",
      "Domain": "",
      "FQDN": "x1000c6b0",
      "Enabled": true,
      "User": "root",
      "Password": "********",
      "RediscoverOnUpdate": false,
      "DiscoveryInfo": {
        "LastDiscoveryAttempt": "2019-08-17T00:37:45.066160Z",
        "LastDiscoveryStatus": "DiscoverOK",
        "RedfishVersion": "1.2.0"
      }
    },
    {
      "ID": "x1000c3b0",
      "Type": "ChassisBMC",
      "Hostname": "",
      "Domain": "",
      "FQDN": "x1000c3b0",
      "Enabled": true,
      "User": "root",
      "Password": "********",
      "RediscoverOnUpdate": false,
      "DiscoveryInfo": {
        "LastDiscoveryAttempt": "2019-08-17T00:37:44.971719Z",
        "LastDiscoveryStatus": "DiscoverOK",
        "RedfishVersion": "1.2.0"
      }
    },
    {
      "ID": "x1000c5b0",
      "Type": "ChassisBMC",
      "Hostname": "",
      "Domain": "",
      "FQDN": "x1000c5b0",
      "Enabled": true,
      "User": "root",
      "Password": "********",
      "RediscoverOnUpdate": false,
      "DiscoveryInfo": {
        "LastDiscoveryAttempt": "2019-08-17T00:37:45.006564Z",
        "LastDiscoveryStatus": "DiscoverOK",
        "RedfishVersion": "1.2.0"
      }
    }
  ]
}
`

func EPORTTestFunc() RoundTripFunc {
	return func(req *http.Request) (*http.Response, error) {
		switch req.URL.String() {
		case "http://localhost:27779/State/Components/Query/x1000?type=chassis",
			"http://localhost:27779/State/Components/Query/s0?type=chassis",
			"http://localhost:27779/State/Components/Query/all?type=chassis",
			"http://localhost:27779/State/Components?id=x1000c3&id=x1000c5&id=x1000c0&id=x1000c6&id=x1000c1&id=x1000c4&id=x1000c7&id=x1000c2":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(testQueryMountainComponentState)),
				Header:     make(http.Header),
			}, nil
		case "http://localhost:27779/Inventory/ComponentEndpoints?id=x1000c3&id=x1000c5&id=x1000c0&id=x1000c6&id=x1000c1&id=x1000c4&id=x1000c7&id=x1000c2":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(testQueryMountainComponentEndpoints)),
				Header:     make(http.Header),
			}, nil
		case "http://localhost:27779/Inventory/RedfishEndpoints?id=x1000c3b0&id=x1000c5b0&id=x1000c7b0&id=x1000c2b0&id=x1000c0b0&id=x1000c6b0&id=x1000c1b0&id=x1000c4b0":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(testQueryMountainRedfishEndpoints)),
				Header:     make(http.Header),
			}, nil
		case "https://x1000c0b0/redfish/v1/Chassis/Enclosure/Actions/Oem/Chassis.EmergencyPower",
			"https://x1000c1b0/redfish/v1/Chassis/Enclosure/Actions/Oem/Chassis.EmergencyPower",
			"https://x1000c2b0/redfish/v1/Chassis/Enclosure/Actions/Oem/Chassis.EmergencyPower",
			"https://x1000c3b0/redfish/v1/Chassis/Enclosure/Actions/Oem/Chassis.EmergencyPower",
			"https://x1000c4b0/redfish/v1/Chassis/Enclosure/Actions/Oem/Chassis.EmergencyPower",
			"https://x1000c5b0/redfish/v1/Chassis/Enclosure/Actions/Oem/Chassis.EmergencyPower",
			"https://x1000c6b0/redfish/v1/Chassis/Enclosure/Actions/Oem/Chassis.EmergencyPower",
			"https://x1000c7b0/redfish/v1/Chassis/Enclosure/Actions/Oem/Chassis.EmergencyPower":
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(testQueryMountainRedfishEndpoints)),
				Header:     make(http.Header),
			}, nil
		default:
			return &http.Response{
				StatusCode: 404,
				Body:       ioutil.NopCloser(bytes.NewBufferString("")),
				Header:     make(http.Header),
			}, nil
		}
	}
}
func TestCapmcDgenEPOList(t *testing.T) {
	var tSvc CapmcD
	var err error
	tSvc.hsmURL, err = url.Parse("http://localhost:27779")
	testClient := NewTestClient(EPORTTestFunc())
	tSvc.rfClient = testClient
	tSvc.smClient = testClient
	if err != nil {
		t.Fatal(err)
	}

	wantArray := []string{
		"x1000c3", "x1000c5", "x1000c0", "x1000c6",
		"x1000c1", "x1000c4", "x1000c7", "x1000c2",
	}
	tests := []struct {
		name    string
		xnames  []string
		want    []string
		wantErr bool
	}{
		{
			"valid system",
			[]string{"s0"},
			wantArray,
			false,
		},
		{
			"valid all",
			[]string{"all"},
			wantArray,
			false,
		},
		{
			"invalid cabinet",
			[]string{"x2000"},
			nil,
			true,
		},
		{
			"valid cabinet",
			[]string{"x1000"},
			wantArray,
			false,
		},
		{
			"valid chassis",
			[]string{"x1000c0"},
			[]string{"x1000c0"},
			false,
		},
		{
			"valid cabinet+chassis",
			[]string{"x1000", "x2000c0"},
			append(wantArray, "x2000c0"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tSvc.genEPOList(tt.xnames)
			if (err != nil) != tt.wantErr {
				t.Errorf("CapmcD.genEPOList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CapmcD.genEPOList() = %v, want %v", got, tt.want)
			}
		})
	}
}

var ssData = []sstorage.MockLookup{
	{
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x1000c0",
			},
		},
	},
	{
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x1000c1",
			},
		},
	},
	{
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x1000c2",
			},
		},
	},
	{
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x1000c3",
			},
		},
	},
	{
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x1000c4",
			},
		},
	},
	{
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x1000c5",
			},
		},
	},
	{
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x1000c6",
			},
		},
	},
	{
		Output: sstorage.OutputLookup{
			Output: &compcreds.CompCredentials{
				Xname: "x1000c7",
			},
		},
	},
}

func TestCapmcDdoEmergencyPowerOff(t *testing.T) {
	var tSvc CapmcD
	var err error
	tSvc.hsmURL, err = url.Parse("http://localhost:27779")
	if err != nil {
		t.Fatal(err)
	}
	testClient := NewTestClient(EPORTTestFunc())
	tSvc.rfClient = testClient
	tSvc.smClient = testClient
	tSvc.config = loadConfig("")
	ss, adapter := sstorage.NewMockAdapter()
	ccs := compcreds.NewCompCredStore("secret/hms-cred", ss)
	tSvc.ss = ss
	tSvc.ccs = ccs

	handler := http.HandlerFunc(tSvc.doEmergencyPowerOff)

	tests := []struct {
		name     string
		method   string
		body     io.Reader
		code     int
		expected string
	}{
		{
			"Get",
			http.MethodGet,
			nil,
			http.StatusMethodNotAllowed,
			GetNotAllowedJSON + "\n",
		},
		{
			"Put",
			http.MethodPut,
			nil,
			http.StatusMethodNotAllowed,
			PutNotAllowedJSON + "\n",
		},
		{
			"Delete",
			http.MethodDelete,
			nil,
			http.StatusMethodNotAllowed,
			DeleteNotAllowedJSON + "\n",
		},
		{
			"Connect",
			http.MethodConnect,
			nil,
			http.StatusMethodNotAllowed,
			ConnectNotAllowedJSON + "\n",
		},
		{
			"Post empty EOF",
			http.MethodPost,
			bytes.NewBufferString(""),
			http.StatusBadRequest,
			NoRequestJSON + "\n",
		},
		{
			"Post empty body",
			http.MethodPost,
			bytes.NewBuffer(json.RawMessage(`{}`)),
			http.StatusBadRequest,
			MissingXnameJSON + "\n",
		},
		{
			"Post unexpected end of body",
			http.MethodPost,
			bytes.NewBuffer(json.RawMessage(`{"xnames":["x0c0"]`)),
			http.StatusBadRequest,
			UnexpectedEOFJSON + "\n",
		},
		{
			"Post missing start of body",
			http.MethodPost,
			bytes.NewBuffer(json.RawMessage(`"xnames":["x0c0"]}`)),
			http.StatusBadRequest,
			XnameControlUnmarshJSON + "\n",
		},
		{
			"Post valid cabinet",
			http.MethodPost,
			bytes.NewBuffer(json.RawMessage(`{"xnames":["x1000"]}`)),
			http.StatusOK,
			SuccessJSON + "\n",
		},
	}
	adapter.LookupData = ssData
	for _, tt := range tests {
		adapter.LookupNum = 0
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, capmc.EmergencyPowerOff, tt.body)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if tt.code != rr.Code {
				t.Errorf("handler returned wrong status code: want %v but got %v",
					tt.code, rr.Code)
			}

			if tt.expected != rr.Body.String() {
				t.Errorf("handler returned unexpected body: want '%v' but got '%v'",
					tt.expected, rr.Body.String())
			}
		})
	}
}

func TestCapmcDexecuteEmergencyPowerOff(t *testing.T) {
	var tSvc CapmcD
	var err error
	tSvc.hsmURL, err = url.Parse("http://localhost:27779")
	if err != nil {
		t.Fatal(err)
	}
	testClient := NewTestClient(EPORTTestFunc())
	tSvc.rfClient = testClient
	tSvc.smClient = testClient
	tSvc.config = loadConfig("")
	ss, adapter := sstorage.NewMockAdapter()
	ccs := compcreds.NewCompCredStore("secret/hms-cred", ss)
	tSvc.ss = ss
	tSvc.ccs = ccs

	tests := []struct {
		name string
		args capmc.XnameControl
		want capmc.XnameControlResponse
	}{
		{
			"Valid cabinet",
			capmc.XnameControl{
				Xnames: []string{"x1000"},
				Reason: "Test valid cabinet",
			},
			capmc.XnameControlResponse{
				ErrResponse: capmc.ErrResponse{
					E:      0,
					ErrMsg: "",
				},
				Xnames: nil,
			},
		},
		{
			"Valid system",
			capmc.XnameControl{
				Xnames: []string{"s0"},
				Reason: "Test valid system",
			},
			capmc.XnameControlResponse{
				ErrResponse: capmc.ErrResponse{
					E:      0,
					ErrMsg: "",
				},
				Xnames: nil,
			},
		},
		{
			"Valid all",
			capmc.XnameControl{
				Xnames: []string{"s0"},
				Reason: "Test valid all",
			},
			capmc.XnameControlResponse{
				ErrResponse: capmc.ErrResponse{
					E:      0,
					ErrMsg: "",
				},
				Xnames: nil,
			},
		},
		{
			"Invalid xname",
			capmc.XnameControl{
				Xnames: []string{"x1000c0s0"},
				Reason: "Test invalid xname",
			},
			capmc.XnameControlResponse{
				ErrResponse: capmc.ErrResponse{
					E:      400,
					ErrMsg: "failure validating xnames: Invalid component",
				},
				Xnames: nil,
			},
		},
		{
			"Xname doesn't exist",
			capmc.XnameControl{
				Xnames: []string{"x2000"},
				Reason: "Test xname doesn't exist",
			},
			capmc.XnameControlResponse{
				ErrResponse: capmc.ErrResponse{
					E:      400,
					ErrMsg: "failure generating component list: ",
				},
				Xnames: nil,
			},
		},
	}
	adapter.LookupData = ssData
	for _, tt := range tests {
		adapter.LookupNum = 0
		t.Run(tt.name, func(t *testing.T) {
			if got := tSvc.executeEmergencyPowerOff(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CapmcD.executeEmergencyPowerOff() = %v, want %v", got, tt.want)
			}
		})
	}
}
