## jason

[![Go Reference](https://pkg.go.dev/badge/github.com/micahasowata/jason.svg)](https://pkg.go.dev/github.com/micahasowata/jason) [![Go Report Card](https://goreportcard.com/badge/github.com/micahasowata/jason)](https://goreportcard.com/report/github.com/micahasowata/jason)

This package is a product of [this article by Alex Edwards](https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body) and [jsoniter](https://github.com/json-iterator/go)

### Installation 📥

```go
go get github.com/micahasowata/jason@latest
```

### Quick Start 💨

```go
package main

import (
	"log/slog"
	"net/http"

	"github.com/micahasowata/jason"
)

func main() {
	j := jason.New(100, true, true)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Name string `json:"name"`
		}

		err := j.Read(w, r, &input)
		if err != nil {
			e, ok := err.(*jason.Err)
			if !ok {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			http.Error(w, e.Error(), e.Code)
			return
		}

		err = j.Write(w, http.StatusOK, jason.Envelope{"data": input}, nil)
		if err != nil {
			e, ok := err.(*jason.Err)
			if !ok {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			http.Error(w, e.Error(), e.Code)
			return
		}
	})

	slog.Info("server started")
	http.ListenAndServe(":8000", nil)
}

```

## Contribution 🗳️

If you run into any issues while using this library or you notice any area where it can be improved feel free to create a pull request and be sure that it would not be ignored

Happy Hacking ⚡⚡⚡
