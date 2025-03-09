package utils

import (
	"bytes"
	"encoding/json"
)

func PrettifyJson(json_ string) string {
	src := []byte(json_)

	dst := &bytes.Buffer{}
	if err := json.Indent(dst, src, "", "  "); err != nil {
		return json_
	}

	return dst.String()
}
