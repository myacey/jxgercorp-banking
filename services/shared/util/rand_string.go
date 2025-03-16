package util

import (
	"strings"

	"golang.org/x/exp/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString(n int) string {
	var sb strings.Builder
	sb.Grow(n)
	for i := 0; i < n; i++ {
		idx := rand.Intn(len(letterBytes))
		sb.WriteByte(letterBytes[idx])
	}
	return sb.String()
}
