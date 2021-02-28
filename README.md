# A Pixeldrain client
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
[![CircleCI](https://circleci.com/gh/jkawamoto/go-pixeldrain.svg?style=svg)](https://circleci.com/gh/jkawamoto/go-pixeldrain)
[![codecov](https://codecov.io/gh/jkawamoto/go-pixeldrain/branch/master/graph/badge.svg?token=ppX3MVIqWA)](https://codecov.io/gh/jkawamoto/go-pixeldrain)
[![Release](https://img.shields.io/badge/release-0.4.0-brightgreen.svg)](https://github.com/jkawamoto/go-pixeldrain/releases/tag/v0.4.0)


## Usage
### Upload a file
`upload` command uploads a file to Pixeldrain and shows a URL to it.


```shell-session
$ pd upload <file path>
```

The uploaded file has the same name as the given file.
`-n` and `--name` options overwrite the file names.


If uploading file is given via STDIN, use `-` instead of a file path.
In this case either `-n` or `--name` option is mandatory.

For example, this command reads `file1.txt` and uploads it with name `uploaded.txt`:

```shell-session
$ cat file1.txt | pd upload --name uploaded.txt -
```


### Download a file
`download` command downloads a file from Pixeldrain and writes it to STDOUT.

```shell-session
$ pd download <file ID | URL>
```

If `-o` option is given with a directory path, the downloaded file is stored in
the directory instead of writing to STDOUT.

For example, this command downloads a file `abcdefg` in `~/Download`:
```shell-session
$ pd download abcdefg -o ~/Download
```

### Upload/Download a directory
This application supports uploading a file from STDIN and downloading a file to STDOUT.
With `tar` command, it's also able to upload/download directories.
For example, this command uploads `~/Documents` directory:

```shell-session
$ tar zcf - ~/Documents | pd upload -n documents.tar.gz -
```

and this command downloads the file:

```shell-session
$ pd download <file id> | tar zxf - -C ~/Downloads
```



## Installation
If you're a Homebrew or Linuxbrew user, you can install this app by the following commands:

```
$ brew tap jkawamoto/pixeldrain
$ brew install pixeldrain
```

To build the newest version, use go get command:

```
$ go get github.com/jkawamoto/go-pixeldrain
```

Otherwise, compiled binaries are also available in [Github](https://github.com/jkawamoto/go-pixeldrain/releases).


## License
This software is released under the MIT License, see [LICENSE](LICENSE.md).
