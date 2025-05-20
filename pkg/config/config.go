package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"blacklist-check/pkg/models"
)

var defaultBlacklists = []string{
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

// LoadConfig loads the configuration from a file or creates a default one
func LoadConfig(configPath string) (*models.Config, error) {
	if configPath == "" {
		configPath = "config.json"
	}

	config := &models.Config{
		Concurrency:  100,
		Timeout:      3,
		Blacklists:   defaultBlacklists,
		OutputFormat: "json",
	}

	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("error reading config file: %v", err)
		}

		if err := json.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("error parsing config file: %v", err)
		}
	} else {
		// Create default config file
		data, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("error creating default config: %v", err)
		}

		if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
			return nil, fmt.Errorf("error creating config directory: %v", err)
		}

		if err := os.WriteFile(configPath, data, 0644); err != nil {
			return nil, fmt.Errorf("error writing default config: %v", err)
		}
	}

	return config, nil
}

// SaveConfig saves the configuration to a file
func SaveConfig(config *models.Config, configPath string) error {
	if configPath == "" {
		configPath = "config.json"
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("error writing config file: %v", err)
	}

	return nil
}
