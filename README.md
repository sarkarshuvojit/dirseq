
# dirseq (Directory Sequencer)


A simple, fast command-line tool written in Go to get the next sequence number unique to the current directory. Useful for sequentially naming files, directories, or tracking versions within specific project folders.

## 🤔 Why?

Ever needed to create files or directories with incrementing numbers within a specific project, like:

* `01-initial-experiment`
* `02-revised-approach`
* `output-run-1.log`
* `output-run-2.log`
* `session_1`
* `session_2`

Manually checking the last number used is tedious and error-prone. `dirseq` solves this by automatically tracking and providing the *next* available sequence number for whatever directory you're currently in.

## ✨ Features

* **Per-Directory Sequences:** Each directory maintains its own independent counter.
* **Persistent:** Remembers the last sequence number used in each directory across sessions.
* **Simple CLI:** Just run `dirseq` to get the next number.
* **Fast:** Written in Go with a minimal SQLite backend.
* **Cross-Platform:** Works on Linux, macOS, and Windows (wherever Go runs).

## 🚀 Installation

There are a few ways to install `dirseq`:

**1. Using `go install` (Recommended if you have Go installed):**

Make sure you have Go version 1.21 or later installed.

```bash
go install github.com/sarkarshuvojit/dirseq@latest
```

**2. From GitHub Releases (If Pre-built Binaries are Provided):**

* Go to the [Releases Page](https://github.com/sarkarshuvojit/dirseq/releases)
* Download the pre-compiled binary suitable for your operating system and architecture.
* Place the downloaded binary somewhere in your system's `PATH` (e.g., `/usr/local/bin` on Linux/macOS, or a custom directory you've added to PATH on Windows).
* Make sure the binary is executable (`chmod +x dirseq` on Linux/macOS).

**3. Build from Source:**

See the [Building from Source](#️-building-from-source) section below.

## ⚙️ Usage

Using `dirseq` is straightforward:

1.  Navigate to the directory where you want a sequence number.
2.  Run the command:

    ```bash
    dirseq
    ```

3.  The tool will print the *next* available sequence number for that specific directory to standard output.

* The first time you run `dirseq` in a new directory, it will output `1`.
* The second time in the *same* directory, it will output `2`, and so on.
* Running `dirseq` in a *different* directory will start its own sequence, outputting `1` on the first run there.

Errors are printed to standard error using structured logging.

## 💡 Examples

`dirseq` is powerful when combined with other shell commands using command substitution (`$()` in bash/zsh, `$(...)` or backticks in other shells).

**Example 1: Basic Sequence**

```bash
# Navigate to your project folder
cd ~/projects/my-cool-project

# Get the first sequence number
dirseq
# Output: 1

# Get the next sequence number
dirseq
# Output: 2

# Navigate to another project
cd ~/projects/another-widget

# Get the first sequence number for *this* directory
dirseq
# Output: 1

# Go back to the first project
cd ~/projects/my-cool-project

# Get the next number (it remembers where you left off)
dirseq
# Output: 3
```

**Example 2: Creating Sequentially Numbered Directories**

This is the example you requested!

```bash
# In your project directory
mkdir $(dirseq)-initial-setup
# Creates directory: 1-initial-setup

mkdir $(dirseq)-feature-branch-work
# Creates directory: 2-feature-branch-work

mkdir $(dirseq)-final-refactor
# Creates directory: 3-final-refactor
```

**Example 3: Creating Sequentially Numbered Files**

```bash
# Save experiment results with a sequence number
run_experiment > results-$(dirseq).log
# Creates file: results-1.log

run_experiment --with-new-param > results-$(dirseq).log
# Creates file: results-2.log
```

**Example 4: Versioning Output Artifacts**

```bash
# Build your project and tag the output with the sequence
build_my_app --output release-build-$(dirseq).zip
# Creates file: release-build-1.zip

# After some changes...
build_my_app --output release-build-$(dirseq).zip
# Creates file: release-build-2.zip
```

## 🛠️ How it Works

`dirseq` stores its state in a simple SQLite database file located at:

`$HOME/.config/dirseq/mem.db`

* It automatically creates the `.config/dirseq` directory if it doesn't exist.
* The `mem.db` file contains a single table (default name `sequences`) with two columns:
    * `abs_path` (VARCHAR, PRIMARY KEY): Stores the absolute path of the directory.
    * `last_seq` (INT): Stores the last sequence number *used* for that path.
* When you run `dirseq`, it:
    1.  Gets the absolute path of the current working directory.
    2.  Looks up this path in the database.
    3.  If found, retrieves `last_seq`. If not found, `last_seq` defaults to `0`.
    4.  Calculates `newSeq = last_seq + 1`.
    5.  Prints `newSeq` to standard output.
    6.  Updates the database, storing `newSeq` as the `last_seq` for the current path (either by updating the existing row or inserting a new one).

## 🤝 Contributing

Contributions are welcome! Please feel free to open an issue or submit a pull request.

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
