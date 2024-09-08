# Tunalog Development Notes

### Using Go Conditional Compilation

Compile All at once:

```
GOOS=windows GOARCH=amd64 go build -o bin/tunalog_windows_x64.exe &&
GOOS=darwin GOARCH=amd64 go build -o bin/tunalog_macos_x64 &&
GOOS=linux GOARCH=amd64 go build -o bin/tunalog_linux_x64 &&
GOOS=windows GOARCH=arm64 go build -o bin/tunalog_windows_arm64.exe &&
GOOS=darwin GOARCH=arm64 go build -o bin/tunalog_macos_arm64 &&
GOOS=linux GOARCH=arm64 go build -o bin/tunalog_linux_arm64
```

or Compile for 64-bit Windows, macOS, Linux:

```
GOOS=windows GOARCH=amd64 go build -o bin/tunalog_windows_x64.exe &&
GOOS=darwin GOARCH=amd64 go build -o bin/tunalog_macos_x64 &&
GOOS=linux GOARCH=amd64 go build -o bin/tunalog_linux_x64
```

or for arm64 Windows, macOS, Linux:

```
GOOS=windows GOARCH=arm64 go build -o bin/tunalog_windows_arm64.exe &&
GOOS=darwin GOARCH=arm64 go build -o bin/tunalog_macos_arm64 &&
GOOS=linux GOARCH=arm64 go build -o bin/tunalog_linux_arm64
```

### Using Makefile

You can also use `make` command to build the whole project. The compiled binary will be exported to the `bin` folder.

```bash
$ cd tunalog
$ make
```

### Using Docker

```bash
$ docker build -t yamiodymel/tunalog:latest -t yamiodymel/tunalog:v1.0.2 .
$ docker push yamiodymel/tunalog
$ docker push yamiodymel/tunalog:v1.0.2
```
