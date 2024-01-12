// Copyright The KCL Authors. All rights reserved.

package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func httpPost(urlpath string, input, output interface{}) error {
	const method = "POST"

	reqBody, err := json.Marshal(input)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, urlpath, bytes.NewReader(reqBody))
	if err != nil {
		return err
	}

	// python -m kclvm.program.rpc-server -http=<addr> do not support application/json Content-Type
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{
		Transport: &http.Transport{
			DisableCompression: true,
		},
	}

	r, err := client.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	bodyData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bodyData, output)
	if err != nil {
		return fmt.Errorf("json.Unmarshal failed: bodyData = %v", string(bodyData))
	}

	return nil
}
