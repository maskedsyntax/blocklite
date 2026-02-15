package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte
}

func NewWallet() *Wallet {
	private, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	public := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return &Wallet{private, public}
}

func (w *Wallet) GetAddress() string {
	return hex.EncodeToString(w.PublicKey)
}

func Sign(privateKey *ecdsa.PrivateKey, data string) (string, error) {
	hash := sha256.Sum256([]byte(data))
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return "", err
	}

	signature := append(r.Bytes(), s.Bytes()...)
	return hex.EncodeToString(signature), nil
}

func Verify(publicKeyHex string, data string, signatureHex string) bool {
	publicKeyBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return false
	}

	sigBytes, err := hex.DecodeString(signatureHex)
	if err != nil || len(sigBytes) < 64 {
		return false
	}

	rBytes := sigBytes[:len(sigBytes)/2]
	sBytes := sigBytes[len(sigBytes)/2:]

	r := new(big.Int).SetBytes(rBytes)
	s := new(big.Int).SetBytes(sBytes)

	xBytes := publicKeyBytes[:len(publicKeyBytes)/2]
	yBytes := publicKeyBytes[len(publicKeyBytes)/2:]

	x := new(big.Int).SetBytes(xBytes)
	y := new(big.Int).SetBytes(yBytes)

	rawPubKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	hash := sha256.Sum256([]byte(data))
	return ecdsa.Verify(&rawPubKey, hash[:], r, s)
}
