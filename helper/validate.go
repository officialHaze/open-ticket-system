package helper

import (
	"fmt"
	"regexp"
	"strings"
)

func ValidateEmail(email string) error {
	emailregex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if strings.TrimSpace(email) == "" {
		return fmt.Errorf("no email id provided")
	}

	if !emailregex.MatchString(email) {
		return fmt.Errorf("invalid email - %s", email)
	}

	return nil
}

func ValidatePhone(phone string) error {
	formattingRegex := regexp.MustCompile(`[()\-\.\s]`)  // + is intentionally NOT removed
	e164Regex := regexp.MustCompile(`^\+[1-9]\d{1,14}$`) // total 2–15 digits after +

	if strings.TrimSpace(phone) == "" {
		return fmt.Errorf("no phone number provided")
	}

	// Step 1: Remove only formatting chars (spaces, dashes, dots, parentheses) — keep +
	cleaned := formattingRegex.ReplaceAllString(phone, "")

	// Step 2: If there's no +, but it starts with 1–3 digits that could be a country code, try to fix
	if !strings.HasPrefix(cleaned, "+") {
		// Remove leading zeros
		cleaned = regexp.MustCompile(`^0+`).ReplaceAllString(cleaned, "")
		if cleaned != "" {
			cleaned = "+" + cleaned
		}
	}

	// Step 3: Final E.164 validation: + followed by 1–15 digits, no leading zero in country code
	if e164Regex.MatchString(cleaned) && len(cleaned) <= 16 { // + + up to 15 digits
		return nil
	}

	return fmt.Errorf("invalid phone no - %s", phone)
}
