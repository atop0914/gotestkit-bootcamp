// Package fake provides fake data generation for testing.
package fake

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
	"unicode"
)

var (
	globalRand *lockedRand
	once       sync.Once
)

func initRand() {
	once.Do(func() {
		globalRand = &lockedRand{rng: rand.New(rand.NewSource(time.Now().UnixNano()))}
	})
}

type lockedRand struct {
	mu sync.Mutex
	rng *rand.Rand
}

func (r *lockedRand) read(p []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	n, err = r.rng.Read(p)
	return n, err
}

func (r *lockedRand) Intn(n int) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.rng.Intn(n)
}

func (r *lockedRand) Int63() int64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.rng.Int63()
}

func (r *lockedRand) Uint32() uint32 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.rng.Uint32()
}

// Person generates fake person data
type Person struct{}

func (p *Person) Name() string {
	initRand()
	return fmt.Sprintf("%s %s",
		firstNames[globalRand.Intn(len(firstNames))],
		lastNames[globalRand.Intn(len(lastNames))],
	)
}

func (p *Person) FirstName() string {
	initRand()
	return firstNames[globalRand.Intn(len(firstNames))]
}

func (p *Person) LastName() string {
	initRand()
	return lastNames[globalRand.Intn(len(lastNames))]
}

func (p *Person) Email() string {
	initRand()
	domain := emailDomains[globalRand.Intn(len(emailDomains))]
	first := firstNames[globalRand.Intn(len(firstNames))]
	last := lastNames[globalRand.Intn(len(lastNames))]
	username := strings.ToLower(first + last)
	n := globalRand.Intn(999)
	return fmt.Sprintf("%s%d@%s", username, n, domain)
}

func (p *Person) Phone() string {
	initRand()
	return fmt.Sprintf("+1-%d-%d-%d",
		200+globalRand.Intn(800),
		100+globalRand.Intn(900),
		100+globalRand.Intn(9000),
	)
}

func (p *Person) Address() string {
	initRand()
	return fmt.Sprintf("%d %s %s, %s %s",
		100+globalRand.Intn(9900),
		streets[globalRand.Intn(len(streets))],
		cities[globalRand.Intn(len(cities))],
		states[globalRand.Intn(len(states))],
		zipCodes[globalRand.Intn(len(zipCodes))],
	)
}

// Internet generates fake internet data
type Internet struct{}

func (i *Internet) Email() string {
	return (&Person{}).Email()
}

func (i *Internet) Username() string {
	initRand()
	adj := adjectives[globalRand.Intn(len(adjectives))]
	noun := nouns[globalRand.Intn(len(nouns))]
	n := globalRand.Intn(99)
	return fmt.Sprintf("%s%s%d", adj, noun, n)
}

func (i *Internet) Password() string {
	initRand()
	var sb strings.Builder
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%"
	for j := 0; j < 12; j++ {
		sb.WriteByte(chars[globalRand.Intn(len(chars))])
	}
	return sb.String()
}

func (i *Internet) URL() string {
	initRand()
	subdomain := strings.ToLower(adjectives[globalRand.Intn(len(adjectives))])
	domain := domains[globalRand.Intn(len(domains))]
	return fmt.Sprintf("https://%s.%s", subdomain, domain)
}

func (i *Internet) IPv4() string {
	initRand()
	return fmt.Sprintf("%d.%d.%d.%d",
		1+globalRand.Intn(255),
		globalRand.Intn(256),
		globalRand.Intn(256),
		globalRand.Intn(256),
	)
}

func (i *Internet) IPv6() string {
	initRand()
	return fmt.Sprintf("%x:%x:%x:%x:%x:%x:%x:%x",
		globalRand.Uint32()&0xFFFF,
		globalRand.Uint32()&0xFFFF,
		globalRand.Uint32()&0xFFFF,
		globalRand.Uint32()&0xFFFF,
		globalRand.Uint32()&0xFFFF,
		globalRand.Uint32()&0xFFFF,
		globalRand.Uint32()&0xFFFF,
		globalRand.Uint32()&0xFFFF,
	)
}

func (i *Internet) MacAddress() string {
	initRand()
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
		globalRand.Uint32()&0xFF,
		globalRand.Uint32()&0xFF,
		globalRand.Uint32()&0xFF,
		globalRand.Uint32()&0xFF,
		globalRand.Uint32()&0xFF,
		globalRand.Uint32()&0xFF,
	)
}

// Number generates fake numbers
type Number struct{}

func (n *Number) Int(min, max int) int {
	initRand()
	return min + globalRand.Intn(max-min)
}

func (n *Number) Int64(min, max int64) int64 {
	initRand()
	return min + int64(globalRand.Int63()%int64(max-min))
}

func (n *Number) Float(min, max float64, decimals int) float64 {
	initRand()
	f := float64(globalRand.Int63()) / (1 << 63)
	f = f*(max-min) + min
	pow := 1
	for i := 0; i < decimals; i++ {
		pow *= 10
	}
	return float64(int(f*float64(pow))) / float64(pow)
}

func (n *Number) Bool() bool {
	initRand()
	return globalRand.Intn(2) == 1
}

// Text generates fake text
type Text struct{}

func (t *Text) Word() string {
	initRand()
	return words[globalRand.Intn(len(words))]
}

func (t *Text) Words(count int) string {
	initRand()
	var result []string
	for i := 0; i < count; i++ {
		result = append(result, words[globalRand.Intn(len(words))])
	}
	return strings.Join(result, " ")
}

func (t *Text) Sentence(wordCount int) string {
	initRand()
	ws := make([]string, wordCount)
	for i := range ws {
		w := words[globalRand.Intn(len(words))]
		runes := []rune(w)
		runes[0] = unicode.ToUpper(runes[0])
		ws[i] = string(runes)
	}
	return strings.Join(ws, " ") + "."
}

func (t *Text) Paragraph(sentenceCount int) string {
	initRand()
	var sentences []string
	for i := 0; i < sentenceCount; i++ {
		wordCount := 5 + globalRand.Intn(10)
		sentences = append(sentences, t.Sentence(wordCount))
	}
	return strings.Join(sentences, " ")
}

// PersonData generates fake person data
func PersonData() *Person { return &Person{} }

// Data arrays
var firstNames = []string{
	"James", "Mary", "John", "Patricia", "Robert", "Jennifer", "Michael", "Linda",
	"William", "Barbara", "David", "Elizabeth", "Richard", "Susan", "Joseph", "Jessica",
	"Thomas", "Sarah","Charles", "Karen", "Christopher", "Nancy", "Daniel", "Lisa",
	"Matthew", "Betty", "Anthony", "Margaret", "Mark", "Sandra", "Donald", "Ashley",
}

var lastNames = []string{
	"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis",
	"Rodriguez", "Martinez", "Hernandez", "Lopez", "Gonzalez", "Wilson", "Anderson",
	"Thomas", "Taylor", "Moore", "Jackson", "Martin", "Lee", "Perez", "Thompson",
	"White", "Harris", "Sanchez", "Clark", "Ramirez", "Lewis", "Robinson",
}

var emailDomains = []string{"gmail.com", "yahoo.com", "outlook.com", "example.com"}
var domains = []string{"com", "io", "dev", "net", "org"}
var streets = []string{
	"Main St", "Oak Ave", "Maple Dr", "Cedar Ln", "Pine Rd",
	"Washington Blvd", "Park Ave", "Lake Dr", "Hill Rd", "River St",
}
var cities = []string{
	"New York", "Los Angeles", "Chicago", "Houston", "Phoenix",
	"Philadelphia", "San Antonio", "San Diego", "Dallas", "San Jose",
}
var states = []string{"CA", "NY", "TX", "FL", "IL", "PA", "OH", "GA", "NC", "MI"}
var zipCodes = []string{"10001", "90001", "60601", "77001", "85001", "19101"}
var adjectives = []string{"quick", "bright", "bold", "calm", "eager", "fair", "kind", "wise"}
var nouns = []string{"fox", "bear", "wolf", "hawk", "lion", "tiger", "eagle", "deer"}
var words = []string{
	"apple", "banana", "cherry", "date", "elderberry",
	"fig", "grape", "honeydew", "kiwi", "lemon",
	"mango", "nectarine", "orange", "papaya", "quince",
	"raspberry", "strawberry", "tangerine", "watermelon", "xylophone",
	"yellow", "zebra", " zebra", "abandon", "ability", "able",
	"about", "above", "abroad", "absense", "absolute", "absorb",
}
