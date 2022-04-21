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

// Package registryproxy is the package for implementing a HTTP proxy towards a Docker Registry.
//
// This can be intertwined with anything that can be related to a local registry but this
// was made for charted-server's use case for OCI registries with Helm. You can read more
// here: https://charts.noelware.org/docs/libraries/registry-proxy.
//
// You can use the library if you want to implement an authorization server for the
// registry itself.
//
//
// USAGE:
//
//    package main
//
//    import (
//        "log"
//        "net/http"
//        proxy "github.com/charted-dev/registry-proxy"
//    )
//
//    func main() {
//        p, err := proxy.NewDefault() // with default options
//        if err != nil {
//             log.Fatal(err)
//        }
//
//        server := &http.Server{
//             Handler: // some handler here!
//        }
//
//        if err := server.ListenAndServe(); err != nil {
//            log.Fatal(err)
//        }
//    }
package registryproxy

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

// Options represents the base options for the proxy itself.
type Options struct {
	// Host returns the host of the registry itself. It will also look for the "REGISTRY_HOST"
	// environment variable.
	Host string

	// Port returns the port of the registry itself. It will also look for the "REGISTRY_PORT"
	// environment variable.
	Port int

	// Represents the authentication type. Uses `nil` as the default.
	Auth *Auth

	// If the server is secured by HTTPS, which will use `https://` instead
	// of `http://`
	Secure bool

	// Client is the http.Client to extend.
	Client *http.Client
}

type Proxy struct {
	options *Options
}

// Auth represents the authentication structure for connecting to a registry.
// Available options are: BasicAuth, TokenAuth, SillyAuth, NoAuth.
type Auth interface {
	// Name returns the name of this authentication type.
	Name() string
	
	// IsNoAuth returns a bool if the authentication type is `NoAuth`.
	IsNoAuth() bool

	// Configure configures the authentication type for the HTTP client.
	Configure(headers http.Header) error
}

type BasicAuth struct {
	Username string
	Password string
}

type NoAuth struct {}

// NewBasicAuth creates a new Auth object being a *BasicAuth object.
func NewBasicAuth(username string, password string) Auth {
	return &BasicAuth{
		Username: username,
		Password: password,
	}
}

func (*BasicAuth) Name() string {
	return "basic authentication with username + password"
}

func (*BasicAuth) IsNoAuth() bool {
	return false
}

func (b *BasicAuth) Configure(headers http.Header) error {
	return nil
}

func NewNoAuth() Auth {
	return &NoAuth{}
}

func (*NoAuth) Name() string {
	return "no authentication provided"
}

func (*NoAuth) IsNoAuth() bool {
	return true
}

func (*NoAuth) Configure(headers http.Header) error {
	return nil
}

// NewDefault creates a new *Proxy object with the default options.
func NewDefault() (*Proxy, error) {
	host := ""
	if ho, ok := os.LookupEnv("REGISTRY_HOST"); ok {
		host = ho
	}

	if host == "" {
		host = "0.0.0.0"
	}

	port := 5000
	if po, ok := os.LookupEnv("REGISTRY_PORT"); ok {
		p, err := strconv.Atoi(po)
		if err != nil {
			return nil, err
		}

		port = p
	}

	client := &http.Client{}
	auth := NewNoAuth()

	proxy := &Proxy{
		options: &Options{
			Host: host,
			Port: port,
			Auth: &auth,
			Client: client,
		},
	}

	if err := proxy.TryConnect(); err != nil {
		return nil, err
	}

	return proxy, nil
}

// New creates a new *Proxy object with custom options.
func New(options *Options) (*Proxy, error) {
	proxy := &Proxy{options}

	if err := proxy.TryConnect(); err != nil {
		return nil, err
	}

	return proxy, nil
}

// TryConnect makes a connection to the registry to check if it's a successful connection.
func (proxy *Proxy) TryConnect() error {
	protocol := "http"
	if proxy.options.Secure {
		protocol = "https"
	}

	req, err := http.NewRequest("HEAD", fmt.Sprintf("%s://%s:%d/v2", protocol, proxy.options.Host, proxy.options.Port), nil)
	if err != nil {
		return err
	}

	if proxy.options.Auth != nil {
		auth := *proxy.options.Auth
		err = auth.Configure(req.Header)
		if err != nil {
			return err
		}
	}

	res, err := proxy.options.Client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		defer res.Body.Close()
		buf := &bytes.Buffer{}

		if _, err := io.Copy(buf, res.Body); err != nil {
			return err
		}

		return fmt.Errorf("received %d, not 200 :: %s", res.StatusCode, buf.String())
	}

	return nil
}
