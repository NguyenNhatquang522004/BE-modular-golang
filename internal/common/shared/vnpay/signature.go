package vnpay

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// ─────────────────────────────────────────────
// Core Signature Utilities
// ─────────────────────────────────────────────

// hmacSHA512 computes HMAC-SHA512 of data using secret and returns the
// lowercase hex-encoded digest.
func hmacSHA512(secret, data string) string {
	h := hmac.New(sha512.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// buildSignatureData constructs the canonical string that is signed / verified.
//
// VNPay signing rules:
//  1. Exclude "vnp_SecureHash" and "vnp_SecureHashType" from the parameter set.
//  2. Sort remaining parameter names lexicographically (ascending).
//  3. URL-encode each VALUE with standard query encoding.
//  4. Concatenate as "key1=value1&key2=value2&…".
func buildSignatureData(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		if k != "vnp_SecureHash" && k != "vnp_SecureHashType" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var sb strings.Builder
	for i, k := range keys {
		if i > 0 {
			sb.WriteByte('&')
		}
		sb.WriteString(k)
		sb.WriteByte('=')
		sb.WriteString(url.QueryEscape(params[k]))
	}
	return sb.String()
}

// Sign generates the HMAC-SHA512 signature for a VNPay parameter map.
// It does NOT modify the map; the caller should add the result as
// "vnp_SecureHash" before building the final URL.
func Sign(params map[string]string, hashSecret string) string {
	return hmacSHA512(hashSecret, buildSignatureData(params))
}

// VerifySignature validates the "vnp_SecureHash" value present in params.
// Returns false if the field is absent or the digest does not match.
//
// Uses constant-time comparison to prevent timing attacks.
func VerifySignature(params map[string]string, hashSecret string) bool {
	received, ok := params["vnp_SecureHash"]
	if !ok || received == "" {
		return false
	}
	expected := Sign(params, hashSecret)
	// Normalise to lowercase before constant-time comparison.
	return hmac.Equal(
		[]byte(strings.ToLower(expected)),
		[]byte(strings.ToLower(received)),
	)
}

// ─────────────────────────────────────────────
// URL Builders
// ─────────────────────────────────────────────

// BuildPaymentURL assembles the full VNPay payment URL from the base endpoint
// and a signed parameter map.
//
// It adds "vnp_SecureHashType" and "vnp_SecureHash" automatically, then
// appends all parameters as a raw (pre-encoded) query string so that VNPay
// receives exactly the bytes that were signed.
func BuildPaymentURL(baseURL string, params map[string]string, hashSecret string) (string, error) {
	if baseURL == "" {
		return "", fmt.Errorf("vnpay: payment base URL must not be empty")
	}

	// Compute signature before encoding into the URL.
	params["vnp_SecureHashType"] = SecureHashType
	params["vnp_SecureHash"] = Sign(params, hashSecret)

	// Sort params for deterministic, human-readable URLs (not required by VNPay).
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	sb.WriteString(baseURL)
	sb.WriteByte('?')

	for i, k := range keys {
		if i > 0 {
			sb.WriteByte('&')
		}
		sb.WriteString(url.QueryEscape(k))
		sb.WriteByte('=')
		sb.WriteString(url.QueryEscape(params[k]))
	}

	return sb.String(), nil
}

// ─────────────────────────────────────────────
// Return URL / IPN Helpers
// ─────────────────────────────────────────────

// ParseReturnQuery parses the raw query string appended by VNPay to the
// return / IPN URL into a flat map[string]string.
// Each key appears at most once; if duplicated the first value wins.
func ParseReturnQuery(rawQuery string) (map[string]string, error) {
	values, err := url.ParseQuery(rawQuery)
	if err != nil {
		return nil, fmt.Errorf("vnpay: failed to parse return query: %w", err)
	}
	result := make(map[string]string, len(values))
	for k, vs := range values {
		if len(vs) > 0 {
			result[k] = vs[0]
		}
	}
	return result, nil
}
