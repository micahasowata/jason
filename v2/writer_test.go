package v2

import (
	"net/http"
	"testing"
)

func TestIndentWrite(t *testing.T) {
	t.Run("indent", func(t *testing.T) {
		parser, w, _ := setupTest(t, "")
		parser.Indent = true
		err := parser.Write(w, 200, map[string]any{
			"name": "jason",
		})

		if err != nil {
			t.Fatal(err.Error())
		}
	})

	t.Run("error", func(t *testing.T) {
		parser, w, _ := setupTest(t, "")
		parser.Indent = true

		intChan := make(chan int)

		err := parser.Write(w, http.StatusOK, intChan)
		if err == nil {
			t.Fatal("unsupported type checks not properly handled")
		}
	})
}

func TestWrite(t *testing.T) {
	t.Run("write", func(t *testing.T) {
		parser, w, _ := setupTest(t, "")
		err := parser.Write(w, 512, map[string]any{
			"name": "jason",
		})

		if err != nil {
			t.Fatal(err.Error())
		}
	})

	t.Run("write_error", func(t *testing.T) {
		parser, w, _ := setupTest(t, "")
		intChan := make(chan int)

		err := parser.Write(w, http.StatusOK, intChan)
		if err == nil {
			t.Fatal("unsupported type checks not properly handled")
		}
	})
}
