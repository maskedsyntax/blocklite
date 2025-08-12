package blockchain

import (
	"bytes"
	"io"
	"os"
	"testing"
)

// TestCalculateHash verifies that the CalculateHash method produces the expected SHA-256 hash
// for various block configurations, ensuring correct concatenation and hashing of block fields.
func TestCalculateHash(t *testing.T) {
	// Test cases cover standard scenarios, zero proof, and empty previous hash for robustness.
	tests := []struct {
		name     string
		block    Block
		expected string
	}{
		{
			name: "Standard block",
			block: Block{
				Index:        1,
				Timestamp:    "2025-07-06T12:00:00Z",
				Proof:        12345,
				PreviousHash: "0",
			},
			// Expected hash computed for "12025-07-06T12:00:00Z123450" using SHA-256
			expected: "3a6cd2d1894f38cd050281be09ef785723fe45ddc291ed9f4b78d13f12571eb0",
		},
		{
			name: "Zero proof",
			block: Block{
				Index:        2,
				Timestamp:    "2025-07-06T13:00:00Z",
				Proof:        0,
				PreviousHash: "abc123",
			},
			// Expected hash computed for "22025-07-06T13:00:00Z0abc123" using SHA-256
			expected: "dbc830895f2083649ffbe715d9aeae0216c597003655d63802f0da02b7e0e4c2",
		},
		{
			name: "Empty previous hash",
			block: Block{
				Index:        3,
				Timestamp:    "2025-07-06T14:00:00Z",
				Proof:        999,
				PreviousHash: "",
			},
			// Expected hash computed for "32025-07-06T14:00:00Z999" using SHA-256
			expected: "c507012b088a24d63be367859a1c51f619a2e1c045771581339b0cdba87cb095",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.block.CalculateHash()
			if got != tt.expected {
				t.Errorf("CalculateHash() = %q; want %q", got, tt.expected)
			}
		})
	}
}

// TestPrint validates the Print method formats block details correctly, capturing and validating
// stdout output for various block configurations.
func TestPrint(t *testing.T) {
	// Test cases include standard blocks, zero proof, and empty previous hash to verify formatting.
	tests := []struct {
		name     string
		block    Block
		expected string
	}{
		{
			name: "Standard block",
			block: Block{
				Index:        1,
				Timestamp:    "2025-07-06T12:00:00Z",
				Proof:        12345,
				PreviousHash: "0",
			},
			expected: "{Index: 1, Timestamp: 2025-07-06T12:00:00Z, Proof: 12345, PreviousHash: 0}\n",
		},
		{
			name: "Zero proof",
			block: Block{
				Index:        2,
				Timestamp:    "2025-07-06T13:00:00Z",
				Proof:        0,
				PreviousHash: "abc123",
			},
			expected: "{Index: 2, Timestamp: 2025-07-06T13:00:00Z, Proof: 0, PreviousHash: abc123}\n",
		},
		{
			name: "Empty previous hash",
			block: Block{
				Index:        3,
				Timestamp:    "2025-07-06T14:00:00Z",
				Proof:        999,
				PreviousHash: "",
			},
			expected: "{Index: 3, Timestamp: 2025-07-06T14:00:00Z, Proof: 999, PreviousHash: }\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up pipe to capture stdout output
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("os.Pipe failed: %v", err)
			}

			// Redirect stdout, restoring it after test
			originalStdout := os.Stdout
			os.Stdout = w
			t.Cleanup(func() {
				os.Stdout = originalStdout
				err := r.Close()
				if err != nil {
					return
				}
			})

			// Execute Print method
			tt.block.Print()

			// Close the writer to flush the pipe
			if err := w.Close(); err != nil {
				t.Fatalf("Failed to close writer: %v", err)
			}

			// Read the output
			var output bytes.Buffer
			if _, err := io.Copy(&output, r); err != nil {
				t.Fatalf("Failed to read output: %v", err)
			}

			// Compare actual output with expected
			if got := output.String(); got != tt.expected {
				t.Errorf("Print() = %q; want %q", got, tt.expected)
			}
		})
	}
}
