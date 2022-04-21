# 🪟 Docker Registry Proxy
> *Pluggable Docker Registry proxy for HTTP servers, made for charted-server*

## Why?
During the development of **charted-server**, it was hard to dissect the source code for [registry](https://github.com/distribution/distribution/tree/main/registry) without actually... diving into it, so this is a library to abstract the proxying for OCI registries.

This is useful for **charted-server** so people can use a local OCI registry instead of a "chart-based" repository that was first introduced in Helm 2.8. This is compatible with the `helm push` command.

For the official instance of **charted-server**, you can push your charts to the server as:

```shell
$ helm push <package>.tar.gz --url=oci://charts.noelware.org/<owner>/<repo>
```

and it will proxy it as `charts.noelware.org/api/v2/<owner>/<repo>`, possible with this library.

## Usage
```go
package main

import (
    "log"
    "net/http"
    regproxy "github.com/charted-dev/registry-proxy"
)

func main() {
    // Create your proxy
    proxy, err := regproxy.New(&regproxy.Options{
        Host: "localhost",
        Port: 5000,
        Auth: &regproxy.BasicAuth{
            Username: "<username>",
            Password: "<password>",
        }
    })

    // Create your HTTP server like this so
    // RegistryProxy.RegisterHandlers(*http.Server) can be registered
    // correctly.
    server := &http.Server{
        Handler: // some handler
    }

    proxy.RegisterHandlers(server)
    log.Fatal(server.ListenAndServe())
}
```

## Contributing
Thanks for considering contributing to **registry-proxy**! Before you boop your heart out on your keyboard ✧ ─=≡Σ((( つ•̀ω•́)つ, we recommend you to do the following:

- Read the [Code of Conduct](./.github/CODE_OF_CONDUCT.md)
- Read the [Contributing Guide](./.github/CONTRIBUTING.md)

If you read both if you're a new time contributor, now you can do the following:

- [Fork me! ＊*♡( ⁎ᵕᴗᵕ⁎ ）](https://github.com/charted-dev/registry-proxy/fork)
- Clone your fork on your machine: `git clone https://github.com/your-username/registry-proxy`
- Create a new branch: `git checkout -b some-branch-name`
- BOOP THAT KEYBOARD!!!! ♡┉ˏ͛ (❛ 〰 ❛)ˊˎ┉♡ ✧ ─=≡Σ((( つ•̀ω•́)つ
- Commit your changes onto your branch: `git commit -am "add features （｡>‿‿<｡ ）"`
- Push it to the fork you created: `git push -u origin some-branch-name`
- Submit a Pull Request and then cry! ｡･ﾟﾟ･(థ Д థ。)･ﾟﾟ･｡

## License
**registry-proxy** is released under the **MIT License** by Noelware with love. :purple_heart:
