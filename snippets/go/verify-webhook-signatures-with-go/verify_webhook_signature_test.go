package main

import (
	"encoding/hex"
	"strings"
	"testing"
)

func validSignatureHeader(body []byte, secretKey string) string {
	return "sha256=" + hex.EncodeToString(computeWebhookDigest(body, secretKey))
}

func TestVerifyWebhookSignatureReturnsTrueForValidSignature(t *testing.T) {
	body := []byte(`{"event":"created"}`)
	if !VerifyWebhookSignature(body, validSignatureHeader(body, "secret"), "secret") {
		t.Fatal("expected valid signature to verify")
	}
}

func TestVerifyWebhookSignatureReturnsFalseForInvalidSignature(t *testing.T) {
	body := []byte(`{"event":"created"}`)
	if VerifyWebhookSignature(body, validSignatureHeader(body, "secret"), "other-secret") {
		t.Fatal("expected invalid signature to fail")
	}
}

func TestVerifyWebhookSignatureAllowsEmptyBodyWhenSignatureMatches(t *testing.T) {
	body := []byte{}
	if !VerifyWebhookSignature(body, validSignatureHeader(body, "secret"), "secret") {
		t.Fatal("expected empty body signature to verify")
	}
}

func TestVerifyWebhookSignatureRejectsEmptyHeader(t *testing.T) {
	if VerifyWebhookSignature([]byte(`{"event":"created"}`), "", "secret") {
		t.Fatal("expected empty signature header to fail")
	}
}

func TestVerifyWebhookSignatureAllowsEmptySecretWhenSignatureMatches(t *testing.T) {
	body := []byte(`{"event":"created"}`)
	if !VerifyWebhookSignature(body, validSignatureHeader(body, ""), "") {
		t.Fatal("expected empty secret signature to verify")
	}
}

func TestVerifyWebhookSignatureAcceptsUppercaseHex(t *testing.T) {
	body := []byte(`{"event":"created"}`)
	header := strings.ToUpper(validSignatureHeader(body, "secret"))
	if !VerifyWebhookSignature(body, header, "secret") {
		t.Fatal("expected uppercase hex signature to verify")
	}
}

func TestVerifyWebhookSignatureRejectsMalformedHex(t *testing.T) {
	if VerifyWebhookSignature([]byte(`{"event":"created"}`), "sha256=not-hex", "secret") {
		t.Fatal("expected malformed hex signature to fail")
	}
}

func TestVerifyWebhookSignatureRejectsWrongHeaderPrefix(t *testing.T) {
	body := []byte(`{"event":"created"}`)
	header := "sha1=" + hex.EncodeToString(computeWebhookDigest(body, "secret"))
	if VerifyWebhookSignature(body, header, "secret") {
		t.Fatal("expected wrong signature prefix to fail")
	}
}
