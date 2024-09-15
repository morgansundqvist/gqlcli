package executor

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/morgansundqvist/gqlcli/pkg/utils"
)

func ExecuteGraphQL(query string, variables map[string]interface{}, headers map[string]string, config *utils.Config) (map[string]interface{}, error) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", config.GraphQLURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var respData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return nil, err
	}

	return respData, nil
}
