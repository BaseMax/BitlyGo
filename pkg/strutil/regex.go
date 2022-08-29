package strutil

import "regexp"

// RemoveNonAlphanumerical remove non alphanumerical characters from string
func RemoveNonAlphanumerical(str string) (string, error) {
	pattern, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return "", err
	}

	return pattern.ReplaceAllString(str, ""), nil
}
