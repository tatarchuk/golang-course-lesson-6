// Package validate provides simple syntactic validators for common
// user-input formats.
//
// Homework — Task 2 (Lesson 6: File I/O, JSON and Testing):
// Implement ValidateEmail and/or ValidatePhone below (your mentor may
// ask for just one) and extend the test tables in validate_test.go to
// at least 8 cases each, including edge cases.
package validate

import (
	"strings"
	"unicode"
)

const (
	maxEmailLen       = 254
	maxLocalPartLen   = 64
	maxDomainLabelLen = 63

	localPartSpecials = "._%+-"

	minPhoneDigits = 10
	maxPhoneDigits = 15
)

// ValidateEmail reports whether s is a syntactically valid email address.
func ValidateEmail(s string) bool {
	if s == "" || len(s) > maxEmailLen {
		return false
	}

	for _, r := range s {
		if unicode.IsSpace(r) {
			return false
		}
	}

	if strings.Count(s, "@") != 1 {
		return false
	}
	at := strings.Index(s, "@")
	local := s[:at]
	domain := s[at+1:]

	return validLocalPart(local) && validDomain(domain)
}

// validLocalPart checks the part of the address before "@".
func validLocalPart(local string) bool {
	if local == "" || len(local) > maxLocalPartLen {
		return false
	}
	if hasBadDots(local) {
		return false
	}
	for _, r := range local {
		if !isASCIILetterOrDigit(r) && !strings.ContainsRune(localPartSpecials, r) {
			return false
		}
	}
	return true
}

// validDomain checks the part of the address after "@".
func validDomain(domain string) bool {
	if domain == "" || !strings.Contains(domain, ".") {
		return false
	}
	for _, label := range strings.Split(domain, ".") {
		if !validDomainLabel(label) {
			return false
		}
	}
	return true
}

func validDomainLabel(label string) bool {
	if label == "" || len(label) > maxDomainLabelLen {
		return false
	}
	if strings.HasPrefix(label, "-") || strings.HasSuffix(label, "-") {
		return false
	}
	for _, r := range label {
		if !isASCIILetterOrDigit(r) && r != '-' {
			return false
		}
	}
	return true
}

func hasBadDots(s string) bool {
	return strings.HasPrefix(s, ".") || strings.HasSuffix(s, ".") || strings.Contains(s, "..")
}

func isASCIILetterOrDigit(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}

// ValidatePhone reports whether s is a syntactically valid phone number.
func ValidatePhone(s string) bool {
	if s == "" {
		return false
	}

	if strings.TrimSpace(s) != s {
		return false
	}

	if s[0] == '+' {
		s = s[1:]
	}

	digits := 0
	for _, r := range s {
		switch {
		case r >= '0' && r <= '9':
			digits++
		case r == ' ' || r == '-' || r == '(' || r == ')':
			// allowed separators, just skip them
		default:
			// letters, dots, another "+", unicode, ...
			return false
		}
	}

	return digits >= minPhoneDigits && digits <= maxPhoneDigits
}
