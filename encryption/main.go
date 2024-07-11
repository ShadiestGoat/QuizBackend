package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)

var enc cipher.AEAD
var BlockSize = aes.BlockSize

func Init(secret string) error {
	c, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return err
	}
	
	enc = gcm
	return nil
}

func mkNonce() []byte {
	return bytes.Repeat([]byte{22}, enc.NonceSize())
}

func Encrypt(text string) string {
	dst := enc.Seal([]byte{}, mkNonce(), []byte(text), nil)
	return hex.EncodeToString(dst)
}

func Decrypt(secret string) string {
	b, err := hex.DecodeString(secret)
	if err != nil {
		return ""
	}

	dst, err := enc.Open([]byte{}, mkNonce(), b, nil)
	if err != nil {
		return ""
	}
	return string(dst)
}
