// 🪟 registry-proxy: Pluggable Docker Registry proxy for HTTP servers, made for charted-server.
// Copyright 2022 Noelware <team@noelware.org>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package registryproxy

import (
	"net/http"

	"github.com/charted-dev/registry-proxy/handlers"
)

type handler interface {
	HandleFunc(pattern string, handler http.HandlerFunc)
}

func (proxy *Proxy) AppendServerHandlers(server *http.Server) {
	repoHandler := handlers.RepositoryHandler{Client: proxy.client}

	if mux, ok := server.Handler.(handler); ok {
		println("uwu")
		mux.HandleFunc("/v2", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Header().Set("Content-Length", "0")
		})

		mux.HandleFunc("/v2/_catalog", func(w http.ResponseWriter, r *http.Request) {
			repoHandler.GetRepositories(w, r)
		})
	}

	if mux, ok := server.Handler.(*http.ServeMux); ok {
		println("uwu")
		mux.HandleFunc("/v2", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Header().Set("Content-Length", "0")
		})

		mux.HandleFunc("/v2/_catalog", func(w http.ResponseWriter, r *http.Request) {
			repoHandler.GetRepositories(w, r)
		})
	}
}
