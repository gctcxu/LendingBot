package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"log"
	"os"
)

var logFile, _ = os.OpenFile("kucoin.log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0744)
var logger = log.New(logFile, "INFO ", log.Ldate|log.Ltime)

func SignWithHmac256(key string, payload string) string {

	hash := hmac.New(sha256.New, []byte(key))

	hash.Write([]byte(payload))

	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func SignWithHmac512ToHex(key string, payload string) string {

	hash := hmac.New(sha512.New, []byte(key))

	hash.Write([]byte(payload))

	return hex.EncodeToString(hash.Sum(nil))
}
