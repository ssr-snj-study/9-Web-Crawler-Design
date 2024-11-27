package crawler

import (
	"net/url"
	"strings"
	"unicode"
)

func SplitDomain(link string) (string, string) {
	u, err := url.Parse(link)
	if err != nil {
		return "", ""
	}
	return u.Host, u.Path
}

func DomainToCamel(domain string) string {
	domain, _ = SplitDomain(domain)
	domainParts := strings.Split(domain, ".")

	for i := range domainParts {
		domainParts[i] = strings.ToUpper(domainParts[i][:1]) + domainParts[i][1:]
	}
	return strings.Join(append([]string{"domain"}, domainParts...), "")
}

func CamelToDomain(camel string) string {
	var word []rune

	for _, r := range camel {
		if unicode.IsUpper(r) {
			word = append(word, '.')
		}
		word = append(word, r)
	}
	return strings.ToLower(string(word[1:]))
}
