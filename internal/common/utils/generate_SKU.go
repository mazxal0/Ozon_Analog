package utils

import (
	"crypto/rand"
	"fmt"
)

// GenerateSKU генерирует уникальный числовой SKU длиной 10 цифр
func GenerateSKU() string {
	b := make([]byte, 5) // 5 байт = 10 цифр в hex
	_, _ = rand.Read(b)
	return fmt.Sprintf("%010d", int64(b[0])<<32|int64(b[1])<<24|int64(b[2])<<16|int64(b[3])<<8|int64(b[4]))
}
