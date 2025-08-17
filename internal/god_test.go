package internal

import (
	"strings"
	"sync"
	"testing"
)

func TestNewGod(t *testing.T) {
	god, err := NewGod(32)
	if err != nil {
		t.Fatalf("Failed to create God instance: %v", err)
	}

	if god == nil {
		t.Fatal("God instance is nil")
	}

	if god.GetAmount() != 32 {
		t.Errorf("Expected amount to be 32, got %d", god.GetAmount())
	}

	if god.GetWordsCount() == 0 {
		t.Error("Expected words to be loaded, got 0 words")
	}
}

func TestGodSpeak(t *testing.T) {
	god, err := NewGod(5)
	if err != nil {
		t.Fatalf("Failed to create God instance: %v", err)
	}

	message := god.Speak()
	if message == "" {
		t.Error("Expected non-empty message")
	}

	// Note: The message contains 5 "entries" from Happy.TXT, but each entry
	// can contain multiple words (e.g., "Catastrophic Success", "I'll ask nicely")
	// So we can't test for exact word count, just that we have some words
	words := strings.Fields(message)
	if len(words) == 0 {
		t.Error("Expected at least some words in message")
	}

	// Test that we get different messages (with high probability)
	message2 := god.Speak()
	if message == message2 {
		// This could happen by chance, but it's very unlikely with enough words
		t.Logf("Got same message twice (this could happen by chance): %s", message)
	}
}

func TestGodSetAmount(t *testing.T) {
	god, err := NewGod(10)
	if err != nil {
		t.Fatalf("Failed to create God instance: %v", err)
	}

	err = god.SetAmount(20)
	if err != nil {
		t.Fatalf("Failed to set amount: %v", err)
	}

	if god.GetAmount() != 20 {
		t.Errorf("Expected amount to be 20, got %d", god.GetAmount())
	}
}

func TestGodSetAmountInvalid(t *testing.T) {
	god, err := NewGod(10)
	if err != nil {
		t.Fatalf("Failed to create God instance: %v", err)
	}

	// Test invalid amounts
	testCases := []int{0, -1, 1001, 5000}
	for _, amount := range testCases {
		err = god.SetAmount(amount)
		if err == nil {
			t.Errorf("Expected error for amount %d, got nil", amount)
		}
		if err != ErrInvalidAmount {
			t.Errorf("Expected ErrInvalidAmount for amount %d, got %v", amount, err)
		}
	}

	// Ensure original amount is unchanged
	if god.GetAmount() != 10 {
		t.Errorf("Expected amount to remain 10, got %d", god.GetAmount())
	}
}

func TestReadWords(t *testing.T) {
	words, err := readWords()
	if err != nil {
		t.Fatalf("Failed to read words: %v", err)
	}

	if len(words) == 0 {
		t.Error("Expected words to be loaded, got 0 words")
	}

	// Check that we have some expected words from Terry's list
	expectedWords := []string{"God", "Terry", "TempleOS", "CIA", "FBI"}
	wordMap := make(map[string]bool)
	for _, word := range words {
		wordMap[word] = true
	}

	foundCount := 0
	for _, expected := range expectedWords {
		if wordMap[expected] {
			foundCount++
		}
	}

	if foundCount == 0 {
		t.Error("Expected to find at least some known words from Terry's list")
	}
}

func TestNewGodInvalidAmount(t *testing.T) {
	testCases := []int{0, -1, 1001, 5000}
	for _, amount := range testCases {
		_, err := NewGod(amount)
		if err == nil {
			t.Errorf("Expected error for amount %d, got nil", amount)
		}
		if err != ErrInvalidAmount {
			t.Errorf("Expected ErrInvalidAmount for amount %d, got %v", amount, err)
		}
	}
}

func TestGodSpeakWithAmount(t *testing.T) {
	god, err := NewGod(DefaultAmount)
	if err != nil {
		t.Fatalf("Failed to create God instance: %v", err)
	}

	// Test valid amounts - we test that we get some output, not exact word count
	// because Happy.TXT contains phrases with multiple words
	testCases := []int{1, 5, 10}
	for _, amount := range testCases {
		message, err := god.SpeakWithAmount(amount)
		if err != nil {
			t.Errorf("Unexpected error for amount %d: %v", amount, err)
			continue
		}

		if message == "" {
			t.Errorf("Expected non-empty message for amount %d", amount)
		}

		// Count the number of phrases (separated by spaces between entries)
		// This is a rough check - we expect at least some content
		if len(strings.TrimSpace(message)) == 0 {
			t.Errorf("Expected non-empty trimmed message for amount %d", amount)
		}
	}

	// Test invalid amounts
	invalidCases := []int{0, -1, 1001}
	for _, amount := range invalidCases {
		_, err := god.SpeakWithAmount(amount)
		if err == nil {
			t.Errorf("Expected error for invalid amount %d, got nil", amount)
		}
	}

	// Ensure original amount is unchanged
	if god.GetAmount() != DefaultAmount {
		t.Errorf("Expected amount to remain %d, got %d", DefaultAmount, god.GetAmount())
	}
}

func TestGodConcurrency(t *testing.T) {
	god, err := NewGod(10)
	if err != nil {
		t.Fatalf("Failed to create God instance: %v", err)
	}

	const numGoroutines = 100
	const numOperations = 10

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Test concurrent Speak operations
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				message := god.Speak()
				if message == "" {
					t.Errorf("Got empty message")
				}
			}
		}()
	}

	wg.Wait()
}

func TestValidateAmount(t *testing.T) {
	validCases := []int{1, 10, 100, 500, 1000}
	for _, amount := range validCases {
		if err := validateAmount(amount); err != nil {
			t.Errorf("Expected no error for valid amount %d, got %v", amount, err)
		}
	}

	invalidCases := []int{0, -1, -100, 1001, 5000}
	for _, amount := range invalidCases {
		if err := validateAmount(amount); err == nil {
			t.Errorf("Expected error for invalid amount %d, got nil", amount)
		}
	}
}

func BenchmarkGodSpeak(b *testing.B) {
	god, err := NewGod(32)
	if err != nil {
		b.Fatalf("Failed to create God instance: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = god.Speak()
	}
}

func BenchmarkNewGod(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := NewGod(DefaultAmount)
		if err != nil {
			b.Fatalf("Failed to create God instance: %v", err)
		}
	}
}

func BenchmarkGodSpeakWithAmount(b *testing.B) {
	god, err := NewGod(DefaultAmount)
	if err != nil {
		b.Fatalf("Failed to create God instance: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := god.SpeakWithAmount(10)
		if err != nil {
			b.Fatalf("Failed to speak: %v", err)
		}
	}
}

func BenchmarkGodConcurrentSpeak(b *testing.B) {
	god, err := NewGod(DefaultAmount)
	if err != nil {
		b.Fatalf("Failed to create God instance: %v", err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = god.Speak()
		}
	})
}
