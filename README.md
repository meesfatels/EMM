# EMM (Eidolon Minion Manager)

Modular Go-based CLI/TUI for OpenRouter AI models. Decouples "Agent" personas (instincts) from "Minion" model configurations.

## Installation

### Arch Linux (AUR)
Install using `yay`:
```bash
yay -S emm-git
```

### macOS (Homebrew)
Install directly via the formula:
```bash
brew install --formula https://raw.githubusercontent.com/meesfatels/EMM/main/packaging/brew/emm.rb
```

### Manual Installation (Go)
```bash
go install github.com/meesfatels/emm/cmd/emm@latest
```

## Setup
Once installed, initialize your configuration directory:
```bash
emm init
```
Edit `~/.emm/emm.yaml` and add your **OpenRouter API Key**.

## Usage
Start a session with an agent and minion:
```bash
emm run --agent example --minion example
```
Inside the TUI:
- Type your message and press **Enter** to send.
- Use **`/save name`** to save your session.
- Use **`/load name`** to continue a saved session.
- Press **Ctrl+C** to exit.
