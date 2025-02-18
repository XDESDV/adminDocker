package functions

import (
	"regexp"
	"unicode"
)

// IsEmailValid checks if the email provided passes the required structure and length.
func IsEmailValid(s string) bool {
	return isValid(s,
		regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"))
}

// IsNameValid checks if the name provided respects the required elements
func IsNameValid(s string) bool {
	return isValid(s,
		regexp.MustCompile("^[a-zA-Z]+(([',. -][a-zA-Z ])?[a-zA-Z]*)*$"))
}

// IsUserNameValid checks if the name provided respects the required elements
func IsUserNameValid(s string) bool {
	return isValid(s,
		regexp.MustCompile("^[a-zA-Z]+(([',. -][a-zA-Z ])?[a-zA-Z]*)*$"))
}

// IsPasswordValid checks if the password provided respects the required elements
func IsPasswordValid(s string, i int) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if i == 0 {
		i = 7
	}
	if len(s) >= i {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

// IsLenStringValid check lenght of string
func IsLenStringValid(s string) bool {
	return (len(s) > 3 && len(s) <= 254)
}

// IsValidURL Check URL
func IsValidURL(s string) bool {
	return isValid(s,
		regexp.MustCompile(`^(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?$`))

}

// Called by otherr fonction to check
func isValid(s string, r *regexp.Regexp) bool {
	if !IsLenStringValid(s) {
		return false
	}
	return r.MatchString(s)
}
