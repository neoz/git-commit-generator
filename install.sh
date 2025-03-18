#!/bin/bash

set -e

print_colored() {
    local message=$1
    local color=$2
    case $color in
        "red") code=31 ;;
        "green") code=32 ;;
        "yellow") code=33 ;;
        "cyan") code=36 ;;
        *) code=37 ;;
    esac
    echo -e "\033[${code}m${message}\033[0m"
}

check_command() {
    if ! command -v "$1" &> /dev/null; then
        return 1
    fi
    return 0
}

# Check required dependencies
print_colored "Checking required dependencies..." "cyan"

declare -A required_tools=(
    ["git"]="Git is required. Please install it from your package manager or https://git-scm.com/downloads"
    ["go"]="Go is required. Please install it from your package manager or https://golang.org/dl/"
    ["ollama"]="Ollama is required. Please install it from https://ollama.com"
)

missing_tools=0
for tool in "${!required_tools[@]}"; do
    if ! check_command "$tool"; then
        print_colored "❌ $tool not found: ${required_tools[$tool]}" "red"
        missing_tools=1
    else
        print_colored "✅ $tool is installed" "green"
    fi
done

if [ $missing_tools -eq 1 ]; then
    print_colored "Please install the missing dependencies and run this script again." "yellow"
    exit 1
fi

# Install Git Commit Generator
print_colored "Installing Git Commit Generator using go install..." "cyan"
export GO111MODULE=on
if ! go install github.com/neoz/git-commit-generator@latest; then
    print_colored "Failed to install Git Commit Generator." "red"
    exit 1
fi

# Find the Go binary path
GOPATH=$(go env GOPATH)
GOBIN=$(go env GOBIN)
if [ -z "$GOBIN" ]; then
    GOBIN="$GOPATH/bin"
fi
EXECUTABLE_PATH="$GOBIN/git-commit-generator"

# Verify the executable exists
if [ ! -f "$EXECUTABLE_PATH" ]; then
    print_colored "❌ Could not find the installed git-commit-generator executable at $EXECUTABLE_PATH" "red"
    exit 1
fi

# Pull Ollama models
print_colored "Pulling required Ollama models..." "cyan"
ollama pull tavernari/git-commit-message:reasoning
ollama pull tavernari/git-commit-message:merge_commits

# Set up git alias
print_colored "Setting up git alias..." "cyan"
git config --global alias.commit-gen "!$EXECUTABLE_PATH"

if [ $? -eq 0 ]; then
    print_colored "✅ Git alias 'commit-gen' has been configured successfully!" "green"
    print_colored "You can now use 'git commit-gen' to generate commit messages." "green"
else
    print_colored "❌ Failed to set up the git alias." "red"
    exit 1
fi

# Show installation summary
print_colored "\nInstallation Summary:" "cyan"
print_colored "---------------------" "cyan"
print_colored "✅ Git Commit Generator installed at: $EXECUTABLE_PATH" "green"
print_colored "✅ Git alias 'commit-gen' configured" "green"
print_colored "✅ Required Ollama models pulled" "green"
print_colored "\nUsage:" "yellow"
print_colored "------" "yellow"
print_colored "1. Stage your changes: git add <files>" "white"
print_colored "2. Run: git commit-gen" "white"
print_colored "3. Follow the prompts to generate and use your commit message" "white"
