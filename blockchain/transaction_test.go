package blockchain

import (
	"os"
	"reflect"
	"testing"
)

func TestAddTransaction(t *testing.T) {
	// Clean up before and after test
	os.Remove(BlockchainFile)
	defer os.Remove(BlockchainFile)
	
	bc := NewBlockChain()
	
	sender := "address1"
	receiver := "address2"
	amount := 50.0
	signature := "dummy_signature"

	index := bc.AddTransaction(sender, receiver, amount, signature)

	if index != 2 {
		t.Errorf("Expected next block index to be 2, got %d", index)
	}

	if len(bc.CurrentTransactions) != 1 {
		t.Errorf("Expected 1 pending transaction, got %d", len(bc.CurrentTransactions))
	}

	expectedTx := Transaction{
		Sender:    sender,
		Receiver:  receiver,
		Amount:    amount,
		Signature: signature,
	}

	if !reflect.DeepEqual(bc.CurrentTransactions[0], expectedTx) {
		t.Errorf("Transaction mismatch. Got %+v, want %+v", bc.CurrentTransactions[0], expectedTx)
	}
}

func TestMinedTransactions(t *testing.T) {
	// Clean up before and after test
	os.Remove(BlockchainFile)
	defer os.Remove(BlockchainFile)
	
	bc := NewBlockChain()
	
	bc.AddTransaction("A", "B", 10.0, "sig1")
	bc.AddTransaction("C", "D", 20.0, "sig2")
	
	latestBlock := bc.GetLatestBlock()
	proof := ProofOfWork(latestBlock.Proof)
	previousHash := latestBlock.CalculateHash()
	
	newBlock := bc.CreateBlock(proof, previousHash)
	
	if len(newBlock.Transactions) != 2 {
		t.Errorf("Expected 2 transactions in the new block, got %d", len(newBlock.Transactions))
	}
	
	if len(bc.CurrentTransactions) != 0 {
		t.Errorf("Expected mempool to be empty after mining, got %d", len(bc.CurrentTransactions))
	}
	
	if newBlock.Transactions[0].Sender != "A" || newBlock.Transactions[1].Sender != "C" {
		t.Errorf("Transactions order or data mismatch in block")
	}
}
