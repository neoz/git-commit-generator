#Requires -Version 5.0

# Git Commit Generator Installer
# This script installs the Git Commit Generator tool and configures it as a git alias

$ErrorActionPreference = "Stop"

function Write-ColorOutput {
    param(
        [Parameter(Mandatory = $true)]
        [string]$Message,
        
        [Parameter(Mandatory = $false)]
        [string]$ForegroundColor = "White"
    )
    
    Write-Host $Message -ForegroundColor $ForegroundColor
}

function Check-Command {
    param (
        [string]$Command
    )
    
    try {
        $null = Get-Command $Command -ErrorAction Stop
        return $true
    } catch {
        return $false
    }
}

# Check if required tools are installed
Write-ColorOutput "Checking required dependencies..." "Cyan"

$requiredTools = @{
    "git" = "Git is required. Please install it from https://git-scm.com/downloads"
    "go" = "Go is required. Please install it from https://golang.org/dl/"
    "ollama" = "Ollama is required. Please install it from https://ollama.com"
}

$missingTools = @()
foreach ($tool in $requiredTools.Keys) {
    if (-not (Check-Command $tool)) {
        $missingTools += $tool
        Write-ColorOutput "❌ $tool not found: $($requiredTools[$tool])" "Red"
    } else {
        Write-ColorOutput "✅ $tool is installed" "Green"
    }
}

if ($missingTools.Count -gt 0) {
    Write-ColorOutput "Please install the missing dependencies and run this script again." "Yellow"
    exit 1
}

# Install the Git Commit Generator using go install
Write-ColorOutput "Installing Git Commit Generator using go install..." "Cyan"
try {
    # Use go install to get the package
    $env:GO111MODULE = "on"
    & go install github.com/neoz/git-commit-generator@latest
    
    if ($LASTEXITCODE -ne 0) {
        Write-ColorOutput "Failed to install Git Commit Generator." "Red"
        exit 1
    }
} catch {
    Write-ColorOutput "Error installing Git Commit Generator: $_" "Red"
    exit 1
}

# Find the Go binary path
$goPath = & go env GOPATH
$gobin = if ((& go env GOBIN) -ne "") { & go env GOBIN } else { Join-Path -Path $goPath -ChildPath "bin" }
$executablePath = Join-Path -Path $gobin -ChildPath "git-commit-generator.exe"

# Verify the executable exists
if (-not (Test-Path $executablePath)) {
    Write-ColorOutput "❌ Could not find the installed git-commit-generator executable at $executablePath" "Red"
    exit 1
}

# Pull Ollama models
Write-ColorOutput "Pulling required Ollama models..." "Cyan"
ollama pull tavernari/git-commit-message:reasoning
ollama pull tavernari/git-commit-message:merge_commits

# Set up the git alias
Write-ColorOutput "Setting up git alias..." "Cyan"
$gitAliasCommand = "!`"$executablePath`""
git config --global alias.commit-gen $gitAliasCommand

if ($LASTEXITCODE -eq 0) {
    Write-ColorOutput "✅ Git alias 'commit-gen' has been configured successfully!" "Green"
    Write-ColorOutput "You can now use 'git commit-gen' to generate commit messages." "Green"
} else {
    Write-ColorOutput "❌ Failed to set up the git alias." "Red"
    exit 1
}

# Show installation summary
Write-ColorOutput "`nInstallation Summary:" "Cyan"
Write-ColorOutput "---------------------" "Cyan"
Write-ColorOutput "✅ Git Commit Generator installed at: $executablePath" "Green"
Write-ColorOutput "✅ Git alias 'commit-gen' configured" "Green"
Write-ColorOutput "✅ Required Ollama models pulled" "Green"
Write-ColorOutput "`nUsage:" "Yellow"
Write-ColorOutput "------" "Yellow"
Write-ColorOutput "1. Stage your changes: git add <files>" "White"
Write-ColorOutput "2. Run: git commit-gen" "White"
Write-ColorOutput "3. Follow the prompts to generate and use your commit message" "White"
