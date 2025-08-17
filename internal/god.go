// Package internal provides the core god says functionality.
package internal

import (
	"bufio"
	"embed"
	"errors"
	"math/rand"
	"strings"
	"sync"
	"time"
)

//go:embed Happy.TXT
var happyFS embed.FS

// God represents the god says functionality with thread-safe operations.
type God struct {
	words  []string
	amount int
	rng    *rand.Rand
	mu     sync.RWMutex // Protects amount field and RNG
}

const (
	MinAmount     = 1
	MaxAmount     = 1000
	DefaultAmount = 32
)

// ErrInvalidAmount is returned when an invalid amount is provided
var ErrInvalidAmount = errors.New("amount must be between 1 and 1000")

// NewGod creates a new God instance with the specified amount of words to generate.
func NewGod(amount int) (*God, error) {
	if err := validateAmount(amount); err != nil {
		return nil, err
	}

	words, err := readWords()
	if err != nil {
		return nil, err
	}

	// Create a new random source with current time as seed
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	return &God{
		words:  words,
		amount: amount,
		rng:    rng,
	}, nil
}

// validateAmount checks if the provided amount is within valid range
func validateAmount(amount int) error {
	if amount < MinAmount || amount > MaxAmount {
		return ErrInvalidAmount
	}
	return nil
}

// readWords reads the Happy.TXT file and returns a slice of words
func readWords() ([]string, error) {
	data, err := happyFS.ReadFile("Happy.TXT")
	if err != nil {
		return nil, err
	}

	var words []string
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			words = append(words, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return words, nil
}

// Speak generates a random message by selecting words from the word list.
func (g *God) Speak() string {
	if len(g.words) == 0 {
		return ""
	}

	g.mu.RLock()
	amount := g.amount
	g.mu.RUnlock()

	return g.generateMessage(amount)
}

// generateMessage generates a message with the specified amount of words
func (g *God) generateMessage(amount int) string {
	selectedWords := make([]string, 0, amount)
	for i := 0; i < amount; i++ {
		g.mu.Lock()
		randomIndex := g.rng.Intn(len(g.words))
		word := g.words[randomIndex]
		g.mu.Unlock()

		if word != "" {
			selectedWords = append(selectedWords, word)
		}
	}

	return strings.Join(selectedWords, " ")
}

// SpeakWithAmount generates a random message with a specific amount of words.
func (g *God) SpeakWithAmount(amount int) (string, error) {
	if err := validateAmount(amount); err != nil {
		return "", err
	}

	if len(g.words) == 0 {
		return "", nil
	}

	return g.generateMessage(amount), nil
}

// SetAmount sets the number of words to generate.
func (g *God) SetAmount(amount int) error {
	if err := validateAmount(amount); err != nil {
		return err
	}

	g.mu.Lock()
	g.amount = amount
	g.mu.Unlock()
	return nil
}

// GetAmount returns the current amount of words to generate.
func (g *God) GetAmount() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.amount
}

// GetWordsCount returns the total number of words available
func (g *God) GetWordsCount() int {
	return len(g.words)
}
