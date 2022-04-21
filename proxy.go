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
	"net/http"
	"os"
	"strconv"

	"github.com/charted-dev/registry-proxy/auth"
	"github.com/charted-dev/registry-proxy/client"
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
	Auth *auth.BaseAuth

	// If the server is secured by HTTPS, which will use `https://` instead
	// of `http://`
	Secure bool

	// Client is the http.Client to extend.
	Client *http.Client
}

type Proxy struct {
	client *client.RegistryClient
	options *Options
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

	httpClient := &http.Client{}
	auth := auth.NewNoAuth()
	registryClient, err := client.NewRegistryClient(
		httpClient,
		false,
		host,
		port,
		&auth,
	)

	if err != nil {
		return nil, err
	}

	proxy := &Proxy{
		client: registryClient,
		options: &Options{
			Host: host,
			Port: port,
			Auth: &auth,
			Secure: false,
		},
	}

	return proxy, nil
}

// New creates a new *Proxy object with custom options.
func New(options *Options) (*Proxy, error) {
	if options.Client == nil {
		options.Client = &http.Client{}
	}
	
	registryClient, err := client.NewRegistryClient(
		options.Client,
		options.Secure,
		options.Host,
		options.Port,
		options.Auth,
	)

	if err != nil {
		return nil, err
	}

	proxy := &Proxy{
		registryClient,
		options,
	}

	return proxy, nil
}
