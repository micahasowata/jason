package jason

import (
	"bytes"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	ContentType     = "Content-Type"
	ContentTypeJSON = "application/json"
)

// isBodyJSON ensures the media type of the request is set to application/json
func isBodyJSON(r *http.Request) bool {
	contentTypeHeader := r.Header.Get(ContentType)
	mediaType := strings.ToLower(strings.TrimSpace(strings.Split(contentTypeHeader, ";")[0]))
	return mediaType == ContentTypeJSON
}

func isDstPointer(dst interface{}) bool {
	return reflect.TypeOf(dst).Kind() == reflect.Ptr
}

func isBodyTooLarge(err error) bool {
	return bytes.Contains([]byte(err.Error()), []byte("http: request body too large"))
}

func isBadlyFormedJSON(err error) bool {
	return bytes.Contains([]byte(err.Error()), []byte("readObjectStart:")) ||
		bytes.HasPrefix([]byte(err.Error()), []byte("skipFourBytes: expect null")) ||
		bytes.Contains([]byte(err.Error()), []byte(": readStringSlowPath: unexpected end of input"))
}

func containsImproperAssignment(err error) bool {
	return bytes.Contains([]byte(err.Error()), []byte(": Read")) ||
		bytes.HasPrefix([]byte(err.Error()), []byte("Read"))
}

func isEmpty(err error) bool {
	return bytes.Equal([]byte(err.Error()), []byte("EOF"))
}

func findErrorLocation(input string) int {
	if !strings.Contains(input, "#") {
		return 0
	}

	if !strings.Contains(input, "byte") {
		return 0
	}

	startIndex := strings.Index(input, "#") + 1
	endIndex := strings.Index(input, "byte") - 1

	cleanedValue := strings.TrimSpace(input[startIndex:endIndex])

	pos, _ := strconv.Atoi(cleanedValue)

	return pos
}

func isUnknownField(err error) bool {
	return bytes.HasPrefix([]byte(err.Error()), []byte("ReadObject: found unknown field:"))
}

func getFieldName(input string) string {
	pattern := `\{([^:]+):`

	re := regexp.MustCompile(pattern)

	match := re.FindStringSubmatch(input)

	if len(match) < 2 {
		return ""
	}

	return match[1]
}

func getStmt(input string) string {
	pattern := `\|([^|]+)\|`
	re := regexp.MustCompile(pattern)

	endIndex := len(input)
	if !strings.Contains(input, "bigger") {
		return ""
	}

	startIndex := strings.Index(input, "bigger")
	if startIndex == -1 {
		return ""
	}

	match := re.FindStringSubmatch(input[startIndex:endIndex])

	if len(match) < 2 {
		return ""
	}

	return match[1]
}
