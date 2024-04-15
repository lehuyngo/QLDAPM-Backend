package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
)

func GetMD5Hash(text string) string {
    hasher := md5.New()
    hasher.Write([]byte(text))
    return hex.EncodeToString(hasher.Sum(nil))
}

var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func encode(b []byte) string {
    return base64.StdEncoding.EncodeToString(b)
}

func decode(s string) ([]byte, error) {
    return base64.StdEncoding.DecodeString(s)
}

func Encrypt(text, secretKey string) (string, error) {
    block, err := aes.NewCipher([]byte(secretKey))
    if err != nil {
        return "", err
    }
    
    plainText := []byte(text)
    cfb := cipher.NewCFBEncrypter(block, bytes)
    cipherText := make([]byte, len(plainText))
    cfb.XORKeyStream(cipherText, plainText)
    return encode(cipherText), nil
}

func Decrypt(text, secretKey string) (string, error) {
    block, err := aes.NewCipher([]byte(secretKey))
    if err != nil {
        return "", err
    }

    cipherText, err := decode(text)
    if err != nil {
        return "", err
    }

    cfb := cipher.NewCFBDecrypter(block, bytes)
    plainText := make([]byte, len(cipherText))
    cfb.XORKeyStream(plainText, cipherText)
    return string(plainText), nil
}
