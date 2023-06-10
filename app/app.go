package app

import (
	"bufio"
	"bytes"
	"cody-gateway-cli/config"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type EmbeddingsResponse struct {
	Embeddings []struct {
		Index int       `json:"index"`
		Data  []float64 `json:"data"`
	} `json:"embeddings"`
	Model      string `json:"model"`
	Dimensions int    `json:"dimensions"`
}

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
func GetVersionInfo(host string, debugSecretToken string) (version string, err error) {
	client := http.Client{}
	url := fmt.Sprintf("%s/-/__version", host)
	// add request header Authorization as bearer sekret value
	req, _ := http.NewRequest("GET", url, nil)
	req.Header = http.Header{
		"Authorization": {"Bearer " + debugSecretToken},
	}
	resp, err := client.Do(req)
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

// create a function that accepts Host parameter and concatenates "/-/healthz" to the path
// and make a GET request with debugsecrettoken in the authorization header
func HealthCheck(host string, debugSecretToken string) (ok bool, err error) {
	client := http.Client{}
	url := fmt.Sprintf("%s/-/healthz", host)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header = http.Header{
		"Authorization": {"Bearer " + debugSecretToken},
	}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}
	return false, nil
}

// write a function that accepts string array, host, access token to call v1/embeddings API
func EmbeddingsAPI(terms []string, host string, accessToken string) (embeddingResponse EmbeddingsResponse, err error) {
	var embeddingsResponse EmbeddingsResponse

	url := fmt.Sprintf("%s/v1/embeddings", host)
	jsonData, _ := json.Marshal(&map[string]interface{}{
		"model": "openai/text-embedding-ada-002",
		"input": terms,
	})

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("X-Sourcegraph-Feature", "chat_completions")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return embeddingResponse, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = errors.New("failed to call embeddings api")
		return
	}

	// read the response body
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &embeddingsResponse)
	if err != nil {
		// Handle error
		return embeddingResponse, err
	}
	return embeddingsResponse, err
}

func Run(c config.Config) error {

	if c.VersionAPI {
		version, err := GetVersionInfo(c.GatewayHost, c.DebugSecretToken)
		if err != nil {
			return err
		}
		fmt.Printf("Version: %s\n", version)
	} else if c.HealthCheckAPI {
		ok, err := HealthCheck(c.GatewayHost, c.DebugSecretToken)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("health check failed")
		}
	} else if c.EmbeddingsAPI {

		fmt.Println("Enter keyphrases to embed ✨ (supports multi-line and type --END-- to terminate input):")

		// get multiline input from user and save it as string
		// declare string array
		keyphrases := []string{}

		input := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("-> ")
			line, _ := input.ReadString('\n')
			// convert CRLF to LF
			line = strings.Replace(line, "\n", "", -1)

			if strings.Compare("--END--", line) == 0 {
				// break this for loop
				break
			} else {
				// append text to terms array
				keyphrases = append(keyphrases, line)
			}
		}

		resp, err := EmbeddingsAPI(keyphrases, c.GatewayHost, c.GatewayToken)
		if err != nil {
			return err
		}
		//TODO: Improve this print as table with colors - allow option to write embeddings to disk
		// print the output as table resp
		for _, e := range resp.Embeddings {
			fmt.Printf("%d\t%v\n", e.Index, e.Data)
		}
		fmt.Printf("%s\t%d\t%d\n", resp.Model, resp.Dimensions, len(resp.Embeddings))
	}
	return nil
}
