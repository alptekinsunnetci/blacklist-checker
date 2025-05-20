package utils

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

// ParseInput parses an IP address or subnet input
func ParseInput(arg string) ([]string, string, error) {
	if !strings.Contains(arg, "/") {
		ip := net.ParseIP(arg)
		if ip == nil || ip.To4() == nil {
			return nil, "", fmt.Errorf("invalid IP address: %s", arg)
		}
		return []string{ip.String()}, ip.String(), nil
	}

	ip, ipNet, err := net.ParseCIDR(arg)
	if err != nil {
		return nil, "", fmt.Errorf("invalid subnet format: %s", arg)
	}

	ones, bits := ipNet.Mask.Size()
	if bits != 32 || ones != 24 {
		return nil, "", fmt.Errorf("only /24 subnet is supported")
	}

	var ips []string
	base := ip.To4()
	for i := 0; i < 256; i++ {
		ips = append(ips, fmt.Sprintf("%d.%d.%d.%d", base[0], base[1], base[2], i))
	}
	return ips, ipNet.String(), nil
}

// FormatResults formats the results for output
func FormatResults(results map[string][]string, format string) (string, error) {
	switch format {
	case "json":
		return formatJSON(results)
	case "text":
		return formatText(results)
	default:
		return "", fmt.Errorf("unsupported output format: %s", format)
	}
}

func formatJSON(results map[string][]string) (string, error) {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error formatting JSON: %v", err)
	}
	return string(data), nil
}

func formatText(results map[string][]string) (string, error) {
	var sb strings.Builder
	for ip, blacklists := range results {
		sb.WriteString(fmt.Sprintf("IP: %s\n", ip))
		sb.WriteString("Blacklisted in:\n")
		for _, bl := range blacklists {
			sb.WriteString(fmt.Sprintf("  - %s\n", bl))
		}
		sb.WriteString("\n")
	}
	return sb.String(), nil
}
