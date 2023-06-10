package app

import (
	"cody-gateway-cli/config"
	"fmt"
	"io/ioutil"
	"net/http"
)

// write a function to make HTTP Get request
func MakeGetRequest(url string) (response string, err error) {
	// TODO: implement me
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	// Read the response
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Return the response
	return string(body), nil
}

// create a function that accepts Host parameter and concatenates "-/__version" to the path
// and makes a GET request
func GetVersionInfo(host string) (version string, err error) {
	url := fmt.Sprintf("%s/-/__version", host)
	// add request header Authorization as bearer sekret value
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "bearer XXX")
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func Run(c config.Config) error {
	version, err := GetVersionInfo(c.GatewayHost)
	if err != nil {
		return err
	}
	fmt.Printf("Version: %s\n", version)
	return nil
}
