package calchash

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"strings"
)

type MyConstraint interface {
	int64 | float64
}

func Calculate[T MyConstraint](key string, mType string, id string, value T) string {
	if len(strings.TrimSpace(key)) == 0 {
		return ""
	}
	var data string
	if mType == "gauge" {
		data = fmt.Sprintf("%s:gauge:%f", id, value)
	} else {
		data = fmt.Sprintf("%s:counter:%d", id, value)
	}
	return calcHash(data, key)
}

func calcHash(data string, key string) string {
	src := []byte(data)
	h := hmac.New(sha256.New, []byte(key))
	h.Write(src)
	dst := h.Sum(nil)
	return fmt.Sprintf("%x", dst)
}
