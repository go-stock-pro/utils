package token

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/pkg/errors"
)

func derToPemString(key, keyType string) string {
	return fmt.Sprintf("-----BEGIN %s-----\n%s\n-----END %s-----", keyType, key, keyType)
}

func LoadPublicKey(publicKey string) (*rsa.PublicKey, error) {
	publicKey = derToPemString(publicKey, "PUBLIC KEY")

	var parsedKey interface{}
	dPem, _ := pem.Decode([]byte(publicKey))

	if dPem.Type == "RSA PUBLIC KEY" {
		return nil, errors.Errorf("invalid key type: %s", dPem.Type)
	}

	parsedKey, err := x509.ParsePKIXPublicKey(dPem.Bytes)
	if err != nil {
		return nil, err
	}

	if pubKey, ok := parsedKey.(*rsa.PublicKey); ok {
		return pubKey, nil
	}

	return nil, errors.Errorf("failed to parse rsa public key")
}

func LoadPrivateKey(privateKey, privateKeyPassword string) (priKey *rsa.PrivateKey, err error) {
	privateKey = derToPemString(privateKey, "PRIVATE KEY")
	privPem, _ := pem.Decode([]byte(privateKey))

	if privPem.Type == "RSA PRIVATE KEY" {
		return nil, errors.Errorf("invalid key type: %s", privPem.Type)
	}

	var privPemBytes []byte
	if privateKeyPassword != "" {
		privPemBytes, err = x509.DecryptPEMBlock(privPem, []byte(privateKeyPassword))
		if err != nil {
			return nil, err
		}
	} else {
		privPemBytes = privPem.Bytes
	}

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(privPemBytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(privPemBytes); err != nil {
			return nil, err
		}
	}

	if priKey, ok := parsedKey.(*rsa.PrivateKey); ok {
		return priKey, nil
	}

	return nil, errors.Errorf("failed to parse rsa private key")
}
