// Homework — Task 2: extend emailCases and/or phoneCases below to at
// least 8 cases each (your mentor may ask for just one of the two
// functions), then implement validate.go until every subtest passes.
//
// Like todo_test.go, this file uses t.Run per case and t.Errorf (not
// t.Fatalf) for the actual assertions, so one wrong case never hides
// the others. Run `go test -v ./validate/...` and read every FAIL line.
package validate

import (
	"strings"
	"testing"
)

const minCases = 8

// emailCases is the table of test cases for ValidateEmail.
var emailCases = []struct {
	name  string
	input string
	want  bool
}{
	{"valid simple", "student@softserve.academy", true},
	{"missing at sign", "student-softserve.academy", false},
	{"empty string", "", false},

	// valid addresses
	{"valid with dot and plus in local part", "first.last+tag@softserve.academy", true},
	{"valid with subdomain", "student@mail.softserve.academy", true},
	{"valid with digits and hyphen", "student42@soft-serve.academy", true},
	{"valid with percent and underscore in local part", "stu_dent%1@softserve.academy", true},
	{"valid uppercase letters", "Student@SoftServe.Academy", true},
	{"valid local part of exactly 64 characters", strings.Repeat("a", 64) + "@softserve.academy", true},

	// missing or duplicated parts
	{"only at sign", "@", false},
	{"two at signs", "student@@softserve.academy", false},
	{"two at signs with text between", "student@soft@serve.academy", false},
	{"missing local part", "@softserve.academy", false},
	{"missing domain", "student@", false},
	{"domain without a dot", "student@localhost", false},

	// whitespace
	{"whitespace inside local part", "stu dent@softserve.academy", false},
	{"leading whitespace", " student@softserve.academy", false},
	{"trailing newline", "student@softserve.academy\n", false},

	// dots
	{"leading dot in local part", ".student@softserve.academy", false},
	{"trailing dot in local part", "student.@softserve.academy", false},
	{"trailing dot in domain", "student@softserve.academy.", false},
	{"consecutive dots in local part", "stu..dent@softserve.academy", false},
	{"consecutive dots in domain", "student@softserve..academy", false},

	// domain labels (the parts between the dots) - added after the AI blind
	// test in Task 3 found that "user@-example.com" was accepted
	{"leading hyphen in domain label", "student@-softserve.academy", false},
	{"trailing hyphen in domain label", "student@softserve-.academy", false},
	{"domain label of exactly 63 characters", "student@" + strings.Repeat("a", 63) + ".academy", true},
	{"domain label of 64 characters is too long", "student@" + strings.Repeat("a", 64) + ".academy", false},

	// unicode and length
	{"unicode in local part", "студент@softserve.academy", false},
	{"unicode in domain", "student@софтсерв.academy", false},
	{"local part of 65 characters is too long", strings.Repeat("a", 65) + "@softserve.academy", false},
	// 64 + 1 + 63 + 1 + 63 + 1 + 57 + 4 = 254 (each label stays within 63 characters)
	{"whole address of exactly 254 characters", strings.Repeat("a", 64) + "@" + strings.Repeat("b", 63) + "." + strings.Repeat("c", 63) + "." + strings.Repeat("d", 57) + ".com", true},
	{"whole address of 255 characters is too long", strings.Repeat("a", 64) + "@" + strings.Repeat("b", 63) + "." + strings.Repeat("c", 63) + "." + strings.Repeat("d", 58) + ".com", false},
}

func TestValidateEmail(t *testing.T) {
	if len(emailCases) < minCases {
		t.Fatalf(
			"emailCases has %d case(s), need at least %d — add more edge cases to validate_test.go before this test can pass",
			len(emailCases), minCases,
		)
	}

	for _, tc := range emailCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := ValidateEmail(tc.input)
			if got != tc.want {
				t.Errorf("ValidateEmail(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

// phoneCases is the table of test cases for ValidatePhone.
var phoneCases = []struct {
	name  string
	input string
	want  bool
}{
	{"valid with plus", "+380501234567", true},
	{"contains letters", "050-abc-4567", false},
	{"empty string", "", false},

	// valid formats
	{"valid local with dashes", "050-123-4567", true},
	{"valid local without separators", "0501234567", true},
	{"valid international with spaces and parentheses", "+38 (050) 123 45 67", true},
	{"valid with 15 digits (max length)", "+123456789012345", true},

	// wrong number of digits
	{"too short", "12345", false},
	{"9 digits is one too few", "050123456", false},
	{"16 digits is one too many", "+1234567890123456", false},
	{"only plus sign", "+", false},
	{"only separators", "- () -", false},

	// wrong characters / wrong position
	{"plus in the middle", "380+501234567", false},
	{"two plus signs", "++380501234567", false},
	{"dots as separators", "050.123.4567", false},
	{"leading whitespace", " +380501234567", false},
	{"trailing whitespace", "0501234567 ", false},
	{"unicode letters", "050-абв-4567", false},
}

func TestValidatePhone(t *testing.T) {
	if len(phoneCases) < minCases {
		t.Fatalf(
			"phoneCases has %d case(s), need at least %d — add more edge cases to validate_test.go before this test can pass",
			len(phoneCases), minCases,
		)
	}

	for _, tc := range phoneCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := ValidatePhone(tc.input)
			if got != tc.want {
				t.Errorf("ValidatePhone(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}
