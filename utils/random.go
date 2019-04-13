package utils

import (
    "github.com/dryairship/online-election-manager/config"
    "math/rand"
    "time"
)

func InitializeRandomSeed() {
    rand.Seed(time.Now().UnixNano())
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Function to return a random string to be used as authentication code.
func GetRandomAuthCode() string {
    b := make([]byte, 6)
    for i := range b {
        b[i] = letters[rand.Int63()%int64(len(letters))]
    }
    return string(b)
}

// Function to get a random time delay after which the vote will be added to the database.
func GetRandomTimeDelay() time.Duration {
    return time.Duration(rand.Intn(config.MaxTimeDelay)) * time.Second
}
