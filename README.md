## jason

jason is a convenience package that provides proper error handling while working with JSON APIs

### install

- install the package in a project with a valid [go.mod](https://go.dev/doc/modules/gomod-ref) file

```shell
$ go get github.com/micahasowata/jason@latest
```

- use it to read and write json

```go

package main

import (
    "github.com/micahasowata/jason"
)

func main() {
   j := jason.New(100, true, true)

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
        var input struct {
            Name string `json:"name"`
        }

        err := j.Read(w, r, &input)
        if err != nil {
            http.Error(w, err.Error, http.StatusBadRequest)
            return
        }

        err := j.Write(w, http.StatusOK, jason.Envelope{"input", input}, nil)
        if err != nil {
            http.Error(w, err.Error, http.StatusInternalServerError)
            return
        }
    })

   http.ListenAndServe(":8000", nil)
}
```

## Contribution

this package is not in its best possible shape
there is a need for tests and better documentation so if you can contribute that, please do so.

Thank you
