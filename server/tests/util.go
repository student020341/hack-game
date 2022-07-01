package tests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func SimpleGet(path string, token *string) (*int, []byte, error) {

	client := http.DefaultClient
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	if token != nil {
		req.Header.Set("token", *token)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()
	return &res.StatusCode, body, nil
}

func SimplePost(path string, iface interface{}, token *string) (*int, []byte, error) {
	// post body
	var asJson []byte
	if iface != nil {
		var err error
		asJson, err = json.MarshalIndent(iface, "", "  ")
		if err != nil {
			return nil, nil, err
		}
	}

	// for headers
	client := http.DefaultClient
	req, err := http.NewRequest("POST", path, bytes.NewBuffer(asJson))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if token != nil {
		req.Header.Set("token", *token)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()
	return &res.StatusCode, body, nil
}
