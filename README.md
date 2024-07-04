# jason

[![Go Reference](https://pkg.go.dev/badge/github.com/micahasowata/jason/v2@v2.0.0.svg)](https://pkg.go.dev/github.com/micahasowata/jason/v2@v2.0.0) [![Go Report Card](https://goreportcard.com/badge/github.com/micahasowata/jason/v2)](https://goreportcard.com/report/github.com/micahasowata/jason/v2)

jason is a zero dependency http library wrapper for JSON.

## Installation 📥

```go
go get github.com/micahasowata/jason/v2@latest
```

## Quick Start 💨

### With `Default`

```go
package main

import (
	"net/http"
	"github.com/micahasowata/jason/v2"
)

func main() {
	parser := jason.Default()
	names := map[string]string{}

	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Name string `json:"name"`
		}

		err := parser.Read(w, r, &input)
		if err != nil {
			err, ok := err.(*jason.Err)
			if !ok {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			http.Error(w, err.Error(), err.Code)
			return
		}

		names["my_name"] = input.Name
	})

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		err := parser.Write(w, http.StatusOK, names["my_name"])
		if err != nil {
			err, ok := err.(*jason.Err)
			if !ok {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			http.Error(w, err.Error(), err.Code)
			return
		}
	})

	http.ListenAndServe(":8080", nil)
}
```

### With `NewParser`

```go
package main

import (
	"net/http"
	"github.com/micahasowata/jason/v2"
)

func main() {
	// Creates a parser that reads up to 10MB of data from stream and indents response
	parser := jason.NewParser(jason.OneMB*10, true)
	names := map[string]string{}

	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Name string `json:"name"`
		}

		err := parser.Read(w, r, &input)
		if err != nil {
			err, ok := err.(*jason.Err)
			if !ok {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			http.Error(w, err.Error(), err.Code)
			return
		}

		names["my_name"] = input.Name
	})

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		err := parser.Write(w, http.StatusOK, names["my_name"])
		if err != nil {
			err, ok := err.(*jason.Err)
			if !ok {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			http.Error(w, err.Error(), err.Code)
			return
		}
	})

	http.ListenAndServe(":8080", nil)
}
```

## Note

If you get any error when you import the v2 you can try aliasing it like this:

```go
import (
	jason "github.com/micahasowata/jason/v2"
)

// instead of

import (
	"github.com/micahasowata/jason/v2"
)

// should probably fix it. If it still doesn't get fixed open an issue.

```

## Contributing 🗳️

If you run into any issues while using this library or you notice any area where it can be improved feel free to create a pull request and be sure that it would not be ignored

Happy Hacking ⚡⚡⚡
