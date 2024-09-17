package encrypt

import "encoding/base64"

func Base64Decode(base64Str string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(base64Str)
}
