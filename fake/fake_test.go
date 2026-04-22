package fake

import (
	"strings"
	"testing"
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

func TestInternetEmail(t *testing.T) {
	email := (&Internet{}).Email()
	if !strings.Contains(email, "@") {
		t.Errorf("expected email to contain @, got %s", email)
	}
}

func TestInternetUsername(t *testing.T) {
	username := (&Internet{}).Username()
	if len(username) < 3 {
		t.Errorf("expected username length >= 3, got %d", len(username))
	}
}

func TestInternetPassword(t *testing.T) {
	password := (&Internet{}).Password()
	if len(password) != 12 {
		t.Errorf("expected password length 12, got %d", len(password))
	}
}

func TestInternetIPv4(t *testing.T) {
	ip := (&Internet{}).IPv4()
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		t.Errorf("expected 4 parts, got %d", len(parts))
	}
}

func TestNumberInt(t *testing.T) {
	n := (&Number{}).Int(10, 20)
	if n < 10 || n >= 20 {
		t.Errorf("expected 10 <= n < 20, got %d", n)
	}
}

func TestNumberFloat(t *testing.T) {
	f := (&Number{}).Float(0, 100, 2)
	if f < 0 || f > 100 {
		t.Errorf("expected 0 <= f <= 100, got %f", f)
	}
}

func TestTextWord(t *testing.T) {
	word := (&Text{}).Word()
	if word == "" {
		t.Error("expected non-empty word")
	}
}

func TestTextSentence(t *testing.T) {
	sentence := (&Text{}).Sentence(5)
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
	para := (&Text{}).Paragraph(3)
	sentences := strings.Split(para, ". ")
	if len(sentences) < 2 {
		t.Errorf("expected multiple sentences, got %d", len(sentences))
	}
}
