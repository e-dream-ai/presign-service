package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
)

func main() {
	fmt.Println("=== API Key Generator ===\n")

	fmt.Println("Base64 API Key:")
	base64Key, err := generateBase64Key(32)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   %s\n\n", base64Key)

	fmt.Println("Usage:")
	fmt.Printf("export API_KEY=\"%s\"\n", base64Key)
	fmt.Println("\nNote: Store this key securely and never commit it to version control!")
}

// generateBase64Key generates a random base64 encoded API key
func generateBase64Key(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}
