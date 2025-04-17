// From Alex Edward's "Let's Go Further":
//     Note: Your secret key should be a cryptographically secure random
//     string with an underlying entropy of at least 32 bytes (256 bits).
//
// This code is to generate such a secret key.

package main

import (
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "log"
)

func main() {
    key := make([]byte, 32)

    if _, err := rand.Read(key); err != nil {
        log.Fatal(err)
    }

    secretKey := base64.URLEncoding.EncodeToString(key)

    fmt.Println("Your secret key:", secretKey)
}

