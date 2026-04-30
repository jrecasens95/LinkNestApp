package security

import (
	"context"
	"errors"
	"net"
	"net/netip"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	ErrURLTooLong           = errors.New("url is too long")
	ErrURLInvalid           = errors.New("url is invalid")
	ErrURLUnsupportedScheme = errors.New("url must start with http:// or https://")
	ErrURLBlockedHost       = errors.New("url host is not allowed")
	ErrURLPrivateAddress    = errors.New("url points to a private or internal address")
)

type URLValidator struct {
	maxLength          int
	blacklistedDomains map[string]struct{}
}

func NewURLValidator(maxLength int, blacklistedDomains []string) *URLValidator {
	if maxLength <= 0 {
		maxLength = 2048
	}

	blacklist := make(map[string]struct{}, len(blacklistedDomains))
	for _, domain := range blacklistedDomains {
		normalized := normalizeHost(domain)
		if normalized != "" {
			blacklist[normalized] = struct{}{}
		}
	}

	return &URLValidator{
		maxLength:          maxLength,
		blacklistedDomains: blacklist,
	}
}

func (v *URLValidator) ValidateCreateURL(rawURL string) error {
	parsed, err := v.parseAndValidate(rawURL)
	if err != nil {
		return err
	}

	host := normalizeHost(parsed.Hostname())
	if err := v.validateHost(host); err != nil {
		return err
	}

	return v.validateResolvedAddresses(host)
}

func (v *URLValidator) ValidateRedirectURL(rawURL string) error {
	parsed, err := v.parseAndValidate(rawURL)
	if err != nil {
		return err
	}

	return v.validateHost(normalizeHost(parsed.Hostname()))
}

func (v *URLValidator) parseAndValidate(rawURL string) (*url.URL, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" || !utf8.ValidString(rawURL) || strings.ContainsAny(rawURL, "\r\n\t\\") {
		return nil, ErrURLInvalid
	}
	if len(rawURL) > v.maxLength {
		return nil, ErrURLTooLong
	}

	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Host == "" || parsed.Hostname() == "" {
		return nil, ErrURLInvalid
	}
	if parsed.User != nil {
		return nil, ErrURLInvalid
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, ErrURLUnsupportedScheme
	}
	if parsed.Port() != "" {
		port, err := strconv.Atoi(parsed.Port())
		if err != nil || port < 1 || port > 65535 {
			return nil, ErrURLInvalid
		}
	}

	return parsed, nil
}

func (v *URLValidator) validateHost(host string) error {
	if host == "" {
		return ErrURLInvalid
	}

	if _, blocked := v.blacklistedDomains[host]; blocked {
		return ErrURLBlockedHost
	}
	for blocked := range v.blacklistedDomains {
		if strings.HasSuffix(host, "."+blocked) {
			return ErrURLBlockedHost
		}
	}

	if addr, err := netip.ParseAddr(host); err == nil {
		if isPrivateOrInternal(addr) {
			return ErrURLPrivateAddress
		}
		return nil
	}
	if !strings.Contains(host, ".") {
		return ErrURLBlockedHost
	}

	return nil
}

func (v *URLValidator) validateResolvedAddresses(host string) error {
	if _, err := netip.ParseAddr(host); err == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	addresses, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil {
		return ErrURLInvalid
	}

	for _, address := range addresses {
		addr, ok := netip.AddrFromSlice(address.IP)
		if !ok || isPrivateOrInternal(addr.Unmap()) {
			return ErrURLPrivateAddress
		}
	}

	return nil
}

func normalizeHost(host string) string {
	return strings.TrimSuffix(strings.ToLower(strings.TrimSpace(host)), ".")
}

func isPrivateOrInternal(addr netip.Addr) bool {
	return addr.IsLoopback() ||
		addr.IsPrivate() ||
		addr.IsLinkLocalUnicast() ||
		addr.IsLinkLocalMulticast() ||
		addr.IsMulticast() ||
		addr.IsUnspecified()
}
