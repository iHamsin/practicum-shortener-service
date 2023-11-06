package util

import (
	"math/rand"
	"strings"
)

func randomString(n int) string {
	// словарь
	alphabet := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	alphabetSize := len(alphabet)
	var sb strings.Builder
	for i := 0; i < n; i++ {
		ch := alphabet[rand.Intn(alphabetSize)]
		sb.WriteRune(ch)
	}
	return sb.String()
}
