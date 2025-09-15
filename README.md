# 🌙 Luna - AI Commit Generator

Luna is an intelligent Git commit message generator powered by Google's Gemini AI. It automatically analyzes your staged changes and creates meaningful, emoji-enhanced commit messages that follow conventional commit standards.

## ✨ Features

- 🤖 **AI-Powered**: Uses Google Gemini 2.0 Flash for intelligent commit message generation
- 🎨 **Emoji Enhancement**: Automatically adds relevant emojis to make commits more expressive
- 📝 **Conventional Commits**: Follows standard commit message conventions (feat:, fix:, chore:, etc.)
- 🚀 **One-Click Generation**: Generates and commits messages for all staged files at once
- 🔒 **Smart Filtering**: Automatically skips binary files (exe, dll, images, etc.)
- 🎯 **Concise Messages**: Keeps commit messages under 60 characters for better readability
- 🌈 **Colorful Output**: Beautiful colored terminal output for better user experience

## 🚀 Quick Start

### Prerequisites

- Windows operating system
- Git installed and configured
- Google Gemini API key ([Get one here](https://aistudio.google.com/app/apikey))

### Installation

1. **Download the executable**
   ```bash
   # Clone this repository
   git clone https://github.com/6hax/Luna.git luna
   cd luna
   ```

2. **Add Luna to Windows PATH**
   
   **Option A: Using System Properties (Recommended)**
   - Press `Win + R`, type `sysdm.cpl` and press Enter
   - Click "Environment Variables"
   - Under "System Variables", find and select "Path", then click "Edit"
   - Click "New" and add the full path to your Luna folder (e.g., `C:\Users\YourUsername\luna`)
   - Click "OK" on all dialogs
   
   **Option B: Using Command Prompt (Temporary)**
   ```cmd
   setx PATH "%PATH%;C:\path\to\your\luna\folder"
   ```
   
   **Option C: Using PowerShell (Temporary)**
   ```powershell
   $env:PATH += ";C:\path\to\your\luna\folder"
   ```

3. **Set up your API key**
   ```bash
   LunaApikey YOUR_GEMINI_API_KEY
   ```
   > **Note**: Close and reopen your terminal after setting the API key

4. **Verify installation**
   ```bash
   LunaHelp
   ```

## 📖 Usage

### Basic Commands

| Command | Description |
|---------|-------------|
| `LunaHelp` | Show help information and available commands |
| `LunaCommit` | Generate and commit messages for all staged files |
| `LunaApikey YOUR_KEY` | Set your Gemini API key |

### Step-by-Step Workflow

1. **Stage your changes**
   ```bash
   git add .
   # or stage specific files
   git add file1.js file2.py
   ```

2. **Generate commit messages**
   ```bash
   LunaCommit
   ```

3. **That's it!** Luna will:
   - Analyze each staged file's changes
   - Generate appropriate commit messages using AI
   - Add random emojis and conventional commit prefixes
   - Create individual commits for each file

### Example Output

```
Generating commit for file: src/main.go
Committed src/main.go with message:
🚀 feat: add user authentication system

Generating commit for file: README.md
Committed README.md with message:
📝 docs: update installation instructions
```

## 🎯 How It Works

1. **File Analysis**: Luna scans your staged Git changes using `git diff --cached`
2. **AI Processing**: Sends the diff to Google Gemini 2.0 Flash API for analysis
3. **Message Generation**: AI generates concise, meaningful commit messages
4. **Enhancement**: Adds random emojis and conventional commit prefixes
5. **Commit Creation**: Automatically commits each file with its generated message

## 🔧 Configuration

### Environment Variables

- `GEMINI_API_KEY`: Your Google Gemini API key (set via `LunaApikey` command)

### Supported File Types

Luna automatically processes most text-based files and skips:
- Binary files (`.exe`, `.dll`)
- Image files (`.png`, `.jpg`, `.jpeg`, `.gif`)
- Other binary formats

## 🎨 Commit Message Format

Luna generates commit messages in this format:
```
[EMOJI] [PREFIX]: [AI-GENERATED MESSAGE]
```

**Examples:**
- `✨ feat: add dark mode toggle`
- `🐛 fix: resolve memory leak in data processing`
- `📝 docs: update API documentation`
- `🔧 refactor: optimize database queries`

### Available Prefixes
- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation changes
- `refactor:` - Code refactoring
- `test:` - Test-related changes
- `chore:` - Maintenance tasks


## 🐛 Troubleshooting

### Common Issues

**"Error: Set GEMINI_API_KEY using LunaApikey first"**
- Solution: Run `LunaApikey YOUR_API_KEY` and restart your terminal

**"Error running git status"**
- Solution: Make sure you're in a Git repository and Git is installed

**"No staged changes to commit"**
- Solution: Stage your files first with `git add .` or `git add filename`

**API Key not working**
- Solution: Verify your API key is correct and has Gemini API access enabled

## 🙏 Acknowledgments

- Google Gemini AI for providing the intelligent commit message generation
- The Go community for excellent tooling and libraries

## ⭐ Support the Project

If Luna has helped you write better commit messages, please consider:

- ⭐ **Star this repository** - It helps others discover the tool
- 🐛 **Report bugs** - Help us improve by reporting issues
- 💡 **Suggest features** - Tell us what you'd like to see next
- 👥 **Follow the author** - Stay updated with new projects and updates

---

**Made with ❤️ for hax**

*Luna - Making Git commits more intelligent, one message at a time.*
