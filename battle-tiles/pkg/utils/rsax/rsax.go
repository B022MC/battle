package rsax

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
)

// ======== Key parse helpers ========

// Base64 -> *rsa.PublicKey (X.509 SubjectPublicKeyInfo / PKIX)
func parsePublicKeyFromBase64(b64 string) (*rsa.PublicKey, error) {
	der, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, fmt.Errorf("decode pub base64: %w", err)
	}
	pubAny, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		return nil, fmt.Errorf("parse PKIX pub: %w", err)
	}
	pub, ok := pubAny.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not RSA public key")
	}
	return pub, nil
}

// Base64 -> *rsa.PrivateKey (PKCS#1 优先, 否则尝试 PKCS#8)
func parsePrivateKeyFromBase64(b64 string) (*rsa.PrivateKey, error) {
	der, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, fmt.Errorf("decode pri base64: %w", err)
	}
	if pri1, err1 := x509.ParsePKCS1PrivateKey(der); err1 == nil {
		return pri1, nil
	}
	priAny, err := x509.ParsePKCS8PrivateKey(der)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}
	pri, ok := priAny.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("not RSA private key")
	}
	return pri, nil
}

func parsePublicKeyFromHex(h string) (*rsa.PublicKey, error) {
	der, err := hex.DecodeString(h)
	if err != nil {
		return nil, fmt.Errorf("decode pub hex: %w", err)
	}
	pubAny, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		return nil, fmt.Errorf("parse PKIX pub: %w", err)
	}
	pub, ok := pubAny.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not RSA public key")
	}
	return pub, nil
}

func parsePrivateKeyFromHex(h string) (*rsa.PrivateKey, error) {
	der, err := hex.DecodeString(h)
	if err != nil {
		return nil, fmt.Errorf("decode pri hex: %w", err)
	}
	if pri1, err1 := x509.ParsePKCS1PrivateKey(der); err1 == nil {
		return pri1, nil
	}
	priAny, err := x509.ParsePKCS8PrivateKey(der)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}
	pri, ok := priAny.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("not RSA private key")
	}
	return pri, nil
}

// ======== Encrypt / Decrypt (PKCS#1 v1.5) ========

// 兼容你的命名：公钥加密 -> Base64
func RsaEncryptBase64(plain []byte, pubBase64 string) (string, error) {
	return RsaEncryptToBase64(plain, pubBase64)
}

// 别名：公钥加密 -> Base64
func RsaEncryptToBase64(plain []byte, pubBase64 string) (string, error) {
	pub, err := parsePublicKeyFromBase64(pubBase64)
	if err != nil {
		return "", err
	}
	ct, err := rsa.EncryptPKCS1v15(rand.Reader, pub, plain)
	if err != nil {
		return "", fmt.Errorf("rsa encrypt: %w", err)
	}
	return base64.StdEncoding.EncodeToString(ct), nil
}

// 公钥加密 -> Hex
func RsaEncryptToHex(plain []byte, pubHex string) (string, error) {
	pub, err := parsePublicKeyFromHex(pubHex)
	if err != nil {
		return "", err
	}
	ct, err := rsa.EncryptPKCS1v15(rand.Reader, pub, plain)
	if err != nil {
		return "", fmt.Errorf("rsa encrypt: %w", err)
	}
	return hex.EncodeToString(ct), nil
}

// 兼容你的命名：私钥解密 <- Base64 密文
func RsaDecryptByBase64(cipherBase64, priBase64 string) ([]byte, error) {
	pri, err := parsePrivateKeyFromBase64(priBase64)
	if err != nil {
		return nil, err
	}
	ct, err := base64.StdEncoding.DecodeString(cipherBase64)
	if err != nil {
		return nil, fmt.Errorf("decode cipher base64: %w", err)
	}
	pt, err := rsa.DecryptPKCS1v15(rand.Reader, pri, ct)
	if err != nil {
		return nil, fmt.Errorf("rsa decrypt: %w", err)
	}
	return pt, nil
}

// 私钥解密 <- Hex 密文
func RsaDecryptByHex(cipherHex, priHex string) ([]byte, error) {
	pri, err := parsePrivateKeyFromHex(priHex)
	if err != nil {
		return nil, err
	}
	ct, err := hex.DecodeString(cipherHex)
	if err != nil {
		return nil, fmt.Errorf("decode cipher hex: %w", err)
	}
	pt, err := rsa.DecryptPKCS1v15(rand.Reader, pri, ct)
	if err != nil {
		return nil, fmt.Errorf("rsa decrypt: %w", err)
	}
	return pt, nil
}

// ======== Key generation (DER -> Base64/Hex) ========

type Base64Key struct {
	PublicKey  string // base64(SubjectPublicKeyInfo)
	PrivateKey string // base64(PKCS#1)
}

type HexKey struct {
	PublicKey  string // hex(SubjectPublicKeyInfo)
	PrivateKey string // hex(PKCS#1)
}

func GenerateRsaKeyBase64(bits int) (*Base64Key, error) {
	if bits < 512 || bits > 16384 {
		return nil, fmt.Errorf("invalid key size: %d", bits)
	}
	pri, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}
	pubDer, err := x509.MarshalPKIXPublicKey(&pri.PublicKey)
	if err != nil {
		return nil, err
	}
	priDer := x509.MarshalPKCS1PrivateKey(pri)

	return &Base64Key{
		PublicKey:  base64.StdEncoding.EncodeToString(pubDer),
		PrivateKey: base64.StdEncoding.EncodeToString(priDer),
	}, nil
}

func GenerateRsaKeyHex(bits int) (*HexKey, error) {
	k, err := GenerateRsaKeyBase64(bits)
	if err != nil {
		return nil, err
	}
	pubDer, err := base64.StdEncoding.DecodeString(k.PublicKey)
	if err != nil {
		return nil, err
	}
	priDer, err := base64.StdEncoding.DecodeString(k.PrivateKey)
	if err != nil {
		return nil, err
	}
	return &HexKey{
		PublicKey:  hex.EncodeToString(pubDer),
		PrivateKey: hex.EncodeToString(priDer),
	}, nil
}

// ======== Sign / Verify (SHA256 + PKCS#1 v1.5) ========

func RsaSignBase64(data []byte, priBase64 string) (string, error) {
	pri, err := parsePrivateKeyFromBase64(priBase64)
	if err != nil {
		return "", err
	}
	h := sha256.Sum256(data)
	sig, err := rsa.SignPKCS1v15(rand.Reader, pri, crypto.SHA256, h[:])
	if err != nil {
		return "", fmt.Errorf("rsa sign: %w", err)
	}
	return base64.StdEncoding.EncodeToString(sig), nil
}

func RsaVerifySignBase64(data []byte, signBase64, pubBase64 string) bool {
	pub, err := parsePublicKeyFromBase64(pubBase64)
	if err != nil {
		return false
	}
	sig, err := base64.StdEncoding.DecodeString(signBase64)
	if err != nil {
		return false
	}
	h := sha256.Sum256(data)
	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, h[:], sig) == nil
}

func RsaSignHex(data []byte, priHex string) (string, error) {
	pri, err := parsePrivateKeyFromHex(priHex)
	if err != nil {
		return "", err
	}
	h := sha256.Sum256(data)
	sig, err := rsa.SignPKCS1v15(rand.Reader, pri, crypto.SHA256, h[:])
	if err != nil {
		return "", fmt.Errorf("rsa sign: %w", err)
	}
	return hex.EncodeToString(sig), nil
}

func RsaVerifySignHex(data []byte, signHex, pubHex string) bool {
	pub, err := parsePublicKeyFromHex(pubHex)
	if err != nil {
		return false
	}
	sig, err := hex.DecodeString(signHex)
	if err != nil {
		return false
	}
	h := sha256.Sum256(data)
	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, h[:], sig) == nil
}
