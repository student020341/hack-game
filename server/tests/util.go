package tests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func SimpleGet(path string) (*int, []byte, error) {
	res, err := http.Get(path)
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

func SimplePost(path string, iface interface{}) (*int, []byte, error) {
	asJson, err := json.MarshalIndent(iface, "", "  ")
	if err != nil {
		return nil, nil, err
	}

	res, err := http.Post(path, "application/json", bytes.NewBuffer(asJson))
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
