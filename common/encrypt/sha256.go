package encrypt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

// Sha256Encrypt encrypts the given string with sha256.
func Sha256Encrypt(s string) string {
	secret := "tiktok89757"
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(s))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}
