// Alptekin Sünnetci (C) 2025
// This program checks if a given IP address or a /24 subnet is blacklisted in various DNS-based blackhole lists (DNSBLs).
// blacklist-checker v.1.1

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"blacklist-check/pkg/config"
	"blacklist-check/pkg/dnsbl"
	"blacklist-check/pkg/utils"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "config.json", "Path to configuration file")
	outputFormat := flag.String("format", "json", "Output format (json or text)")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Override output format if specified
	if *outputFormat != "" {
		cfg.OutputFormat = *outputFormat
	}

	// Check if IP/subnet argument is provided
	if len(flag.Args()) != 1 {
		fmt.Println("Usage: blacklist-checker [options] <IP> or <subnet/24>")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Parse input
	ipsToCheck, subnetOrIP, err := utils.ParseInput(flag.Arg(0))
	if err != nil {
		fmt.Printf("Error parsing input: %v\n", err)
		os.Exit(1)
	}

	// Create checker
	checker := dnsbl.NewChecker(cfg)
	defer checker.Close()

	// Setup concurrency control
	var wg sync.WaitGroup
	sem := make(chan struct{}, cfg.Concurrency)

	// Start result collector
	results := make(map[string][]string)
	go func() {
		for result := range checker.GetResultChan() {
			if result.Error != "" {
				fmt.Printf("Error checking %s: %s\n", result.IP, result.Error)
				continue
			}
			results[result.IP] = append(results[result.IP], result.Blacklist)
		}
	}()

	// Start checking IPs
	fmt.Printf("Checking %d IPs against %d blacklists...\n", len(ipsToCheck), len(cfg.Blacklists))
	startTime := time.Now()

	for _, ip := range ipsToCheck {
		wg.Add(1)
		go checker.CheckIP(ip, &wg, sem)
	}

	wg.Wait()

	// Format and save results
	if len(results) == 0 {
		fmt.Println("No blacklisted IPs found.")
		return
	}

	formattedResults, err := utils.FormatResults(results, cfg.OutputFormat)
	if err != nil {
		fmt.Printf("Error formatting results: %v\n", err)
		os.Exit(1)
	}

	// Save results to file
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s.%s",
		strings.ReplaceAll(subnetOrIP, "/", "-"),
		timestamp,
		cfg.OutputFormat)

	if err := os.WriteFile(filename, []byte(formattedResults), 0644); err != nil {
		fmt.Printf("Error saving results: %v\n", err)
		os.Exit(1)
	}

	duration := time.Since(startTime)
	fmt.Printf("Check completed in %v\n", duration)
	fmt.Printf("Results saved to %s\n", filename)
}
