package iapi

import (
	"encoding/json"
	"fmt"
)

// GetHost ...
func (server *Server) GetHost(hostname string) ([]HostStruct, error) {

	var hosts []HostStruct

	results, err := server.NewAPIRequest("GET", "/objects/hosts/"+hostname, nil)
	if err != nil {
		return nil, err
	}

	// Contents of the results is an interface object. Need to convert it to json first.
	jsonStr, marshalErr := json.Marshal(results.Results)
	if marshalErr != nil {
		return nil, marshalErr
	}

	// then the JSON can be pushed into the appropriate struct.
	// Note : Results is a slice so much push into a slice.

	if unmarshalErr := json.Unmarshal(jsonStr, &hosts); unmarshalErr != nil {
		return nil, unmarshalErr
	}

	return hosts, err

}

// CreateHost ...
func (server *Server) CreateHost(hostname, address, checkCommand string, variables map[string]interface{}, templates []string, groups []string, zone string) ([]HostStruct, error) {
	attrs := map[string]interface{}{}
	attrs["address"] = address
	attrs["check_command"] = checkCommand

	if variables != nil {
		for key, value := range variables {
			attrs["vars." + key] = value
		}
	}
	if groups != nil {
	  attrs["groups"] = groups
	}
	attrs["zone"] = zone

	var newHost NewHostStruct
	newHost.Name = hostname
	newHost.Type = "Host"
	newHost.Attrs = attrs
	newHost.Templates = templates

	// Create JSON from completed struct
	payloadJSON, marshalErr := json.Marshal(newHost)
	if marshalErr != nil {
		return nil, marshalErr
	}

	// Make the API request to create the hosts.
	results, err := server.NewAPIRequest("PUT", "/objects/hosts/"+hostname, []byte(payloadJSON))
	if err != nil {
		return nil, err
	}

	if results.Code == 200 {
		hosts, err := server.GetHost(hostname)
		return hosts, err
	}

	return nil, fmt.Errorf("%s", results.ErrorString)

}

// DeleteHost ...
func (server *Server) DeleteHost(hostname string) error {

	results, err := server.NewAPIRequest("DELETE", "/objects/hosts/"+hostname+"?cascade=1", nil)
	if err != nil {
		return err
	}

	if results.Code == 200 {
		return nil
	} else {
		return fmt.Errorf("%s", results.ErrorString)
	}

}
