
# Blacklist Checker

A high-performance DNS-based blackhole list (DNSBL) checker written in Go. This tool checks if IP addresses or /24 subnets are listed in various DNSBLs.

## Features

- Check individual IP addresses or /24 subnets
- Configurable concurrency and timeout settings
- Support for custom DNSBL lists
- Multiple output formats (JSON and text)
- Detailed error reporting
- Performance optimized DNS lookups
- Configurable through command-line flags and config file

## Installation

```bash
git clone https://github.com/alptekinsunnetci/blacklist-checker.git
cd blacklist-checker
go build
```

## Usage

Basic usage:
```bash
./blacklist-checker <IP> or <subnet/24>
```

With options:
```bash
./blacklist-checker -config custom_config.json -format text <IP> or <subnet/24>
```

### Command-line Options

- `-config`: Path to configuration file (default: config.json)
- `-format`: Output format (json or text, default: json)

### Configuration File

The configuration file (config.json) supports the following options:

```json
{
  "concurrency": 100,
  "timeout": 3,
  "blacklists": [
    "access.redhawk.org",
    "b.barracudacentral.org",
    ...
  ],
  "output_format": "json",
  "custom_blacklists": []
}
```

## Output

The tool generates a timestamped output file in the specified format. The filename follows the pattern:
`<IP-or-subnet>_<timestamp>.<format>`

### JSON Output Example
```json
{
  "192.168.1.1": [
    "dnsbl.sorbs.net",
    "bl.spamcop.net"
  ]
}
```

### Text Output Example
```
IP: 192.168.1.1
Blacklisted in:
  - dnsbl.sorbs.net
  - bl.spamcop.net
```

## Performance

The tool is optimized for performance with:
- Configurable concurrency limits
- Efficient DNS lookups
- Goroutine-based parallel processing
- Memory-efficient result collection

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

1. **Personal Use**: You are free to use this project for personal, non-commercial purposes. This includes viewing, downloading, and utilizing the software for personal learning, experimentation, or non-commercial projects.
   
2. **Non-Commercial Use**: You may not use this project for commercial purposes. This means you cannot sell, offer for sale, or use the project in any manner that is intended for commercial gain or financial benefit.

3. **No Commercial Use of Derivatives**: Even if you receive permission to modify the project, you are prohibited from using the modified project for commercial purposes unless explicitly authorized by the project’s creators.

4. **Attribution**: If you use the project, you must provide appropriate attribution to the original authors. You may not imply that the authors endorse your use of the project.
