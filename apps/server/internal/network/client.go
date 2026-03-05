package network

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Config holds network security configuration
type Config struct {
	AllowHosts          []string
	BlockPrivateMetaIPs bool
	Timeout             time.Duration
}

// Client is an SSRF-safe HTTP client
type Client struct {
	config      Config
	allowCIDRs  []*net.IPNet
	allowHosts  map[string]bool
	blockCIDRs  []*net.IPNet
	httpClient  *http.Client
}

// NewClient creates a new SSRF-safe HTTP client
func NewClient(cfg Config) (*Client, error) {
	c := &Client{
		config:     cfg,
		allowHosts: make(map[string]bool),
	}

	// Parse allow hosts/CIDRs
	for _, h := range cfg.AllowHosts {
		if strings.Contains(h, "/") {
			_, cidr, err := net.ParseCIDR(h)
			if err != nil {
				return nil, fmt.Errorf("invalid CIDR %q: %w", h, err)
			}
			c.allowCIDRs = append(c.allowCIDRs, cidr)
		} else {
			c.allowHosts[h] = true
		}
	}

	// Build blocklist (always blocked)
	c.blockCIDRs = buildBlocklist(cfg.BlockPrivateMetaIPs)

	// Create HTTP client with custom transport
	c.httpClient = &http.Client{
		Timeout: cfg.Timeout,
		Transport: &http.Transport{
			DialContext:         c.dialContext,
			DisableCompression:  false,
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	return c, nil
}

// buildBlocklist creates the list of always-blocked networks
func buildBlocklist(blockPrivate bool) []*net.IPNet {
	blocked := []*net.IPNet{
		// Cloud metadata (always blocked)
		mustParseCIDR("169.254.169.254/32"),
		// Link-local (always blocked)
		mustParseCIDR("169.254.0.0/16"),
		// Loopback (always blocked)
		mustParseCIDR("127.0.0.0/8"),
		mustParseCIDR("::1/128"),
	}

	if blockPrivate {
		// Private networks (blocked when flag is set)
		blocked = append(blocked,
			mustParseCIDR("10.0.0.0/8"),
			mustParseCIDR("172.16.0.0/12"),
			mustParseCIDR("192.168.0.0/16"),
			mustParseCIDR("fc00::/7"), // IPv6 private
		)
	}

	return blocked
}

func mustParseCIDR(s string) *net.IPNet {
	_, cidr, err := net.ParseCIDR(s)
	if err != nil {
		panic(err)
	}
	return cidr
}

// dialContext wraps dialer with SSRF checks
func (c *Client) dialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, fmt.Errorf("invalid address: %w", err)
	}

	// Resolve hostname
	ips, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil {
		return nil, fmt.Errorf("DNS lookup failed: %w", err)
	}

	// Check all resolved IPs
	for _, ip := range ips {
		if err := c.checkIP(ip.IP); err != nil {
			return nil, fmt.Errorf("SSRF check failed for %s (%s): %w", host, ip.IP, err)
		}
	}

	// Use first IP for connection
	dialer := &net.Dialer{
		Timeout:   c.config.Timeout,
		KeepAlive: 30 * time.Second,
	}
	return dialer.DialContext(ctx, network, net.JoinHostPort(ips[0].IP.String(), port))
}

// checkIP verifies an IP is allowed
func (c *Client) checkIP(ip net.IP) error {
	// Check blocklist first
	for _, cidr := range c.blockCIDRs {
		if cidr.Contains(ip) {
			return fmt.Errorf("IP %s is in blocked range %s", ip, cidr)
		}
	}

	// If no allowlist, only localhost is allowed (for testing)
	if len(c.allowCIDRs) == 0 && len(c.allowHosts) == 0 {
		if !ip.IsLoopback() {
			return fmt.Errorf("no allowlist configured, only localhost allowed")
		}
		return nil
	}

	// Check allowlist
	for _, cidr := range c.allowCIDRs {
		if cidr.Contains(ip) {
			return nil
		}
	}

	return fmt.Errorf("IP %s is not in any allowed range", ip)
}

// Get performs an SSRF-safe HTTP GET
func (c *Client) Get(ctx context.Context, rawURL string) (*http.Response, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Check if hostname is explicitly allowed
	if len(c.allowHosts) > 0 {
		host := parsedURL.Hostname()
		if !c.allowHosts[host] {
			// Check if it resolves to an allowed IP
			ips, err := net.DefaultResolver.LookupIPAddr(ctx, host)
			if err != nil {
				return nil, fmt.Errorf("DNS lookup failed: %w", err)
			}
			
			allowed := false
			for _, ip := range ips {
				if c.checkIP(ip.IP) == nil {
					allowed = true
					break
				}
			}
			if !allowed {
				return nil, fmt.Errorf("host %q is not in allowlist", host)
			}
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}

	return c.httpClient.Do(req)
}

// Post performs an SSRF-safe HTTP POST
func (c *Client) Post(ctx context.Context, rawURL string, contentType string, body io.Reader) (*http.Response, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Check if hostname is explicitly allowed
	if len(c.allowHosts) > 0 {
		host := parsedURL.Hostname()
		if !c.allowHosts[host] {
			ips, err := net.DefaultResolver.LookupIPAddr(ctx, host)
			if err != nil {
				return nil, fmt.Errorf("DNS lookup failed: %w", err)
			}
			
			allowed := false
			for _, ip := range ips {
				if c.checkIP(ip.IP) == nil {
					allowed = true
					break
				}
			}
			if !allowed {
				return nil, fmt.Errorf("host %q is not in allowlist", host)
			}
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, rawURL, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)

	return c.httpClient.Do(req)
}

// Do executes an SSRF-safe HTTP request
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	if req.URL == nil {
		return nil, fmt.Errorf("request URL is nil")
	}

	// Verify URL is allowed
	host := req.URL.Hostname()
	if len(c.allowHosts) > 0 && !c.allowHosts[host] {
		// Check resolved IPs
		ips, err := net.DefaultResolver.LookupIPAddr(req.Context(), host)
		if err != nil {
			return nil, fmt.Errorf("DNS lookup failed: %w", err)
		}
		
		allowed := false
		for _, ip := range ips {
			if c.checkIP(ip.IP) == nil {
				allowed = true
				break
			}
		}
		if !allowed {
			return nil, fmt.Errorf("host %q is not in allowlist", host)
		}
	}

	return c.httpClient.Do(req)
}

// CheckHost verifies a host is allowed
func (c *Client) CheckHost(ctx context.Context, host string) error {
	ips, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil {
		return fmt.Errorf("DNS lookup failed: %w", err)
	}

	for _, ip := range ips {
		if err := c.checkIP(ip.IP); err == nil {
			return nil // At least one IP is allowed
		}
	}

	return fmt.Errorf("host %q has no allowed IP addresses", host)
}

// IsAllowed checks if an IP/CIDR string is allowed
func (c *Client) IsAllowed(ipOrCIDR string) bool {
	if strings.Contains(ipOrCIDR, "/") {
		_, cidr, err := net.ParseCIDR(ipOrCIDR)
		if err != nil {
			return false
		}
		// Check if this CIDR overlaps with any allowlist CIDR
		for _, allowed := range c.allowCIDRs {
			if cidrsOverlap(cidr, allowed) {
				return true
			}
		}
		return false
	}

	ip := net.ParseIP(ipOrCIDR)
	if ip == nil {
		return false
	}

	// Check if it's explicitly in allowHosts (as a hostname/IP string)
	if c.allowHosts[ipOrCIDR] {
		return true
	}

	// Check against allowCIDRs
	for _, cidr := range c.allowCIDRs {
		if cidr.Contains(ip) {
			return true
		}
	}

	return false
}

// cidrsOverlap checks if two CIDRs overlap
func cidrsOverlap(a, b *net.IPNet) bool {
	return a.Contains(b.IP) || b.Contains(a.IP)
}
