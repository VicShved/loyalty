package common

import (
	"crypto/sha256"
	"encoding/base64"
	"hash/fnv"
	"strconv"
)

// hash
func Hash32(s string) string {
	h := fnv.New32a()
	_, err := h.Write([]byte(s))
	if err != nil {
		return ""
	}
	return strconv.Itoa(int(h.Sum32()))
}

func HashSha256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	result := h.Sum(nil)
	return base64.RawStdEncoding.EncodeToString(result)
}
