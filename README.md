
# Blacklist Checker

## Overview

This Go application checks whether a given IP address or a /24 subnet is blacklisted across multiple DNS-based blacklists. The application takes either a single IP address or a /24 subnet as input and outputs the blacklisted IPs into a JSON file.

## Features

- Supports both single IP addresses and /24 subnets as input.
- Checks IPs against a wide range of DNS-based blacklists.
- Outputs the results in a JSON format, grouped by IP address with their corresponding blacklists.
- Handles concurrent DNS lookups with a limit on the number of concurrent requests.
- Supports automatic generation of output filenames with timestamps.

## Prerequisites

- Go version 1.16 or higher.
- Access to the internet to perform DNS lookups.

## Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/alptekisunnetci/blacklist-checker.git
   cd blacklist-checker
   ```

2. **Build the Go application:**

   ```bash
   go build -o blacklist_checker main.go
   ```

3. **Run the application:**

   ```bash
   ./blacklist_checker <IP> or <subnet/24>
   ```

   Replace `<IP>` with the IP address or `<subnet/24>` with the subnet you want to check. For example:
   - `./blacklist_checker 218.92.0.211`
   - `./blacklist_checker 218.92.0.0/24`

## Usage

The application supports the following input formats:

- **Single IP address:** Provide an IP address as an argument. The format should be in the standard IPv4 format (e.g., `218.92.0.211`).
- **/24 Subnet:** Provide a subnet in CIDR format (e.g., `218.92.0.0/24`). Only `/24` subnets are supported.

### Example Commands

- Checking a single IP address:
  
  ```bash
  ./blacklist_checker 218.92.0.211
  ```

- Checking a subnet:

  ```bash
  ./blacklist_checker 218.92.0.0/24
  ```

## Output

- The application generates a JSON file with the results of the blacklist check.
- The output file is named in the format `subnet_or_ip_timestamp.json`, where `subnet_or_ip` is the provided input (with `/` replaced by `-`), and `timestamp` is the current timestamp.

Example filename: `218.92.0.0-24_20230514_153000.json`

### JSON Output Structure

The JSON output will contain a list of IPs as keys, with their respective blacklists as values.

Example:

```json
{
  "218.92.0.1": [
    "dnsbl-2.uceprotect.net"
  ],
  "218.92.0.220": [
    "dnsbl-1.uceprotect.net",
    "dnsbl-2.uceprotect.net",
    "dnsbl.dronebl.org"
  ],
  "218.92.0.221": [
    "dnsbl-1.uceprotect.net",
    "dnsbl-2.uceprotect.net",
    "dnsbl.dronebl.org"
  ],
  "218.92.0.222": [
    "dnsbl-1.uceprotect.net",
    "dnsbl-2.uceprotect.net"
  ],
  "218.92.0.223": [
    "dnsbl-1.uceprotect.net",
    "dnsbl-2.uceprotect.net",
    "dnsbl.dronebl.org"
  ],
  "218.92.0.224": [
    "dnsbl-2.uceprotect.net"
  ],
  "218.92.0.225": [
    "dnsbl-1.uceprotect.net",
    "dnsbl-2.uceprotect.net",
    "dnsbl.dronebl.org"
  ],
  "218.92.0.226": [
    "dnsbl-1.uceprotect.net",
    "dnsbl-2.uceprotect.net",
    "dnsbl.dronebl.org"
  ],
  "218.92.0.227": [
    "dnsbl-1.uceprotect.net",
    "dnsbl-2.uceprotect.net"
  ],
  "218.92.0.228": [
    "dnsbl-1.uceprotect.net",
    "dnsbl-2.uceprotect.net",
    "dnsbl.dronebl.org"
  ],
  "218.92.0.229": [
    "dnsbl-1.uceprotect.net",
    "dnsbl-2.uceprotect.net",
    "dnsbl.dronebl.org"
  ],
  "218.92.0.230": [
    "dnsbl-1.uceprotect.net",
    "dnsbl-2.uceprotect.net",
    "dnsbl.dronebl.org"
  ],
  "218.92.0.231": [
    "dnsbl-1.uceprotect.net",
    "dnsbl-2.uceprotect.net",
    "dnsbl.dronebl.org"
  ],
  "218.92.0.232": [
    "dnsbl-1.uceprotect.net",
    "dnsbl-2.uceprotect.net",
    "dnsbl.dronebl.org"
  ],
  "218.92.0.233": [
    "dnsbl-1.uceprotect.net",
    "dnsbl-2.uceprotect.net"
  ],
  "218.92.0.234": [
    "dnsbl-2.uceprotect.net"
  ],
  "218.92.0.235": [
    "dnsbl-1.uceprotect.net",
    "dnsbl-2.uceprotect.net",
    "dnsbl.dronebl.org"
  ],
  "218.92.0.236": [
    "dnsbl-1.uceprotect.net",
    "dnsbl-2.uceprotect.net",
    "dnsbl.dronebl.org"
  ],
  "218.92.0.237": [
    "dnsbl-1.uceprotect.net",
    "dnsbl-2.uceprotect.net",
    "dnsbl.dronebl.org"
  ],
  "218.92.0.252": [
    "dnsbl-1.uceprotect.net",
    "dnsbl-2.uceprotect.net"
  ],
  "218.92.0.52": [
    "dnsbl-2.uceprotect.net",
    "dnsbl.dronebl.org"
  ],
  "218.92.0.53": [
    "dnsbl-2.uceprotect.net",
    "dnsbl.dronebl.org",
    "rbl.rtbh.com.tr"
  ],
  "218.92.0.56": [
    "b.barracudacentral.org",
    "dnsbl-2.uceprotect.net",
    "rbl.rtbh.com.tr"
  ],
  "218.92.0.57": [
    "dnsbl-2.uceprotect.net"
  ],
  "218.92.0.99": [
    "dnsbl-2.uceprotect.net",
    "rbl.rtbh.com.tr"
  ]
}
```

## Concurrency and Performance

- The application uses goroutines to perform concurrent DNS lookups, with a limit on the number of concurrent requests (`50` by default).
- You can adjust the concurrency by modifying the `concurrency` variable in the `main` function.

## License

1. **Personal Use**: You are free to use this project for personal, non-commercial purposes. This includes viewing, downloading, and utilizing the software for personal learning, experimentation, or non-commercial projects.
   
2. **Non-Commercial Use**: You may not use this project for commercial purposes. This means you cannot sell, offer for sale, or use the project in any manner that is intended for commercial gain or financial benefit.

3. **No Commercial Use of Derivatives**: Even if you receive permission to modify the project, you are prohibited from using the modified project for commercial purposes unless explicitly authorized by the project’s creators.

4. **Attribution**: If you use the project, you must provide appropriate attribution to the original authors. You may not imply that the authors endorse your use of the project.
