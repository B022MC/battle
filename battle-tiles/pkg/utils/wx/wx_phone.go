package wx

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
)

type PhoneNumberInfo struct {
	PhoneNumber     string `json:"phoneNumber"`
	PurePhoneNumber string `json:"purePhoneNumber"`
	CountryCode     string `json:"countryCode"`
	Watermark       struct {
		AppID     string `json:"appid"`
		Timestamp int64  `json:"timestamp"`
	} `json:"watermark"`
}

// 解密手机号信息（微信获取手机号）iu斤斤计较急急急急急急急急急急急急vbbbbbbbb
func DecryptPhoneNumber(sessionKey, encryptedData, iv string) (*PhoneNumberInfo, error) {
	// base64 decode
	encryptedDataBytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, errors.New("encryptedData decode error: " + err.Error())
	}
	keyBytes, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return nil, errors.New("sessionKey decode error: " + err.Error())
	}
	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, errors.New("iv decode error: " + err.Error())
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, errors.New("new cipher error: " + err.Error())
	}
	if len(encryptedDataBytes)%block.BlockSize() != 0 {
		return nil, errors.New("encryptedData is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, ivBytes)
	plainText := make([]byte, len(encryptedDataBytes))
	mode.CryptBlocks(plainText, encryptedDataBytes)

	// remove PKCS#7 padding
	plainText = pkcs7Unpad(plainText)
	var phoneInfo PhoneNumberInfo
	if err := json.Unmarshal(plainText, &phoneInfo); err != nil {
		return nil, errors.New("unmarshal decrypted phone info error: " + err.Error())
	}
	return &phoneInfo, nil
}

func pkcs7Unpad(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	if unpadding > length {
		return data
	}
	return data[:(length - unpadding)]
}
