package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

var sercret = ""

func GenerateVerificationCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}
func GenerateHash(code string) string {
	mac := hmac.New(sha256.New, []byte(sercret))
	mac.Write([]byte(code))
	return hex.EncodeToString(mac.Sum(nil))
}

func CompareHash(code, hashed string) bool {
	codeHex := GenerateHash(code)

	a, errA := hex.DecodeString(codeHex)
	b, errB := hex.DecodeString(hashed)
	if errA != nil || errB != nil {
		return false
	}

	return hmac.Equal(a, b)
}
