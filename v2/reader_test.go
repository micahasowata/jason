package v2

import (
	"testing"
)

func TestNonPointerDestination(t *testing.T) {
	parser, w, r := setupTest(t, "")

	err := parser.Read(w, r, map[string]any{})
	if err == nil {
		t.Fatal("destination checks not properly handled")
	}
}

func TestInvalidContentType(t *testing.T) {
	parser, w, r := setupTest(t, "")
	r.Header.Set("Content-Type", "text/plain")

	err := parser.Read(w, r, map[string]any{})
	if err == nil {
		t.Fatal("content type checks not properly handled")
	}
}

func TestVeryLargeBody(t *testing.T) {
	parser, w, r := setupTest(t, `{"name": "jason"}`)
	parser.Limit = 2
	r.Header.Set("Content-Type", "application/json")

	var input struct {
		Name string `json:"name"`
	}

	err := parser.Read(w, r, &input)
	if err == nil {
		t.Fatal("body size limit checks not properly handled")
	}
}

func TestRequestBodyErrors(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		{
			name: "empty body",
			body: "{}",
			want: "body is empty",
		},
		{
			name: "has syntax error",
			body: `<?xml version="1.0" encoding="UTF-8"?><note><to>Jason</to></note>`,
			want: "body contains invalid syntax at character 1",
		},
		{
			name: "has field with invalid value",
			body: `{"name": 123}'`,
			want: `body contains invalid value type for "name" at character 12`,
		},
		{
			name: "has invalid structure",
			body: `["foo", "bar"]`,
			want: "body contains invalid value type at character 1",
		},
		{
			name: "contains unknown field",
			body: `{"email": "p@happy.com"}`,
			want: `contains unknown field "email"`,
		},
		{
			name: "has multiple objects",
			body: `{"name":"Jason"}{"name":"Xoxo"}`,
			want: "body contains multiple json objects",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser, w, r := setupTest(t, tt.body)
			r.Header.Set("Content-Type", "application/json")

			var input struct {
				Name string `json:"name"`
			}

			err := parser.Read(w, r, &input)
			if err != nil {
				if err.Error() != tt.want {
					t.Logf("wanted %s got %s", tt.want, err.Error())
				}
			}
		})
	}
}

func TestRead(t *testing.T) {
	parser, w, r := setupTest(t, `{"name": "jason"}`)
	r.Header.Set("Content-Type", "application/json")

	var input struct {
		Name string `json:"name"`
	}

	err := parser.Read(w, r, &input)
	if err != nil {
		t.Fatal(err.Error())
	}

	if input.Name != "jason" {
		t.Fatal("request not decoded")
	}
}
