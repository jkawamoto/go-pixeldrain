# A Pixeldrain client
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
[![Go application](https://github.com/jkawamoto/go-pixeldrain/actions/workflows/ci.yaml/badge.svg)](https://github.com/jkawamoto/go-pixeldrain/actions/workflows/ci.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/jkawamoto/go-pixeldrain.svg)](https://pkg.go.dev/github.com/jkawamoto/go-pixeldrain)
[![codecov](https://codecov.io/gh/jkawamoto/go-pixeldrain/branch/master/graph/badge.svg?token=ppX3MVIqWA)](https://codecov.io/gh/jkawamoto/go-pixeldrain)
[![Release](https://img.shields.io/badge/release-0.7.3-brightgreen.svg)](https://github.com/jkawamoto/go-pixeldrain/releases/tag/v0.7.3)


## Usage
### Upload files
```shell
pd upload <path[:name]>...
```

`upload` command uploads files specified by the given `path`s to Pixeldrain and shows a URL to download them.
Each `path` can have an optional `name`. If a name is given, uploaded file will be renamed with it.

For example, this command reads `img.png` and uploads it as `another.png`:

```shell
pd upload img.png:another.png
```

If `path` is `-`, the uploading file is read from stdin. In this case, it's recommended to give a file name.
To avoid being interpreted as an option flag, it is necessary to prepend `--` before the argument in this case.
For example, this command reads data from stdin and uploads it as `output.log`:

```shell
pd upload -- -:output.log
```

If multiple files are given, an album consists of them will be created. By default, the album has a random name.
`--album` flag can specify the name.
For example, this command uploads two files and creates an album named `screenshots`:

```shell
pd upload --album screenshots img1.png img2.png
```

#### Upload a directory
Since this application supports uploading a file from STDIN, you can upload a directory with `tar` command.
To ensure the argument is not interpreted as an option flag, prepend `--` before the argument.
For example, this command uploads `~/Documents` directory:

```shell
tar zcf - ~/Documents | pd upload -- -:documents.tar.gz
```

#### Upload to your account
If you want to upload files to your account, give your API key with `--api-key` flag or via `PIXELDRAIN_API_KEY`
environment variable.

An API key can be obtained from https://pixeldrain.com/user/api_keys.


### Download files
```shell
pd download <URL>...
```

`download` command downloads files from Pixeldrain and stores it in the current directory by default.

If `--dir` or `-o` option is given with a directory path, the downloaded file is stored in the directory.

If the given URL refers an album which consists of multiple files, this command asks which file you want to download.
If you want to download all files without any interaction, use `--all` flag.

### End-to-end encryption
If recipients are specified with `--recipient` and/or `--recipient-file` flags to upload command,
files will be encrypted before being uploaded by [age](https://github.com/FiloSottile/age).
Encrypted files will have extension `.age`.

A recipient specified with `--recipient` flag can be an age public key generated by `age-keygen` ("age1...")
or an SSH public key ("ssh-ed25519 AAAA...", "ssh-rsa AAAA...").
A recipient file specified with `--recipient-file` flag contains one or more recipients, one per line.
Empty line sand lines starting with "#" are ignored as comments.

If a downloading file has extension `.age` and an identity file is specified with `--identity` flag to download command,
the file will be decrypted.

An identity contains one or more secret keys ("AGE-SECRET-KEY-1..."), one per line, or an SSH key.
Empty lines and lines starting with "#" are ignored as comments.

See [age](https://github.com/FiloSottile/age) for the details of `age` and `age-keygen`.

## Installation
If you're a Homebrew or Linuxbrew user, you can install this app by the following commands:

```
$ brew tap jkawamoto/pixeldrain
$ brew install pixeldrain
```

To build the newest version, use go get command:

```
$ go install github.com/jkawamoto/go-pixeldrain/cmd/pd@latest
```

Otherwise, compiled binaries are also available in [GitHub](https://github.com/jkawamoto/go-pixeldrain/releases).


## License
This software is released under the MIT License, see [LICENSE](LICENSE).
