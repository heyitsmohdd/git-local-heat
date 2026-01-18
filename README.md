# git-local-heat

Visualize your local git commit history as a GitHub-style heatmap in your terminal.

## Features

- Recursive Scanning - Automatically finds all `.git` repositories in any directory tree
- High Performance - Concurrent processing with Go goroutines for blazing-fast scan times
- Beautiful Output - GitHub-style contribution graph with 5 shades of green
- Email Filtering - Filter commits by author email (auto-detects from git config)
- Flexible Timeframes - View 6 or 12 months of commit history
- Configurable - Supports YAML config files and CLI flags
- Robust - Graceful error handling for permission issues and corrupted repos

## Preview

```
📊 Git Contribution Heatmap

Total commits: 123  |  Max in a day: 23

    Oct       Nov         Dec       Jan   
    ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ 
Mon ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ 
    ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ 
Wed ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ 
    ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ 
Fri ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ 
    ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ 

Less ■ ■ ■ ■ ■  More
```

## Installation

### Via Go Install (Recommended)

```bash
go install github.com/heyitsmohdd/git-local-heat/cmd/git-local-heat@latest
```

### From Source

```bash
# Clone the repository
git clone https://github.com/shariqattar/git-local-heat.git
cd git-local-heat

# Build the binary
go build -o git-local-heat cmd/git-local-heat/main.go

# (Optional) Move to PATH
sudo mv git-local-heat /usr/local/bin/
```

### Pre-built Binaries

Download pre-built binaries from the [Releases](https://github.com/shariqattar/git-local-heat/releases) page.

## Usage

### Basic Usage

```bash
# Scan current directory
git-local-heat

# Scan a specific directory
git-local-heat --path ~/Developer

# View last 6 months
git-local-heat --months 6

# Filter by email
git-local-heat --email user@example.com
```

### All Options

```bash
git-local-heat [flags]

Flags:
  -c, --config string   Config file path (default: ~/.git-local-heat.yaml)
  -e, --email string    Filter commits by author email (default: git config user.email)
  -h, --help            help for git-local-heat
  -m, --months int      Number of months to display (6 or 12) (default 12)
  -p, --path string     Root directory to scan (default: current directory)
```

### Configuration File

Create a config file at `~/.git-local-heat.yaml`:

```yaml
email: your@email.com
path: /Users/you/Developer
months: 12
colorScheme:
  - "#161b22"  # No contributions
  - "#0e4429"  # Low (1-2 commits)
  - "#006d32"  # Medium-low (3-5 commits)
  - "#26a641"  # Medium-high (6-10 commits)
  - "#39d353"  # High (11+ commits)
```

**Note:** CLI flags take precedence over config file values.

## Project Structure

```
git-local-heat/
├── cmd/
│   └── git-local-heat/
│       └── main.go              # CLI entry point
├── internal/
│   ├── scanner/
│   │   └── scanner.go           # Concurrent repo discovery
│   ├── parser/
│   │   └── git_parser.go        # Commit history extraction
│   ├── heatmap/
│   │   ├── generator.go         # Grid calculation
│   │   └── renderer.go          # Terminal rendering
│   └── config/
│       └── config.go            # Config management
├── pkg/
│   └── models/
│       └── models.go            # Shared data structures
└── README.md
```

## Technology Stack

- **Go 1.21+** - Core language
- **[cobra](https://github.com/spf13/cobra)** - CLI framework
- **[go-git](https://github.com/go-git/go-git)** - Pure Go git implementation
- **[lipgloss](https://github.com/charmbracelet/lipgloss)** - Terminal styling
- **[yaml.v3](https://gopkg.in/yaml.v3)** - YAML parsing

## Performance

The tool uses Go's concurrency primitives for maximum performance:

- **Worker pool pattern** for directory scanning (50 concurrent workers)
- **Goroutines** for parallel commit parsing across repositories
- **Efficient date aggregation** with in-memory hashmaps

Typical performance on a development machine:
- **100+ repositories** scanned in **< 5 seconds**
- **Thousands of commits** processed in parallel

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Setup

```bash
# Clone the repo
git clone https://github.com/heyitsmohdd/git-local-heat.git
cd git-local-heat

# Install dependencies
go mod download

# Run tests
go test ./...

# Build
go build -o git-local-heat cmd/git-local-heat/main.go
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by GitHub's contribution graph
- Built with amazing Go libraries from the community

## Known Issues & Limitations

- Only displays commits from the default branch (HEAD)
- Does not support bare repositories
- Large monorepos (100k+ commits) may take longer to parse

## Contact

 
GitHub: (https://github.com/heyitsmohdd)

---

If you find this tool useful, please consider starring the repository.
