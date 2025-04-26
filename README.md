# Go-Docs

A lightweight CLI tool to automatically generate API documentation from Go source code using OpenAI GPT.

---

## 🚀 Step 1: Install `godocs`

Install `godocs` directly from the GitHub repository using the Go install command:

```bash
go install github.com/JackBee2912/godocs@latest
```

✅ This command will download, build, and install the `godocs` binary into your `$GOPATH/bin` directory (usually `~/go/bin`).

---

## 🛠 Step 2: Ensure `$GOPATH/bin` is included in your system `$PATH`

After installation, make sure your `$GOPATH/bin` is included in your system's PATH environment variable.

**For zsh users (macOS, Linux, etc.):**

```bash
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
source ~/.zshrc
```

**For bash users:**

```bash
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc
```

✅ This ensures you can run the `godocs` command from anywhere in your terminal.

---

## ⚙️ Step 3: Run `godocs` to generate API documentation

Once installed and PATH is set, you can use the `godocs` CLI to generate documentation.

Basic usage:

```bash
godocs init --sourceRoot="../auth" --apiKey="your-openai-api-key"
```

Where:
- `--sourceRoot`: Path to your Go project source code folder
- `--apiKey`: Your OpenAI API Key (required)

---

## 📦 Example

```bash
godocs init --sourceRoot="../core" --apiKey="sk-xxxxxxxxxxxxxxxxxxxxxxxx"
```

✅ This will scan your Go source code, extract the functions, and auto-generate clean markdown API documentation using GPT.

---

## 📋 Troubleshooting

| Issue | Solution |
|:---|:---|
| `zsh: command not found: godocs` | Ensure your `$GOPATH/bin` is added to `$PATH` and you sourced your shell config (`source ~/.zshrc` or `source ~/.bashrc`). |
| `cannot read prompt.txt` | Use the latest version of `godocs` where prompt content is embedded into the binary (no external file needed). |
| OpenAI API errors (e.g., invalid API key) | Verify that your OpenAI API key is correct and active. |

---

## 🧹 Local development (for contributors)

If you want to modify or contribute to `godocs`:

**Clone the repository:**

```bash
git clone https://github.com/JackBee2912/godocs.git
cd godocs
```

**Install locally:**

```bash
make install
```

**Build manually:**

```bash
make build
```

**Clean build artifacts:**

```bash
make clean
```

---

## 📂 Project Structure

| Folder/File | Description |
|:---|:---|
| `/cmd` | CLI command definitions |
| `/internal/gpt` | GPT integration and prompt template |
| `/internal/markdown` | Markdown file generation utilities |
| `/internal/parser` | Go source code parser |
| `main.go` | Entry point for CLI application |
| `Makefile` | Build and install instructions |
| `go.mod` | Go module definition |

---

## 🏷️ Tags

#golang #cli #documentation #openai #gpt #automation #godocs

---

## ✨ License

This project is open-source under the MIT License.
