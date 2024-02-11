package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type APIGroupList struct {
	Groups []struct {
		Name             string `json:"name"`
		PreferredVersion struct {
			Version string `json:"version"`
		} `json:"preferredVersion"`
	} `json:"groups"`
}

type APIResourceList struct {
	Resources []struct {
		Name  string   `json:"name"`
		Verbs []string `json:"verbs"`
	} `json:"resources"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <server_address>")
		os.Exit(1)
	}
	server := os.Args[1]

	// Ensure the server URL starts with http:// or https://
	if !strings.HasPrefix(server, "http://") && !strings.HasPrefix(server, "https://") {
		fmt.Println("Server address must start with http:// or https://")
		os.Exit(1)
	}

	// Get list of all API groups
	groupList := getAPIGroups(server + "/apis")

	// Process core resources first
	processCoreResources(server)

	// Now process non-core resources
	for _, group := range groupList.Groups {
		api := group.Name
		version := group.PreferredVersion.Version
		resourcesURL := fmt.Sprintf("%s/apis/%s/%s", server, api, version)
		processResources(api, resourcesURL)
	}
}

func getAPIGroups(url string) APIGroupList {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var groupList APIGroupList
	if err := json.Unmarshal(body, &groupList); err != nil {
		panic(err)
	}

	return groupList
}

func processCoreResources(server string) {
	url := server + "/api/v1"
	processResources("core", url)
}

func processResources(api, url string) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var resourceList APIResourceList
	if err := json.Unmarshal(body, &resourceList); err != nil {
		panic(err)
	}

	for _, resource := range resourceList.Resources {
		fmt.Printf("%s %s: %s\n", api, resource.Name, strings.Join(resource.Verbs, " "))
	}
}
