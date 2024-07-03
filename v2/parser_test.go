package v2

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupTest(t *testing.T, body string) (*Parser, http.ResponseWriter, *http.Request) {
	t.Helper()

	parser := Default()

	w := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	if err != nil {
		t.Fatal(err.Error())
	}

	return parser, w, r
}

func TestNewParser(t *testing.T) {
	parser := NewParser(0, false)
	if parser.Limit < 1 {
		t.Fatal("parser limit checks not handled properly")
	}
}
