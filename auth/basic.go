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

package auth

import "net/http"

type BasicAuth struct {
	Username string
	Password string
}

func NewBasicAuth(username string, password string) BaseAuth {
	return &BasicAuth{
		Username: username,
		Password: password,
	}
}

func (*BasicAuth) Name() string {
	return "basic authentication with username + password"
}

func (b *BasicAuth) Configure(headers http.Header) error {
	return nil
}
