package utils

import (
	"golang.org/x/text/unicode/norm"
	"regexp"
	"strings"
	"unicode"
)

// GenerateSlug creates a URL-safe slug from a given string.
func GenerateSlug(value string) string {
	slug := strings.ToLower(value)

	// Replace accented characters with their ASCII equivalent
	slug = removeAccents(slug)

	// Replace all non-alphanumeric characters with a hyphen
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Trim hyphens from beginning and end
	slug = strings.Trim(slug, "-")

	return slug
}

// input: Café Déjà Vu output: Cafe Deja Vu
// input: Đà Nẵng output: Da Nang
func removeAccents(input string) string {
	// Normalize the input using NFD (decomposes characters into base + accent)
	t := norm.NFD.String(input)

	var b strings.Builder
	b.Grow(len(t))

	for _, r := range t {
		// Skip non-spacing marks
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}
