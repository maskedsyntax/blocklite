package blockchain

import (
	"os"
	"testing"
)

func TestPersistence(t *testing.T) {
	filename := "test_blockchain.json"
	defer os.Remove(filename)

	bc := &Blockchain{
		Chain: []Block{},
	}
	bc.CreateBlock(1, "0")
	bc.AddTransaction("sender", "receiver", 100.0, "sig")
	bc.CreateBlock(2, bc.Chain[0].CalculateHash())

	err := bc.Save(filename)
	if err != nil {
		t.Fatalf("Failed to save blockchain: %v", err)
	}

	bc2 := &Blockchain{}
	err = bc2.LoadFromFile(filename)
	if err != nil {
		t.Fatalf("Failed to load blockchain: %v", err)
	}

	if len(bc2.Chain) != len(bc.Chain) {
		t.Errorf("Chain length mismatch: got %d, want %d", len(bc2.Chain), len(bc.Chain))
	}

	if bc2.Chain[1].Transactions[0].Amount != 100.0 {
		t.Errorf("Transaction data mismatch after loading")
	}
}
