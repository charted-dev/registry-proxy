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

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/charted-dev/registry-proxy/client"
)

type RepositoryHandler struct {
	Client *client.RegistryClient
}

type Error struct {
	Message string `json:"message"`
	Code string `json:"code"`
	Detail *string `json:"detail,omitempty"`
}

func (h *RepositoryHandler) GetRepositories(w http.ResponseWriter, req *http.Request) {
	var data map[string]any
	headers, err := h.Client.Request("GET", "/_catalog", nil, &data)
	if err != nil {
		w.WriteHeader(400)
		_, _ = w.Write([]byte("lol fuck you"))

		return
	}

	fmt.Println(data)

	for key, value := range headers {
		w.Header().Set(key, strings.Join(value, ","))
	}

	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}
}
