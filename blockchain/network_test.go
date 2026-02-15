package blockchain

import (
	"os"
	"testing"
)

func TestRegisterNode(t *testing.T) {
	os.Remove(BlockchainFile)
	defer os.Remove(BlockchainFile)

	bc := NewBlockChain()
	node := "localhost:5000"
	bc.RegisterNode(node)

	if !bc.Nodes[node] {
		t.Errorf("Node %s was not registered", node)
	}
}

func TestValidChain(t *testing.T) {
	os.Remove(BlockchainFile)
	defer os.Remove(BlockchainFile)

	bc := NewBlockChain()
	
	// Create some blocks
	proof1 := ProofOfWork(bc.Chain[0].Proof)
	bc.CreateBlock(proof1, bc.Chain[0].CalculateHash())
	
	proof2 := ProofOfWork(bc.Chain[1].Proof)
	bc.CreateBlock(proof2, bc.Chain[1].CalculateHash())
	
	if !bc.ValidChain(bc.Chain) {
		t.Error("ValidChain failed for a valid chain")
	}
	
	// Corrupt the chain
	bc.Chain[1].PreviousHash = "corrupted"
	if bc.ValidChain(bc.Chain) {
		t.Error("ValidChain passed for a corrupted chain")
	}
}
