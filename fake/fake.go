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

func (p *Person) Username() string {
	return (&Internet{}).Username()
}

func (p *Person) Password() string {
	return (&Internet{}).Password()
}

func (p *Person) Age(min, max int) int {
	return (&Number{}).Int(min, max)
}

// Company generates fake company data
type Company struct{}

func (c *Company) Name() string {
	initRand()
	return fmt.Sprintf("%s %s",
		companyAdjectives[globalRand.Intn(len(companyAdjectives))],
		companyNouns[globalRand.Intn(len(companyNouns))],
	)
}

func (c *Company) NameWithSuffix() string {
	initRand()
	suffix := companySuffixes[globalRand.Intn(len(companySuffixes))]
	return c.Name() + " " + suffix
}

func (c *Company) Industry() string {
	initRand()
	return industries[globalRand.Intn(len(industries))]
}

func (c *Company) Email() string {
	initRand()
	return fmt.Sprintf("contact@%s.%s",
		strings.ToLower(companyAdjectives[globalRand.Intn(len(companyAdjectives))]),
		domains[globalRand.Intn(len(domains))],
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

func (i *Internet) Domain() string {
	initRand()
	return fmt.Sprintf("%s.%s",
		strings.ToLower(adjectives[globalRand.Intn(len(adjectives))]),
		domains[globalRand.Intn(len(domains))],
	)
}

func (i *Internet) Slug() string {
	initRand()
	adj := adjectives[globalRand.Intn(len(adjectives))]
	noun := nouns[globalRand.Intn(len(nouns))]
	return strings.ToLower(adj) + "-" + strings.ToLower(noun)
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

// Time generates fake time data
type Time struct{}

func (t *Time) Date() time.Time {
	initRand()
	now := time.Now()
	offset := globalRand.Intn(365*24) * int(time.Hour)
	return now.Add(-time.Duration(offset))
}

func (t *Time) DateBetween(start, end time.Time) time.Time {
	initRand()
	duration := end.Sub(start)
	offset := time.Duration(globalRand.Int63() % int64(duration.Hours()))
	return start.Add(offset)
}

func (t *Time) Timestamp() int64 {
	return t.Date().Unix()
}

func (t *Time) TimestampBetween(start, end int64) int64 {
	initRand()
	return start + globalRand.Int63()%(end-start)
}

func (t *Time) Month() string {
	initRand()
	return months[globalRand.Intn(len(months))]
}

func (t *Time) DayOfWeek() string {
	initRand()
	return daysOfWeek[globalRand.Intn(len(daysOfWeek))]
}

// File generates fake file data
type File struct{}

func (f *File) FileName() string {
	initRand()
	ext := fileExtensions[globalRand.Intn(len(fileExtensions))]
	name := strings.ToLower(nouns[globalRand.Intn(len(nouns))])
	return fmt.Sprintf("%s%s", name, ext)
}

func (f *File) FileNameWithUUID() string {
	initRand()
	ext := fileExtensions[globalRand.Intn(len(fileExtensions))]
	return fmt.Sprintf("%s_%d%s",
		strings.ToLower(nouns[globalRand.Intn(len(nouns))]),
		globalRand.Int63(),
		ext,
	)
}

func (f *File) Extension() string {
	initRand()
	return fileExtensions[globalRand.Intn(len(fileExtensions))]
}

func (f *File) MIMEType() string {
	initRand()
	return mimeTypes[globalRand.Intn(len(mimeTypes))]
}

func (f *File) FilePath() string {
	initRand()
	return fmt.Sprintf("/%s/%s/%s",
		strings.ToLower(adjectives[globalRand.Intn(len(adjectives))]),
		strings.ToLower(nouns[globalRand.Intn(len(nouns))]),
		f.FileName(),
	)
}

// Color generates fake color data
type Color struct{}

func (c *Color) Hex() string {
	initRand()
	return fmt.Sprintf("#%06x", globalRand.Uint32()&0xFFFFFF)
}

func (c *Color) RGB() string {
	initRand()
	return fmt.Sprintf("rgb(%d, %d, %d)",
		globalRand.Intn(256),
		globalRand.Intn(256),
		globalRand.Intn(256),
	)
}

func (c *Color) RGBA() string {
	initRand()
	return fmt.Sprintf("rgba(%d, %d, %d, %.2f)",
		globalRand.Intn(256),
		globalRand.Intn(256),
		globalRand.Intn(256),
		float64(globalRand.Intn(100))/100.0,
	)
}

func (c *Color) Name() string {
	initRand()
	return colorNames[globalRand.Intn(len(colorNames))]
}

// Currency generates fake currency data
type Currency struct{}

func (c *Currency) Code() string {
	initRand()
	return currencyCodes[globalRand.Intn(len(currencyCodes))]
}

func (c *Currency) Name() string {
	initRand()
	return currencyNames[globalRand.Intn(len(currencyNames))]
}

func (c *Currency) Amount(min, max float64) float64 {
	return (&Number{}).Float(min, max, 2)
}

func (c *Currency) Symbol() string {
	initRand()
	return currencySymbols[globalRand.Intn(len(currencySymbols))]
}

// Lorem generates lorem ipsum text
type Lorem struct{}

func (l *Lorem) Word() string {
	initRand()
	return loremWords[globalRand.Intn(len(loremWords))]
}

func (l *Lorem) Words(count int) string {
	initRand()
	var result []string
	for i := 0; i < count; i++ {
		result = append(result, loremWords[globalRand.Intn(len(loremWords))])
	}
	return strings.Join(result, " ")
}

func (l *Lorem) Sentence(wordCount int) string {
	words := l.Words(wordCount)
	runes := []rune(words)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes) + "."
}

func (l *Lorem) Paragraph(sentenceCount int) string {
	initRand()
	var sentences []string
	for i := 0; i < sentenceCount; i++ {
		wordCount := 8 + globalRand.Intn(12)
		sentences = append(sentences, l.Sentence(wordCount))
	}
	return strings.Join(sentences, " ")
}

// PersonData generates fake person data
func PersonData() *Person { return &Person{} }

// CompanyData generates fake company data
func CompanyData() *Company { return &Company{} }

// InternetData generates fake internet data
func InternetData() *Internet { return &Internet{} }

// NumberData generates fake number data
func NumberData() *Number { return &Number{} }

// TextData generates fake text data
func TextData() *Text { return &Text{} }

// TimeData generates fake time data
func TimeData() *Time { return &Time{} }

// FileData generates fake file data
func FileData() *File { return &File{} }

// ColorData generates fake color data
func ColorData() *Color { return &Color{} }

// CurrencyData generates fake currency data
func CurrencyData() *Currency { return &Currency{} }

// LoremData generates fake lorem data
func LoremData() *Lorem { return &Lorem{} }

// Data arrays
var firstNames = []string{
	"James", "Mary", "John", "Patricia", "Robert", "Jennifer", "Michael", "Linda",
	"William", "Barbara", "David", "Elizabeth", "Richard", "Susan", "Joseph", "Jessica",
	"Thomas", "Sarah", "Charles", "Karen", "Christopher", "Nancy", "Daniel", "Lisa",
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
	"yellow", "zebra", "abandon", "ability", "able",
	"about", "above", "abroad", "absence", "absolute", "absorb",
}

var companyAdjectives = []string{"Global", "United", "Prime", "Elite", "Apex", "Nova", "Summit", "Pioneer"}
var companyNouns = []string{"Tech", "Systems", "Solutions", "Industries", "Corp", "Labs", "Dynamics", "Ventures"}
var companySuffixes = []string{"Inc", "LLC", "Ltd", "Corp", "Co", "Group", "Holdings", "Partners"}
var industries = []string{
	"Technology", "Healthcare", "Finance", "Retail", "Manufacturing",
	"Education", "Real Estate", "Transportation", "Energy", "Entertainment",
}

var months = []string{
	"January", "February", "March", "April", "May", "June",
	"July", "August", "September", "October", "November", "December",
}

var daysOfWeek = []string{
	"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday",
}

var fileExtensions = []string{".txt", ".pdf", ".doc", ".docx", ".jpg", ".png", ".gif", ".mp4", ".zip", ".csv"}
var mimeTypes = []string{
	"text/plain", "application/pdf", "image/jpeg", "image/png",
	"image/gif", "video/mp4", "application/zip", "text/csv",
}

var colorNames = []string{
	"red", "blue", "green", "yellow", "orange", "purple", "pink", "black", "white", "gray",
}

var currencyCodes = []string{"USD", "EUR", "GBP", "JPY", "CNY", "AUD", "CAD", "CHF", "INR", "BRL"}
var currencyNames = []string{
	"US Dollar", "Euro", "British Pound", "Japanese Yen", "Chinese Yuan",
	"Australian Dollar", "Canadian Dollar", "Swiss Franc", "Indian Rupee", "Brazilian Real",
}
var currencySymbols = []string{"$", "€", "£", "¥", "¥", "A$", "C$", "Fr", "₹", "R$"}

var loremWords = []string{
	"lorem", "ipsum", "dolor", "sit", "amet", "consectetur", "adipiscing", "elit",
	"sed", "do", "eiusmod", "tempor", "incididunt", "ut", "labore", "et", "dolore",
	"magna", "aliqua", "enim", "ad", "minim", "veniam", "quis", "nostrud",
	"exercitation", "ullamco", "laboris", "nisi", "aliquip", "ex", "ea", "commodo",
	"consequat", "duis", "aute", "irure", "in", "reprehenderit", "voluptate",
	"velit", "esse", "cillum", "fugiat", "nulla", "pariatur", "excepteur", "sint",
	"occaecat", "cupidatat", "non", "proident", "sunt", "culpa", "qui", "officia",
	"deserunt", "mollit", "anim", "id", "est", "laborum",
}
