# 🌙 Luna — AI commit generator 

Luna generates one‑line commit messages for your staged files using the Google Gemini 2.0 Flash API. It reads the diff of each file, asks the model for a concise summary, and creates individual commits with Conventional Commit prefixes. You can optionally include emojis.

## ✨ Features

- 🤖 **Gemini 2.0 Flash**: concise, AI‑powered summaries per file
- 🧩 **Per‑file commits**: one commit for each staged file
- 🏷️ **Conventional Commits**: auto‑adds a prefix if missing
- 🎯 **Concise messages**: prompt asks for < 60 chars; app hard‑cuts to 100
- 🖼️ **Smart filtering**: skips binaries/images (`.exe`, `.dll`, `.png`, `.jpg`, `.jpeg`, `.gif`)
- 😄 **Optional emojis**: enable with `-e`

## 🧩 How it works

1. Collects staged changes with `git status --porcelain` and `git diff --cached -- <file>`
2. Sends each file diff to the Gemini API (`gemini-2.0-flash:generateContent`)
3. If the response does not start with a known prefix, randomly chooses one of `chore:`, `refactor:`, `feat:`, `fix:`, `docs:`, `test:`
4. If you pass `-e`, prepends a random emoji
5. Truncates the final message to 100 characters
6. Runs `git commit -m <message> -- <file>` per file

## 🖥️ Requirements

- Windows
- Git installed and configured
- Google Gemini 2.0 Flash API key (`https://aistudio.google.com/app/apikey`)

## 🔧 Installation

### Option A — Use the prebuilt binary (`Luna.exe`)

1. Place `Luna.exe` in a folder under your user (e.g., `C:\Users\YourUser\Luna`)
2. Add that folder to your Windows `PATH`:
   - `Win + R` → `sysdm.cpl` → “Environment Variables”
   - In “System variables”, edit `Path` → “New” → paste the folder path
   - Save all dialogs

### Option B — Build from source (Go)

```bash
go build -o Luna.exe main.go
```

## 🔐 Set the API key

Set the `GEMINI_API_KEY` environment variable using Luna (persistent on Windows):

```bash
LunaApikey YOUR_GEMINI_KEY
```

After setting it, close and reopen your terminal so the variable takes effect.

## 📦 Usage

You must be inside a Git repository with staged changes.

### Commands and aliases

- `LunaHelp` | `lh` | `-lh`: show help
- `LunaCommit` | `lc` | `-c`: generate and create per‑file commits
- `LunaApikey <YOUR_KEY>` | `lkey <YOUR_KEY>` | `-lkey <YOUR_KEY>`: set the API key

Note: you can call them as executable arguments (e.g., `Luna -c`, `Luna lc`, `Luna -lh`) or directly (`LunaCommit`, etc.).

### Typical flow

```bash
git add .
Luna -c        # or: Luna lc, or: LunaCommit
```

### Optional emojis

```bash
Luna -c -e     # enable emojis in messages
```

## 🧪 Example output

```
Generating commit for file: src/main.go
Committed src/main.go with message:
🚀 feat: add user authentication system

Generating commit for file: README.md
Committed README.md with message:
📝 docs: update installation instructions
```

## 📚 Notes

- Luna skips common binaries/images
- If the model returns an empty response, fallback is `update <file>`
- Supported prefixes: `feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `chore:`
- Uses the `GEMINI_API_KEY` environment variable

## 🛠️ Troubleshooting

- "Error: Set GEMINI_API_KEY using LunaApikey first"
  - Run `LunaApikey YOUR_KEY` and reopen the terminal
- "Error running git status"
  - Ensure you are inside a Git repository and Git is installed
- "No staged changes to commit"
  - Run `git add .` or stage specific files
- API key not working
  - Verify the key is valid and has access to Gemini 2.0 Flash

---

Made with ❤️ by hax — version 1.3 (Beta)

