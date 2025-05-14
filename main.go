// Alptekin Sünnetci (C) 2025
// This program checks if a given IP address or a /24 subnet is blacklisted in various DNS-based blackhole lists (DNSBLs).

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type BlacklistResult struct {
	IP        string
	Blacklist string
}

var blacklists = []string{
	"access.redhawk.org",
	"b.barracudacentral.org",
	"bl.spamcop.net",
	"blackholes.mail-abuse.org",
	"bogons.cymru.com",
	"cdl.anti-spam.org.cn",
	"db.wpbl.info",
	"dnsbl-1.uceprotect.net",
	"dnsbl-2.uceprotect.net",
	"dnsbl.dronebl.org",
	"dnsbl.sorbs.net",
	"drone.abuse.ch",
	"dul.dnsbl.sorbs.net",
	"http.dnsbl.sorbs.net",
	"httpbl.abuse.ch",
	"ips.backscatterer.org",
	"ix.dnsbl.manitu.net",
	"multi.surbl.org",
	"netblock.pedantic.org",
	"psbl.surriel.com",
	"query.senderbase.org",
	"rbl-plus.mail-abuse.org",
	"rbl.efnetrbl.org",
	"rbl.spamlab.com",
	"relays.mail-abuse.org",
	"short.rbl.jp",
	"smtp.dnsbl.sorbs.net",
	"socks.dnsbl.sorbs.net",
	"spam.dnsbl.sorbs.net",
	"spamguard.leadmon.net",
	"spamrbl.imp.ch",
	"ubl.unsubscore.com",
	"web.dnsbl.sorbs.net",
	"wormrbl.imp.ch",
	"zombie.dnsbl.sorbs.net",
	"rbl.rtbh.com.tr",
}

func parseInput(arg string) ([]string, string, error) {
	if !strings.Contains(arg, "/") {
		ip := net.ParseIP(arg)
		if ip == nil || ip.To4() == nil {
			return nil, "", fmt.Errorf("Invalid IP address: %s", arg)
		}
		return []string{ip.String()}, ip.String(), nil
	}

	ip, ipNet, err := net.ParseCIDR(arg)
	if err != nil {
		return nil, "", fmt.Errorf("Invalid subnet format: %s", arg)
	}

	ones, bits := ipNet.Mask.Size()
	if bits != 32 || ones != 24 {
		return nil, "", fmt.Errorf("Only /24 subnet is supported")
	}

	var ips []string
	base := ip.To4()
	for i := 0; i < 256; i++ {
		ips = append(ips, fmt.Sprintf("%d.%d.%d.%d", base[0], base[1], base[2], i))
	}
	return ips, ipNet.String(), nil
}

func reverseIP(ip string) string {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return ""
	}
	return fmt.Sprintf("%s.%s.%s.%s", parts[3], parts[2], parts[1], parts[0])
}

func checkBlacklist(ip string, wg *sync.WaitGroup, sem chan struct{}, resultChan chan BlacklistResult) {
	defer wg.Done()
	sem <- struct{}{}

	resolver := net.Resolver{PreferGo: true}
	reversed := reverseIP(ip)

	for _, bl := range blacklists {
		fqdn := fmt.Sprintf("%s.%s", reversed, bl)
		fmt.Printf("Checking %s against %s...\n", ip, bl)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		_, err := resolver.LookupHost(ctx, fqdn)
		cancel()

		if err == nil {
			resultChan <- BlacklistResult{IP: ip, Blacklist: bl}
		}
	}

	<-sem
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./blacklist_checker <IP> or <subnet/24>")
		os.Exit(1)
	}

	ipsToCheck, subnetOrIP, err := parseInput(os.Args[1])
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	concurrency := 50
	sem := make(chan struct{}, concurrency)
	resultChan := make(chan BlacklistResult, 256)

	grouped := make(map[string][]string)

	go func() {
		for res := range resultChan {
			grouped[res.IP] = append(grouped[res.IP], res.Blacklist)
		}
	}()

	for _, ip := range ipsToCheck {
		wg.Add(1)
		go checkBlacklist(ip, &wg, sem, resultChan)
	}

	wg.Wait()
	close(resultChan)

	if len(grouped) == 0 {
		fmt.Println("No blacklisted IPs found.")
		return
	}

	jsonData, err := json.MarshalIndent(grouped, "", "  ")
	if err != nil {
		fmt.Println("JSON error:", err)
		return
	}

	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s.json", strings.ReplaceAll(subnetOrIP, "/", "-"), timestamp)

	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		fmt.Println("File write error:", err)
		return
	}

	fmt.Printf("Blacklisted IPs written to %s\n", filename)
}
