package dnsbl

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"blacklist-check/pkg/models"
)

// Checker represents a DNSBL checker
type Checker struct {
	config     *models.Config
	resultChan chan models.BlacklistResult
}

// NewChecker creates a new DNSBL checker
func NewChecker(config *models.Config) *Checker {
	return &Checker{
		config:     config,
		resultChan: make(chan models.BlacklistResult, 256),
	}
}

// reverseIP reverses the IP address for DNSBL lookup
func reverseIP(ip string) string {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return ""
	}
	return fmt.Sprintf("%s.%s.%s.%s", parts[3], parts[2], parts[1], parts[0])
}

// CheckIP checks if an IP is blacklisted
func (c *Checker) CheckIP(ip string, wg *sync.WaitGroup, sem chan struct{}) {
	defer wg.Done()
	sem <- struct{}{}
	defer func() { <-sem }()

	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{Timeout: time.Duration(c.config.Timeout) * time.Second}
			return d.DialContext(ctx, network, address)
		},
	}

	reversed := reverseIP(ip)
	if reversed == "" {
		c.resultChan <- models.BlacklistResult{
			IP:        ip,
			Error:     "Invalid IP address format",
			Timestamp: time.Now().Format(time.RFC3339),
		}
		return
	}

	allBlacklists := append(c.config.Blacklists, c.config.CustomBlacklists...)
	for _, bl := range allBlacklists {
		fqdn := fmt.Sprintf("%s.%s", reversed, bl)

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.config.Timeout)*time.Second)
		_, err := resolver.LookupHost(ctx, fqdn)
		cancel()

		if err == nil {
			c.resultChan <- models.BlacklistResult{
				IP:        ip,
				Blacklist: bl,
				Timestamp: time.Now().Format(time.RFC3339),
			}
		}
	}
}

// GetResultChan returns the result channel
func (c *Checker) GetResultChan() <-chan models.BlacklistResult {
	return c.resultChan
}

// Close closes the result channel
func (c *Checker) Close() {
	close(c.resultChan)
}
