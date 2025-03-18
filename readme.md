# Git Commit Generator

## Prerequisites

Before setting up the Git Commit Generator, ensure you have the following installed:

1. **Ollama**:  
   Ollama is required for this tool. To install it, follow the instructions on the [official Ollama website](https://ollama.com).

## Installation Guide

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
