// 🪟 registry-proxy: Pluggable Docker Registry proxy for HTTP servers, made for charted-server.
// Copyright 2022 Noelware <team@noelware.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package client represents the registry client for connecting the registry
// to the HTTP server.
package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/charted-dev/registry-proxy/auth"
)

// RegistryClient represents the main registry client structure.
type RegistryClient struct {
	client *http.Client
	secure bool
	auth *auth.BaseAuth
	host string
	port int
}

type repositoriesResult struct {
	Repositories []string `json:"repositories"`
}

func NewRegistryClient(client *http.Client, secure bool, host string, port int, auth *auth.BaseAuth) (*RegistryClient, error) {
	cli := &RegistryClient{client, secure, auth, host, port}

	if err := cli.tryConnect(); err != nil {
		return nil, err
	}

	return cli, nil
}

func (client *RegistryClient) Repositories() ([]string, error) {
	var data repositoriesResult
	if _, err := client.Request("GET", "/_catalog", nil, &data); err != nil {
		return []string{}, err
	}

	return data.Repositories, nil
}

func (client *RegistryClient) tryConnect() error {
	_, err := client.Request("HEAD", "/", nil, nil)
	return err
}

func (client *RegistryClient) Request(method string, endpoint string, body io.Reader, data interface{}) (http.Header, error) {
	protocol := "http"
	if client.secure {
		protocol = "https"
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s://%s:%d/v2%s", protocol, client.host, client.port, endpoint), body)
	if err != nil {
		return nil, err
	}

	isGetOrHead := false
	if req.Method == "GET" {
		isGetOrHead = true
	}

	if req.Method == "HEAD" {
		isGetOrHead = true
	}

	if isGetOrHead && body != nil {
		return nil, errors.New("cannot append `body` to GET/HEAD requests")
	}

	if client.auth != nil {
		auth := *client.auth
		err = auth.Configure(req.Header)
		if err != nil {
			return nil, err
		}
	}

	res, err := client.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		buf := &bytes.Buffer{}

		if _, err := io.Copy(buf, res.Body); err != nil {
			return res.Header, err
		}

		return res.Header, fmt.Errorf("received %d, not 200 :: %s", res.StatusCode, buf.String())
	}

	// Convert the data to JSON, if any
	contentType := res.Header.Get("Content-Type")
	if contentType == "" && data != nil {
		return res.Header, errors.New("missing `content-type` header but data to deserialise was provided")
	}
	
	if data != nil {
		if contentType != "" && strings.HasPrefix(contentType, "application/json") {
			if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
				return res.Header, err
			}
		}
	}

	return res.Header, nil
}
