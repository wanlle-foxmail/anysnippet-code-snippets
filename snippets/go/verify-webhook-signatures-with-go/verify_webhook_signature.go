package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

func VerifyWebhookSignature(body []byte, signatureHeader string, secretKey string) bool {
	// Flow:
	//   parse sha256=... header
	//      |
	//      +-> malformed header -> return false
	//      `-> compute HMAC-SHA256 -> compare digests in constant time
	providedDigest, ok := parseSignatureHeader(signatureHeader)
	if !ok {
		return false
	}

	computedDigest := computeWebhookDigest(body, secretKey)
	return hmac.Equal(providedDigest, computedDigest)
}

func parseSignatureHeader(signatureHeader string) ([]byte, bool) {
	parts := strings.SplitN(strings.TrimSpace(signatureHeader), "=", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "sha256") {
		return nil, false
	}

	digest, err := hex.DecodeString(parts[1])
	if err != nil {
		return nil, false
	}

	return digest, true
}

func computeWebhookDigest(body []byte, secretKey string) []byte {
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write(body)
	return mac.Sum(nil)
}

func main() {
	body := []byte(`{"event":"ping"}`)
	secretKey := "demo-secret"
	signatureHeader := "sha256=" + hex.EncodeToString(computeWebhookDigest(body, secretKey))

	fmt.Println(VerifyWebhookSignature(body, signatureHeader, secretKey))
}
