package pemutil

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"
)

// ErrParse 解析失败哨兵，可用 errors.Is。
var ErrParse = errors.New("pemutil: parse failed")

// Pair 解析结果。
type Pair struct {
	Cert        *x509.Certificate
	Key         crypto.PrivateKey
	Subject     string
	Issuer      string
	Serial      string
	NotBefore   time.Time
	NotAfter    time.Time
	Fingerprint string
}

// ParsePair 解析证书与私钥 PEM，并用 %w 包装 ErrParse。
func ParsePair(certPEM, keyPEM []byte) (Pair, error) {
	cert, err := ParseCertificate(certPEM)
	if err != nil {
		return Pair{}, fmt.Errorf("pemutil: certificate: %v", err)
	}
	key, err := ParsePrivateKey(keyPEM)
	if err != nil {
		return Pair{}, fmt.Errorf("pemutil: private key: %v", err)
	}
	if err := MatchKey(cert, key); err != nil {
		return Pair{}, fmt.Errorf("pemutil: %v", err)
	}
	return Pair{
		Cert:        cert,
		Key:         key,
		Subject:     cert.Subject.String(),
		Issuer:      cert.Issuer.String(),
		Serial:      serialString(cert.SerialNumber),
		NotBefore:   cert.NotBefore,
		NotAfter:    cert.NotAfter,
		Fingerprint: Fingerprint(cert.Raw),
	}, nil
}

// ParseCertificate 解析首个 CERTIFICATE 块。
func ParseCertificate(pemBytes []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("no pem block")
	}
	if block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("unexpected type %q", block.Type)
	}
	return x509.ParseCertificate(block.Bytes)
}

// ParsePrivateKey 解析 PKCS8 / PKCS1 / EC 私钥。
func ParsePrivateKey(pemBytes []byte) (crypto.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("no pem block")
	}
	if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	if key, err := x509.ParseECPrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	return nil, errors.New("unsupported private key")
}

// MatchKey 检查公私钥配对。
func MatchKey(cert *x509.Certificate, key crypto.PrivateKey) error {
	switch k := key.(type) {
	case *rsa.PrivateKey:
		pub, ok := cert.PublicKey.(*rsa.PublicKey)
		if !ok || pub.N.Cmp(k.N) != 0 {
			return errors.New("rsa key mismatch")
		}
	case *ecdsa.PrivateKey:
		pub, ok := cert.PublicKey.(*ecdsa.PublicKey)
		if !ok || pub.X.Cmp(k.X) != 0 || pub.Y.Cmp(k.Y) != 0 {
			return errors.New("ecdsa key mismatch")
		}
	default:
		return errors.New("unsupported key type")
	}
	return nil
}

// EncodeCertificate PEM 编码证书。
func EncodeCertificate(der []byte) []byte {
	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
}

// EncodePKCS8Key PEM 编码 PKCS8 私钥。
func EncodePKCS8Key(key crypto.PrivateKey) ([]byte, error) {
	der, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}), nil
}

// Fingerprint SHA256 hex。
func Fingerprint(der []byte) string {
	sum := sha256.Sum256(der)
	return hex.EncodeToString(sum[:])
}

func serialString(n *big.Int) string {
	if n == nil {
		return ""
	}
	return strings.ToUpper(n.Text(16))
}
