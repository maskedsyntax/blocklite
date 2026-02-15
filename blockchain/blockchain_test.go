package blockchain

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

// mockTime is a helper to mock time.Now for deterministic timestamps in tests.
func mockTime(t time.Time) func() {
	originalTimeNow := timeNow
	timeNow = func() time.Time { return t }
	return func() { timeNow = originalTimeNow }
}

// TestCreateBlock verifies that CreateBlock constructs a new block with the correct fields
// and appends it to the blockchain.
func TestCreateBlock(t *testing.T) {
	// Clean up before and after test
	os.Remove(BlockchainFile)
	defer os.Remove(BlockchainFile)

	// Test cases cover initial block creation and subsequent blocks with varied inputs.
	tests := []struct {
		name          string
		chain         []Block
		proof         int
		previousHash  string
		timestamp     time.Time
		expectedBlock Block
		expectedLen   int
	}{
		{
			name:         "First block",
			chain:        []Block{},
			proof:        100,
			previousHash: "0",
			timestamp:    time.Date(2025, 7, 6, 12, 0, 0, 0, time.UTC),
			expectedBlock: Block{
				Index:        1,
				Timestamp:    time.Date(2025, 7, 6, 12, 0, 0, 0, time.UTC).Format(time.RFC3339),
				Transactions: []Transaction{},
				Proof:        100,
				PreviousHash: "0",
			},
			expectedLen: 1,
		},
		{
			name: "Second block",
			chain: []Block{
				{Index: 1, Timestamp: "2025-07-06T12:00:00Z", Proof: 100, PreviousHash: "0"},
			},
			proof:        200,
			previousHash: "abc123",
			timestamp:    time.Date(2025, 7, 6, 13, 0, 0, 0, time.UTC),
			expectedBlock: Block{
				Index:        2,
				Timestamp:    time.Date(2025, 7, 6, 13, 0, 0, 0, time.UTC).Format(time.RFC3339),
				Transactions: []Transaction{},
				Proof:        200,
				PreviousHash: "abc123",
			},
			expectedLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock time.Now for deterministic timestamps.
			defer mockTime(tt.timestamp)()

			bc := &Blockchain{
				Chain:               tt.chain,
				CurrentTransactions: []Transaction{},
			}
			got := bc.CreateBlock(tt.proof, tt.previousHash)

			// Verify block fields.
			if !reflect.DeepEqual(got, tt.expectedBlock) {
				t.Errorf("CreateBlock() = %+v; want %+v", got, tt.expectedBlock)
			}

			// Verify chain length.
			if len(bc.Chain) != tt.expectedLen {
				t.Errorf("Chain length = %d; want %d", len(bc.Chain), tt.expectedLen)
			}

			// Verify the block was appended.
			if len(bc.Chain) > 0 && !reflect.DeepEqual(bc.Chain[len(bc.Chain)-1], tt.expectedBlock) {
				t.Errorf("Last block = %+v; want %+v", bc.Chain[len(bc.Chain)-1], tt.expectedBlock)
			}
		})
	}
}

// TestNewBlockChain validates that NewBlockChain initializes a blockchain with a genesis block.
func TestNewBlockChain(t *testing.T) {
	// Clean up before and after test
	os.Remove(BlockchainFile)
	defer os.Remove(BlockchainFile)
	
	// Test cases cover the creation of a new blockchain with a genesis block.
	tests := []struct {
		name          string
		timestamp     time.Time
		expectedBlock Block
	}{
		{
			name:      "Genesis block",
			timestamp: time.Date(2025, 7, 6, 12, 0, 0, 0, time.UTC),
			expectedBlock: Block{
				Index:        1,
				Timestamp:    time.Date(2025, 7, 6, 12, 0, 0, 0, time.UTC).Format(time.RFC3339),
				Transactions: []Transaction{},
				Proof:        1,
				PreviousHash: "0",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock time.Now for deterministic timestamps.
			defer mockTime(tt.timestamp)()

			bc := NewBlockChain()

			// Verify chain has one block.
			if len(bc.Chain) != 1 {
				t.Errorf("Chain length = %d; want 1", len(bc.Chain))
			}

			// Verify genesis block fields.
			if got := bc.Chain[0]; !reflect.DeepEqual(got, tt.expectedBlock) {
				t.Errorf("Genesis block = %+v; want %+v", got, tt.expectedBlock)
			}
		})
	}
}

// TestGetLatestBlock confirms that GetLatestBlock returns the last block and prints its details.
func TestGetLatestBlock(t *testing.T) {
	// Test cases cover single and multi-block chains, plus empty chain panic.
	tests := []struct {
		name        string
		chain       []Block
		expected    Block
		expectedOut string
		expectPanic bool
	}{
		{
			name: "Single block",
			chain: []Block{
				{Index: 1, Timestamp: "2025-07-06T12:00:00Z", Transactions: []Transaction{}, Proof: 1, PreviousHash: "0"},
			},
			expected:    Block{Index: 1, Timestamp: "2025-07-06T12:00:00Z", Transactions: []Transaction{}, Proof: 1, PreviousHash: "0"},
			expectedOut: "{Index: 1, Timestamp: 2025-07-06T12:00:00Z, Transactions: 0, Proof: 1, PreviousHash: 0}\n",
		},
		{
			name: "Multiple blocks",
			chain: []Block{
				{Index: 1, Timestamp: "2025-07-06T12:00:00Z", Transactions: []Transaction{}, Proof: 1, PreviousHash: "0"},
				{Index: 2, Timestamp: "2025-07-06T13:00:00Z", Transactions: []Transaction{}, Proof: 2, PreviousHash: "abc123"},
			},
			expected:    Block{Index: 2, Timestamp: "2025-07-06T13:00:00Z", Transactions: []Transaction{}, Proof: 2, PreviousHash: "abc123"},
			expectedOut: "{Index: 2, Timestamp: 2025-07-06T13:00:00Z, Transactions: 0, Proof: 2, PreviousHash: abc123}\n",
		},
		{
			name:        "Empty chain",
			chain:       []Block{},
			expectPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout output.
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("os.Pipe failed: %v", err)
			}
			originalStdout := os.Stdout
			os.Stdout = w
			t.Cleanup(func() {
				os.Stdout = originalStdout
				err := r.Close()
				if err != nil {
					return
				}
			})

			// Handle panic for empty chain.
			if tt.expectPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Error("GetLatestBlock did not panic on empty chain")
					}
				}()
			}

			// Execute GetLatestBlock.
			bc := &Blockchain{
				Chain:               tt.chain,
				CurrentTransactions: []Transaction{},
			}
			got := bc.GetLatestBlock()

			// Close writer to flush output.
			if err := w.Close(); err != nil {
				t.Fatalf("Failed to close writer: %v", err)
			}

			// Read captured output.
			var output bytes.Buffer
			if _, err := io.Copy(&output, r); err != nil {
				t.Fatalf("Failed to read output: %v", err)
			}

			// Skip output and block checks for panic case.
			if tt.expectPanic {
				return
			}

			// Verify returned block.
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("GetLatestBlock() = %+v; want %+v", got, tt.expected)
			}

			// Verify printed output.
			if gotOut := output.String(); gotOut != tt.expectedOut {
				t.Errorf("GetLatestBlock output = %q; want %q", gotOut, tt.expectedOut)
			}
		})
	}
}

// TestVerifyProof checks that VerifyProof correctly validates proof combinations
// based on the SHA-256 hash's leading zeros.
func TestVerifyProof(t *testing.T) {
	// Test cases cover valid and invalid proofs with precomputed hashes.
	tests := []struct {
		name        string
		proof       int
		lastProof   int
		expected    bool
		expectedHex string
	}{
		{
			name:        "Valid proof",
			proof:       93711,
			lastProof:   1,
			expected:    true,
			expectedHex: "0000e6b89b0f08f0f7faf26f0ce60185e5043f11ed84c8881033ea0580011111",
		},
		{
			name:        "Invalid proof",
			proof:       101,
			lastProof:   200,
			expected:    false,
			expectedHex: "918178f1c1d36996e59ca9fbe7e99e40dc8aa34c80f267766f8bf23cdb0a83c3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValid, gotHex := VerifyProof(tt.proof, tt.lastProof)
			if gotValid != tt.expected {
				t.Errorf("VerifyProof valid = %v; want %v", gotValid, tt.expected)
			}
			if !strings.EqualFold(gotHex, tt.expectedHex) {
				t.Errorf("VerifyProof hex = %q; want %q", gotHex, tt.expectedHex)
			}
		})
	}
}

// TestProofOfWork validates that ProofOfWork finds a proof producing a hash with four leading zeros.
func TestProofOfWork(t *testing.T) {
	// Test cases cover different lastProof values with precomputed proof and hash.
	tests := []struct {
		name       string
		lastProof  int
		expected   int
		hashOutput string
	}{
		{
			name:       "Basic proof",
			lastProof:  1,
			expected:   93711,
			hashOutput: "0000e6b89b0f08f0f7faf26f0ce60185e5043f11ed84c8881033ea0580011111",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout for hash output.
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("os.Pipe failed: %v", err)
			}
			originalStdout := os.Stdout
			os.Stdout = w
			t.Cleanup(func() {
				os.Stdout = originalStdout
				err := r.Close()
				if err != nil {
					return
				}
			})

			// Execute ProofOfWork.
			got := ProofOfWork(tt.lastProof)

			// Close writer to flush output.
			if err := w.Close(); err != nil {
				t.Fatalf("Failed to close writer: %v", err)
			}

			// Read captured output.
			var output bytes.Buffer
			if _, err := io.Copy(&output, r); err != nil {
				t.Fatalf("Failed to read output: %v", err)
			}

			// Verify proof result.
			if got != tt.expected {
				t.Errorf("ProofOfWork() = %d; want %d", got, tt.expected)
			}

			// Verify printed output.
			expectedOut := fmt.Sprintf("Encoded Hex Code: %s\n", tt.hashOutput)
			if gotOut := output.String(); gotOut != expectedOut {
				t.Errorf("ProofOfWork output = %q; want %q", gotOut, expectedOut)
			}
		})
	}
}

// TestIsChainValid confirms that IsChainValid correctly identifies valid and invalid blockchains.
func TestIsChainValid(t *testing.T) {
	// Test cases cover valid chains, invalid indices, hashes, proofs, and timestamps.
	tests := []struct {
		name     string
		chain    []Block
		expected bool
	}{
		{
			name: "Valid chain",
			chain: []Block{
				{Index: 1, Timestamp: "2025-07-06T12:00:00Z", Proof: 1, PreviousHash: "0"},
				{Index: 2, Timestamp: "2025-07-06T13:00:00Z", Proof: 93711, PreviousHash: "5828c510a157edfbb185531ea35578281374e2a407ac87dac1665b0d647423c0"},
			},
			expected: true,
		},
		{
			name: "Invalid index",
			chain: []Block{
				{Index: 1, Timestamp: "2025-07-06T12:00:00Z", Proof: 1, PreviousHash: "0"},
				{Index: 3, Timestamp: "2025-07-06T13:00:00Z", Proof: 52838, PreviousHash: "3a6cd2d1894f38cd050281be09ef785723fe45ddc291ed9f4b78d13f12571eb0"},
			},
			expected: false,
		},
		{
			name: "Invalid previous hash",
			chain: []Block{
				{Index: 1, Timestamp: "2025-07-06T12:00:00Z", Proof: 1, PreviousHash: "0"},
				{Index: 2, Timestamp: "2025-07-06T13:00:00Z", Proof: 52838, PreviousHash: "invalid"},
			},
			expected: false,
		},
		{
			name: "Invalid proof",
			chain: []Block{
				{Index: 1, Timestamp: "2025-07-06T12:00:00Z", Proof: 1, PreviousHash: "0"},
				{Index: 2, Timestamp: "2025-07-06T13:00:00Z", Proof: 101, PreviousHash: "3a6cd2d1894f38cd050281be09ef785723fe45ddc291ed9f4b78d13f12571eb0"},
			},
			expected: false,
		},
		{
			name: "Invalid timestamp",
			chain: []Block{
				{Index: 1, Timestamp: "2025-07-06T13:00:00Z", Proof: 1, PreviousHash: "0"},
				{Index: 2, Timestamp: "2025-07-06T12:00:00Z", Proof: 52838, PreviousHash: "3a6cd2d1894f38cd050281be09ef785723fe45ddc291ed9f4b78d13f12571eb0"},
			},
			expected: false,
		},
		{
			name:     "Single block",
			chain:    []Block{{Index: 1, Timestamp: "2025-07-06T12:00:00Z", Proof: 1, PreviousHash: "0"}},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bc := &Blockchain{
				Chain:               tt.chain,
				CurrentTransactions: []Transaction{},
			}
			got := bc.IsChainValid()
			if got != tt.expected {
				t.Errorf("IsChainValid() = %v; want %v", got, tt.expected)
			}
		})
		fmt.Println("-----------------------------------")
	}
}

// TestPrint validates that Print outputs all blocks in the chain using their Print method.
func TestBlockchainPrint(t *testing.T) {
	// Test cases cover empty, single, and multi-block chains.
	tests := []struct {
		name        string
		chain       []Block
		expectedOut string
	}{
		{
			name:        "Empty chain",
			chain:       []Block{},
			expectedOut: "",
		},
		{
			name: "Single block",
			chain: []Block{
				{Index: 1, Timestamp: "2025-07-06T12:00:00Z", Proof: 1, PreviousHash: "0"},
			},
			expectedOut: "{Index: 1, Timestamp: 2025-07-06T12:00:00Z, Transactions: 0, Proof: 1, PreviousHash: 0}\n",
		},
		{
			name: "Multiple blocks",
			chain: []Block{
				{Index: 1, Timestamp: "2025-07-06T12:00:00Z", Proof: 1, PreviousHash: "0"},
				{Index: 2, Timestamp: "2025-07-06T13:00:00Z", Proof: 2, PreviousHash: "abc123"},
			},
			expectedOut: "{Index: 1, Timestamp: 2025-07-06T12:00:00Z, Transactions: 0, Proof: 1, PreviousHash: 0}\n{Index: 2, Timestamp: 2025-07-06T13:00:00Z, Transactions: 0, Proof: 2, PreviousHash: abc123}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout output.
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("os.Pipe failed: %v", err)
			}
			originalStdout := os.Stdout
			os.Stdout = w
			t.Cleanup(func() {
				os.Stdout = originalStdout
				err := r.Close()
				if err != nil {
					return
				}
			})

			// Execute Print method.
			bc := &Blockchain{
				Chain:               tt.chain,
				CurrentTransactions: []Transaction{},
			}
			bc.Print()

			// Close writer to flush output.
			if err := w.Close(); err != nil {
				t.Fatalf("Failed to close writer: %v", err)
			}

			// Read captured output.
			var output bytes.Buffer
			if _, err := io.Copy(&output, r); err != nil {
				t.Fatalf("Failed to read output: %v", err)
			}

			// Verify printed output.
			if got := output.String(); got != tt.expectedOut {
				t.Errorf("Print output = %q; want %q", got, tt.expectedOut)
			}
		})
	}
}
