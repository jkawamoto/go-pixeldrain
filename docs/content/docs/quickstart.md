---
title: "Quick Start & CLI"
weight: 2
---

Learn how to install and use the `pd` CLI for day-to-day file transfers, album creation, and pipelines.

## Installation

### Homebrew (macOS & Linux)

The easiest way to install `pd` on macOS or Linux is via Homebrew:

```shell
brew tap jkawamoto/pixeldrain
brew install pixeldrain
```

### Go Install

If you have Go installed (Go 1.24 or later recommended):

```shell
go install github.com/jkawamoto/go-pixeldrain/cmd/pd@latest
```

### Binary Releases

Download pre-compiled binaries for macOS, Linux, or Windows from the [GitHub Releases](https://github.com/jkawamoto/go-pixeldrain/releases) page.

---

## Authentication (Optional)

`pd` works anonymously out-of-the-box. However, if you have a Pixeldrain account, configuring your API key allows you to manage files under your user account and keep track of your uploads.

Get your API key at [pixeldrain.com/user/api_keys](https://pixeldrain.com/user/api_keys).

You can supply the API key in either of two ways:

1. **Environment Variable** (Recommended):
   ```shell
   export PIXELDRAIN_API_KEY="your-api-key"
   ```

2. **Command Flag**:
   ```shell
   pd --api-key "your-api-key" upload file.zip
   ```

---

## Basic CLI Commands

### Upload Files

To upload a file and print the public link:

```shell
pd upload artwork.png
# => https://pixeldrain.com/u/xxxxxx
```

To specify a custom display name on Pixeldrain using `path:name`:

```shell
pd upload output_0042.png:cyberpunk-city.png
```

### Create Albums

Combine multiple files into a single album URL using `--album`:

```shell
pd upload --album "Neon Dreams" concept1.png concept2.png concept3.png
# => https://pixeldrain.com/l/xxxxxx
```

### Download Files

Download a file by URL or ID (saved to current directory by default):

```shell
pd download https://pixeldrain.com/u/xxxxxx
```

Specify a custom destination directory or filename:

```shell
pd download -o ./downloads https://pixeldrain.com/u/xxxxxx
pd download --filename custom.png https://pixeldrain.com/u/xxxxxx
```

To download all files from an album non-interactively:

```shell
pd download --all https://pixeldrain.com/l/xxxxxx
```

---

## Advanced Unix Pipelines

### Uploading from STDIN

Upload log files, git diffs, or output from any program by providing `-` as the path (use `--` before `-` so flags aren't confused):

```shell
# Upload build logs
cat build.log | pd upload -- -:build.log

# Upload a git patch
git diff | pd upload -- -:patch.diff
```

### Uploading Entire Directories

Stream directory archives directly to Pixeldrain without intermediate disk writes:

```shell
tar zcf - ~/Projects/my-app | pd upload -- -:my-app.tar.gz
```

And download/extract in a single pipeline:

```shell
pd download <file-id> | tar zxf - -C ~/Downloads
```
