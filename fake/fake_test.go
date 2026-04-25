package fake

import (
	"regexp"
	"strings"
	"testing"
	"time"
	"unicode"
)

func TestPersonName(t *testing.T) {
	name := PersonData().Name()
	parts := strings.Split(name, " ")
	if len(parts) != 2 {
		t.Fatalf("expected 2 parts, got %d", len(parts))
	}
}

func TestPersonEmail(t *testing.T) {
	email := PersonData().Email()
	if !strings.Contains(email, "@") {
		t.Errorf("expected email to contain @, got %s", email)
	}
}

func TestPersonPhone(t *testing.T) {
	phone := PersonData().Phone()
	if !strings.HasPrefix(phone, "+1-") {
		t.Errorf("expected phone to start with +1-, got %s", phone)
	}
}

func TestPersonUsername(t *testing.T) {
	username := PersonData().Username()
	if len(username) < 3 {
		t.Errorf("expected username length >= 3, got %d", len(username))
	}
}

func TestPersonPassword(t *testing.T) {
	password := PersonData().Password()
	if len(password) != 12 {
		t.Errorf("expected password length 12, got %d", len(password))
	}
}

func TestPersonAge(t *testing.T) {
	age := PersonData().Age(18, 65)
	if age < 18 || age >= 65 {
		t.Errorf("expected 18 <= age < 65, got %d", age)
	}
}

func TestInternetEmail(t *testing.T) {
	email := InternetData().Email()
	if !strings.Contains(email, "@") {
		t.Errorf("expected email to contain @, got %s", email)
	}
}

func TestInternetUsername(t *testing.T) {
	username := InternetData().Username()
	if len(username) < 3 {
		t.Errorf("expected username length >= 3, got %d", len(username))
	}
}

func TestInternetPassword(t *testing.T) {
	password := InternetData().Password()
	if len(password) != 12 {
		t.Errorf("expected password length 12, got %d", len(password))
	}
}

func TestInternetIPv4(t *testing.T) {
	ip := InternetData().IPv4()
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		t.Errorf("expected 4 parts, got %d", len(parts))
	}
}

func TestInternetIPv6(t *testing.T) {
	ip := InternetData().IPv6()
	parts := strings.Split(ip, ":")
	if len(parts) != 8 {
		t.Errorf("expected 8 parts, got %d", len(parts))
	}
}

func TestInternetMacAddress(t *testing.T) {
	mac := InternetData().MacAddress()
	parts := strings.Split(mac, ":")
	if len(parts) != 6 {
		t.Errorf("expected 6 parts, got %d", len(parts))
	}
}

func TestInternetDomain(t *testing.T) {
	domain := InternetData().Domain()
	if !strings.Contains(domain, ".") {
		t.Errorf("expected domain to contain dot, got %s", domain)
	}
}

func TestInternetSlug(t *testing.T) {
	slug := InternetData().Slug()
	if !strings.Contains(slug, "-") {
		t.Errorf("expected slug to contain hyphen, got %s", slug)
	}
}

func TestCompanyName(t *testing.T) {
	name := CompanyData().Name()
	if name == "" {
		t.Error("expected non-empty company name")
	}
}

func TestCompanyNameWithSuffix(t *testing.T) {
	name := CompanyData().NameWithSuffix()
	found := false
	for _, suffix := range companySuffixes {
		if strings.HasSuffix(name, suffix) {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected company name to have a suffix")
	}
}

func TestCompanyIndustry(t *testing.T) {
	industry := CompanyData().Industry()
	if industry == "" {
		t.Error("expected non-empty industry")
	}
}

func TestCompanyEmail(t *testing.T) {
	email := CompanyData().Email()
	if !strings.Contains(email, "@") || !strings.Contains(email, "contact@") {
		t.Errorf("expected company email to contain contact@, got %s", email)
	}
}

func TestNumberInt(t *testing.T) {
	n := NumberData().Int(10, 20)
	if n < 10 || n >= 20 {
		t.Errorf("expected 10 <= n < 20, got %d", n)
	}
}

func TestNumberInt64(t *testing.T) {
	n := NumberData().Int64(1000, 2000)
	if n < 1000 || n >= 2000 {
		t.Errorf("expected 1000 <= n < 2000, got %d", n)
	}
}

func TestNumberFloat(t *testing.T) {
	f := NumberData().Float(0, 100, 2)
	if f < 0 || f > 100 {
		t.Errorf("expected 0 <= f <= 100, got %f", f)
	}
}

func TestNumberBool(t *testing.T) {
	b := NumberData().Bool()
	// Just check it doesn't panic
	_ = b
}

func TestTextWord(t *testing.T) {
	word := TextData().Word()
	if word == "" {
		t.Error("expected non-empty word")
	}
}

func TestTextSentence(t *testing.T) {
	sentence := TextData().Sentence(5)
	words := strings.Split(strings.TrimSuffix(sentence, "."), " ")
	if len(words) != 5 {
		t.Errorf("expected 5 words, got %d", len(words))
	}
	// Check first letter is uppercase
	runes := []rune(sentence)
	if !unicode.IsUpper(runes[0]) {
		t.Error("expected first letter to be uppercase")
	}
}

func TestTextParagraph(t *testing.T) {
	para := TextData().Paragraph(3)
	sentences := strings.Split(para, ". ")
	if len(sentences) < 2 {
		t.Errorf("expected multiple sentences, got %d", len(sentences))
	}
}

func TestTimeDate(t *testing.T) {
	date := TimeData().Date()
	if date.IsZero() {
		t.Error("expected non-zero date")
	}
}

func TestTimeDateBetween(t *testing.T) {
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
	date := TimeData().DateBetween(start, end)
	if date.Before(start) || date.After(end) {
		t.Errorf("date %v not between %v and %v", date, start, end)
	}
}

func TestTimeTimestamp(t *testing.T) {
	ts := TimeData().Timestamp()
	now := time.Now().Unix()
	// Should be within a year
	if ts < now-365*24*60*60 || ts > now {
		t.Errorf("timestamp %d seems out of range", ts)
	}
}

func TestTimeTimestampBetween(t *testing.T) {
	start := int64(0)
	end := int64(1000000)
	ts := TimeData().TimestampBetween(start, end)
	if ts < start || ts >= end {
		t.Errorf("timestamp %d not between %d and %d", ts, start, end)
	}
}

func TestTimeMonth(t *testing.T) {
	month := TimeData().Month()
	found := false
	for _, m := range months {
		if m == month {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("got unexpected month: %s", month)
	}
}

func TestTimeDayOfWeek(t *testing.T) {
	day := TimeData().DayOfWeek()
	found := false
	for _, d := range daysOfWeek {
		if d == day {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("got unexpected day: %s", day)
	}
}

func TestFileFileName(t *testing.T) {
	filename := FileData().FileName()
	if !strings.HasPrefix(filename, ".") && !strings.Contains(filename, ".") {
		t.Errorf("expected filename to have extension, got %s", filename)
	}
}

func TestFileFileNameWithUUID(t *testing.T) {
	filename := FileData().FileNameWithUUID()
	matched, _ := regexp.MatchString(`.*_\d+\.\w+`, filename)
	if !matched {
		t.Errorf("expected filename with UUID pattern, got %s", filename)
	}
}

func TestFileExtension(t *testing.T) {
	ext := FileData().Extension()
	if !strings.HasPrefix(ext, ".") {
		t.Errorf("expected extension to start with dot, got %s", ext)
	}
}

func TestFileMIMEType(t *testing.T) {
	mime := FileData().MIMEType()
	if !strings.Contains(mime, "/") {
		t.Errorf("expected MIME type to contain slash, got %s", mime)
	}
}

func TestFilePath(t *testing.T) {
	path := FileData().FilePath()
	if !strings.HasPrefix(path, "/") {
		t.Errorf("expected path to start with slash, got %s", path)
	}
}

func TestColorHex(t *testing.T) {
	hex := ColorData().Hex()
	matched, _ := regexp.MatchString(`^#[0-9a-fA-F]{6}$`, hex)
	if !matched {
		t.Errorf("expected hex color format, got %s", hex)
	}
}

func TestColorRGB(t *testing.T) {
	rgb := ColorData().RGB()
	matched, _ := regexp.MatchString(`^rgb\(\d+, \d+, \d+\)$`, rgb)
	if !matched {
		t.Errorf("expected RGB format, got %s", rgb)
	}
}

func TestColorRGBA(t *testing.T) {
	rgba := ColorData().RGBA()
	matched, _ := regexp.MatchString(`^rgba\(\d+, \d+, \d+, 0\.\d+\)$`, rgba)
	if !matched {
		t.Errorf("expected RGBA format, got %s", rgba)
	}
}

func TestColorName(t *testing.T) {
	name := ColorData().Name()
	found := false
	for _, c := range colorNames {
		if c == name {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("got unexpected color name: %s", name)
	}
}

func TestCurrencyCode(t *testing.T) {
	code := CurrencyData().Code()
	found := false
	for _, c := range currencyCodes {
		if c == code {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("got unexpected currency code: %s", code)
	}
}

func TestCurrencyName(t *testing.T) {
	name := CurrencyData().Name()
	found := false
	for _, c := range currencyNames {
		if c == name {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("got unexpected currency name: %s", name)
	}
}

func TestCurrencySymbol(t *testing.T) {
	symbol := CurrencyData().Symbol()
	found := false
	for _, s := range currencySymbols {
		if s == symbol {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("got unexpected currency symbol: %s", symbol)
	}
}

func TestCurrencyAmount(t *testing.T) {
	amount := CurrencyData().Amount(10, 100)
	if amount < 10 || amount > 100 {
		t.Errorf("expected 10 <= amount <= 100, got %f", amount)
	}
}

func TestLoremWord(t *testing.T) {
	word := LoremData().Word()
	if word == "" {
		t.Error("expected non-empty word")
	}
}

func TestLoremWords(t *testing.T) {
	words := LoremData().Words(10)
	count := len(strings.Split(words, " "))
	if count != 10 {
		t.Errorf("expected 10 words, got %d", count)
	}
}

func TestLoremSentence(t *testing.T) {
	sentence := LoremData().Sentence(8)
	runes := []rune(sentence)
	if !unicode.IsUpper(runes[0]) {
		t.Error("expected first letter to be uppercase")
	}
	if !strings.HasSuffix(sentence, ".") {
		t.Error("expected sentence to end with period")
	}
}

func TestLoremParagraph(t *testing.T) {
	para := LoremData().Paragraph(3)
	count := strings.Count(para, ".") + strings.Count(para, "!") + strings.Count(para, "?")
	if count < 2 {
		t.Errorf("expected multiple sentences, got %d", count)
	}
}
