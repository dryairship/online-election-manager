package utils

import (
    "math/rand"
    "time"
)

func InitializeRandomSeed() {
    rand.Seed(time.Now().UnixNano())
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GetRandomAuthCode() string {
    b := make([]byte, 32)
    for i := range b {
        b[i] = letters[rand.Int63()%int64(len(letters))]
    }
    return string(b)
}
