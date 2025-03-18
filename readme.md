# Git Commit Generator

## Prerequisites

Before setting up the Git Commit Generator, ensure you have the following installed:

1. **Ollama**:  
   Ollama is required for this tool. To install it, follow the instructions on the [official Ollama website](https://ollama.com).
2. **Go Programming Language**:  
   The Go compiler is required to build the tool. You can download it from the [official Go website](https://golang.org/dl/).
3. **Git**:  
   Git is required for version control integration. You can download it from the [official Git website](https://git-scm.com/downloads).

## Installation Guide

### Automatic Installation

#### Windows
1. Run the provided PowerShell installation script:
   ```powershell
   .\install.ps1
   ```

#### Linux/macOS
1. Run the provided shell script:
   ```bash
   ./install.sh
   ```

The installation scripts will:
- Check and verify required dependencies
- Build and install the Go executable
- Pull the specified Ollama models
- Configure the git alias automatically

### Manual Installation

To set up a Git alias for running the `git-commit-generator`, follow these steps:

1. Open your terminal or command prompt.
2. Run the following command to configure the alias:

   ```bash
   git config --global alias.commit-gen '!<path-to-git-commit-generator>/git-commit-generator'
   ```

   Replace `<path-to-git-commit-generator>` with the absolute path to the `git-commit-generator` executable. For example:
   - On Windows: `C:/working/go/git-commit-generator/git-commit-generator.exe`
   - On Linux/macOS: `/path/to/git-commit-generator/git-commit-generator`

3. Verify the alias by running:

   ```bash
   git commit-gen
   ```

   This should execute the `git-commit-generator` file.

## Usage

Once the alias is set up, you can use `git commit-gen` to run the Git Commit Generator tool.

1. Stage your changes using `git add <files>`
2. Run `git commit-gen`
3. Follow the prompts to generate and use your commit message

## Copyright

This project uses resources from the following:

- [Git Commit Generator Script by Tavernari](https://gist.githubusercontent.com/Tavernari/b88680e71c281cfcdd38f46bdb164fee/raw/git-gen-commit)
- [Ollama Git Commit Message Model by Tavernari](https://ollama.com/tavernari/git-commit-message)
