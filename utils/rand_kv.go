package utils

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	randStr = rand.New(rand.NewSource(time.Now().Unix()))
	letters = []byte("abcdefghijklmnopqrstuvwXyZABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)

// GetTestKey 获取测试使用的key
func GetTestKey(i int) []byte {
	return []byte(fmt.Sprintf("bietcask-go-key-%09d", i))
}

// RandomValue生成随机value 勇于测试
func RandomValue(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[randStr.Intn(len(letters))]
	}
	return []byte("bitcask-go-value" + string(b))
}
