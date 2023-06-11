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

	"github.com/fatih/color"
)

type EmbeddingsResponse struct {
	Embeddings []struct {
		Index int       `json:"index"`
		Data  []float64 `json:"data"`
	} `json:"embeddings"`
	Model      string `json:"model"`
	Dimensions int    `json:"dimensions"`
}

type AnthropicResponse struct {
	Completion string      `json:"completion"`
	StopReason string      `json:"stop_reason"`
	Model      string      `json:"model"`
	Truncated  bool        `json:"truncated"`
	Stop       string      `json:"stop"`
	LogID      string      `json:"log_id"`
	Exception  interface{} `json:"exception"`
}

type OpenAIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
	} `json:"choices"`
}

type Message struct {
	Speaker string `json:"speaker"`
	Text    string `json:"text"`
}

type Messages struct {
	Messages []Message `json:"messages"`
}

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIMessages struct {
	Messages []Message `json:"messages"`
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

// write a function make Anthropic Completion API - parameters host, access tokens, message , max token int, temperature float, prompt string
func AnthropicAPI(host string, accessToken string, messages []Message, maxTokens int, temperature float32, prompt string, mode string) (resp AnthropicResponse, err error) {
	var anthropicResponse AnthropicResponse
	// make api call to path /v1/completions/anthropic
	url := fmt.Sprintf("%s/v1/completions/anthropic", host)
	// create jsonData as per above sample json schema
	jsonData, _ := json.Marshal(map[string]interface{}{
		"prompt":               prompt,
		"messages":             messages,
		"model":                "claude-v1",
		"max_tokens_to_sample": maxTokens,
		"temperature":          temperature,
	})
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("X-Sourcegraph-Feature", "chat_completions")

	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return anthropicResponse, err
	}
	defer response.Body.Close()

	fmt.Println(response.StatusCode)

	if response.StatusCode != http.StatusOK {
		err = errors.New("failed to call embeddings api")
		return
	}

	// read the response body
	body, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(body, &anthropicResponse)
	if err != nil {
		// Handle error
		return anthropicResponse, err
	}
	return anthropicResponse, err
}

// write a function make Anthropic Completion API - parameters host, access tokens, message , max token int, temperature float, prompt string
func OpenAIAPI(host string, accessToken string, messages []OpenAIMessage, maxTokens int, temperature float32, mode string) (resp OpenAIResponse, err error) {
	var openaiResponse OpenAIResponse
	// make api call to path /v1/completions/openai
	url := fmt.Sprintf("%s/v1/completions/openai", host)
	// create jsonData as per above sample json schema
	jsonData, _ := json.Marshal(map[string]interface{}{
		"messages":    messages,
		"model":       "gpt-3.5-turbo",
		"max_tokens":  maxTokens,
		"temperature": temperature,
	})
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("X-Sourcegraph-Feature", "chat_completions")

	client := http.Client{}
	response, err := client.Do(req)

	if err != nil {
		return openaiResponse, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = errors.New("failed to call embeddings api")
		return
	}

	// read the response body
	body, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(body, &openaiResponse)
	if err != nil {
		// Handle error
		return openaiResponse, err
	}
	return openaiResponse, err
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
	} else if c.AnthropicCompletionAPI {
		fmt.Println("🪄 🪄 🪄 🪄 🪄 🪄 ")
		fmt.Println("Establishing Session with Anthropic AI 🪄 ✨ (supports multi-line and type --END-- to terminate input):")
		fmt.Println("🪄 🪄 🪄 🪄 🪄 🪄 ")

		// get multiline input from user and save it as string
		// declare string array
		chatHistory := []Message{}

		input := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("You -> ")
			line, _ := input.ReadString('\n')
			// convert CRLF to LF
			line = strings.Replace(line, "\n", "", -1)

			if strings.Compare("--END--", line) == 0 {
				// break this for loop
				break
			} else {
				resp, err := AnthropicAPI(c.GatewayHost, c.GatewayToken, chatHistory, 500, 0.1, line, c.CompletionMode)
				if err != nil {
					return err
				}
				fmt.Println("🪄 🪄 🪄 🪄 🪄 🪄 ")
				fmt.Println("------------")
				color.Green("%s", resp.Completion)
				fmt.Println("------------")
				fmt.Println("🪄 🪄 🪄 🪄 🪄 🪄 ")
				// append text to message interaction
				// create Message object and append speaker as human and keyphrase as text
				// create Message object and append speaker as system adn resp.completion as text
				humanMessage := Message{
					Speaker: "human",
					Text:    line,
				}
				chatHistory = append(chatHistory, humanMessage)
				systemMessage := Message{
					Speaker: "system",
					Text:    resp.Completion,
				}
				chatHistory = append(chatHistory, systemMessage)
			}
		}
	} else if c.OpenAICompletionAPI {
		fmt.Println("🪄 🪄 🪄 🪄 🪄 🪄 ")
		fmt.Println("Establishing Session with OpenAI GPT-3.5-Turbo 🪄 ✨ (supports multi-line and type --END-- to terminate input):")
		fmt.Println("🪄 🪄 🪄 🪄 🪄 🪄 ")

		// get multiline input from user and save it as string
		// declare string array
		chatHistory := []OpenAIMessage{}

		input := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("You -> ")
			line, _ := input.ReadString('\n')
			// convert CRLF to LF
			line = strings.Replace(line, "\n", "", -1)

			if strings.Compare("--END--", line) == 0 {
				// break this for loop
				break
			} else {
				humanMessage := OpenAIMessage{
					Role:    "user",
					Content: line,
				}
				chatHistory = append(chatHistory, humanMessage)

				resp, err := OpenAIAPI(c.GatewayHost, c.GatewayToken, chatHistory, 500, 0.1, c.CompletionMode)
				if err != nil {
					return err
				}
				fmt.Println("🪄 🪄 🪄 🪄 🪄 🪄 ")
				fmt.Println("------------")
				color.Green("%s", resp.Choices[0].Message.Content)
				fmt.Println("------------")
				fmt.Println("🪄 🪄 🪄 🪄 🪄 🪄 ")
				// append text to message interaction
				// create Message object and append speaker as human and keyphrase as text
				// create Message object and append speaker as system adn resp.completion as text

				systemMessage := OpenAIMessage{
					Role:    "assistant",
					Content: resp.Choices[0].Message.Content,
				}
				chatHistory = append(chatHistory, systemMessage)
			}
		}
	}
	return nil
}
