# Gira

A CLI tool to manage your Jira tasks from the terminal.

## Installation

### Download Binary

Download the latest binary for your platform from [Releases](https://github.com/enBonnet/gira/releases):

```bash
# Linux/macOS - make it executable and move to PATH
chmod +x gira
sudo mv gira /usr/local/bin/

# Windows - add to PATH
# Move gira.exe to a folder like C:\Program Files\gira\
# Add that folder to your PATH environment variable
```

### Go Install

```bash
go install github.com/enBonnet/gira@latest
```

The binary will be installed to `$GOPATH/bin` (or `$HOME/go/bin` by default).

### Build from Source

```bash
git clone https://github.com/enBonnet/gira.git
cd gira
go build -o gira .
```

**Add to PATH:**

```bash
# Linux/macOS
sudo mv gira /usr/local/bin/

# Or add to your local bin
mkdir -p ~/bin
mv gira ~/bin/
# Add to your shell config (~/.bashrc, ~/.zshrc, etc.)
export PATH="$HOME/bin:$PATH"
```

## Configuration

Create a config file at `~/.config/gira/gira.yaml` or `~/.gira.yaml`:

```yaml
jira_url: https://your-company.atlassian.net
email: your-email@example.com
api_token: your-jira-api-token
project: PROJECT_KEY  # optional
```

Alternatively, use environment variables:

```bash
export GIRA_JIRA_URL=https://your-company.atlassian.net
export GIRA_EMAIL=your-email@example.com
export GIRA_API_TOKEN=your-jira-api-token
export GIRA_PROJECT=PROJECT_KEY  # optional
```

### Getting your API Token

1. Go to [Atlassian Account Settings](https://id.atlassian.com/manage-profile/security/api-tokens)
2. Click "Create API token"
3. Label it (e.g., "gira-cli")
4. Copy the token and use it in your config

## Usage

### List your tasks

```bash
gira tasks
```

Output:
```
KEY        STATUS   PRIORITY  SUMMARY                                                        URL
LPTS-2028  To Do    Highest   Rewards & Referrals V2 FE                                       https://your-company.atlassian.net/browse/LPTS-2028
LPTS-2345  To Do    Highest   Property Admin's should not be able to add users               https://your-company.atlassian.net/browse/LPTS-2345
```

#### Options

- `-l, --limit int` - Maximum number of tasks to display (default: 5)
- Use `-l 0` to show all tasks

Examples:
```bash
gira tasks              # Show 5 tasks (default)
gira tasks -l 10        # Show 10 tasks
gira tasks -l 0         # Show all tasks
```

Tasks are filtered by status: **To Do**, **BLOCKED**, **In Progress**

Tasks are sorted by priority (Highest → Lowest).

### Show current user info

```bash
gira myself
```

Output:
```
Account ID: 5f1234567890abcdef123456
Name:       John Doe
Email:      john.doe@example.com
Active:     true
Timezone:   America/New_York
Locale:     en_US
```

## Development

### Requirements

- Go 1.21+

### Build

```bash
go build -o gira .
```

### Project Structure

```
gira/
├── main.go              # Entry point
├── cmd/
│   ├── root.go          # Root command
│   ├── tasks.go         # tasks subcommand
│   └── myself.go        # myself subcommand
└── internal/
    ├── config/
    │   └── config.go    # Configuration loading
    └── jira/
        └── client.go    # Jira API client
```
