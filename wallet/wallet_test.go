package wallet

import (
	"testing"
)

func TestWallet(t *testing.T) {
	w := NewWallet()
	address := w.GetAddress()

	if len(address) == 0 {
		t.Error("Address should not be empty")
	}

	data := "some transaction data"
	signature, err := Sign(w.PrivateKey, data)
	if err != nil {
		t.Fatalf("Failed to sign: %v", err)
	}

	if !Verify(address, data, signature) {
		t.Error("Verification failed for valid signature")
	}

	if Verify(address, "wrong data", signature) {
		t.Error("Verification should fail for wrong data")
	}

	w2 := NewWallet()
	if Verify(w2.GetAddress(), data, signature) {
		t.Error("Verification should fail for wrong public key")
	}
}
