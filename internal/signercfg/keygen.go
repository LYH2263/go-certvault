package signercfg

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

func generateEC() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}
