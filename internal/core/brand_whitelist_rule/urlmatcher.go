package brand_whitelist_rule

import (
	"crypto/sha256"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// NormalizeURL normalizes a URL to a canonical lowercase form.
func NormalizeURL(raw string) string {
	if raw == "" {
		return ""
	}
	raw = strings.TrimSpace(raw)
	if !strings.Contains(raw, "://") {
		raw = "https://" + raw
	}
	u, err := url.Parse(raw)
	if err != nil {
		return strings.ToLower(raw)
	}
	u.Host = strings.ToLower(u.Host)
	u.Scheme = strings.ToLower(u.Scheme)
	return u.String()
}

// ExtractHostname extracts the lowercase hostname from a URL.
func ExtractHostname(raw string) string {
	if raw == "" {
		return ""
	}
	if !strings.Contains(raw, "://") {
		raw = "https://" + raw
	}
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	return strings.ToLower(u.Hostname())
}

// NormalizeDomain strips protocol, path, port, and leading "www." from a domain string.
func NormalizeDomain(raw string) string {
	if raw == "" {
		return ""
	}
	s := strings.TrimSpace(raw)
	if idx := strings.Index(s, "://"); idx >= 0 {
		s = s[idx+3:]
	}
	if idx := strings.Index(s, "/"); idx >= 0 {
		s = s[:idx]
	}
	if idx := strings.LastIndex(s, ":"); idx >= 0 {
		s = s[:idx]
	}
	s = strings.TrimPrefix(s, "www.")
	return strings.ToLower(s)
}

// IsSameOrSubdomain checks whether hostname is the same as or a subdomain of allowedDomain.
// It prevents bypass attacks like ".com.evilsite.com" matching ".com".
func IsSameOrSubdomain(hostname, allowedDomain string, allowSubdomains bool) bool {
	hostname = strings.ToLower(strings.TrimPrefix(hostname, "www."))
	allowedDomain = strings.ToLower(strings.TrimPrefix(allowedDomain, "www."))

	if hostname == allowedDomain {
		return true
	}
	if allowSubdomains {
		// Must end with "." + allowedDomain to be a legitimate subdomain.
		return strings.HasSuffix(hostname, "."+allowedDomain)
	}
	return false
}

// HashURL returns a short SHA-256 hex hash of the input string.
func HashURL(s string) string {
	h := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", h)[:16]
}

// matchRule evaluates whether targetURL satisfies the given whitelist rule.
func matchRule(rule *BrandWhitelistRule, targetURL string) (bool, error) {
	switch rule.MatchType {
	case "exact_url":
		return NormalizeURL(targetURL) == NormalizeURL(rule.Value), nil

	case "domain":
		hostname := ExtractHostname(targetURL)
		if hostname == "" {
			return false, nil
		}
		allowed := NormalizeDomain(rule.Value)
		return IsSameOrSubdomain(hostname, allowed, rule.AllowSubdomains), nil

	case "path_prefix":
		normalizedTarget := NormalizeURL(targetURL)
		normalizedRule := NormalizeURL(rule.Value)
		// Also require the hostname to match to prevent cross-domain prefix tricks.
		if ExtractHostname(targetURL) != ExtractHostname(rule.Value) {
			return false, nil
		}
		return strings.HasPrefix(normalizedTarget, normalizedRule), nil

	case "contains":
		// Safe only for non-domain contexts such as url_tags.
		return strings.Contains(targetURL, rule.Value), nil

	case "regex":
		r, err := regexp.Compile(rule.Value)
		if err != nil {
			return false, fmt.Errorf("invalid regex in rule %d: %w", rule.ID, err)
		}
		return r.MatchString(targetURL), nil
	}
	return false, nil
}
