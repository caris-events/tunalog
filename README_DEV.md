# Tunalog Development Notes

### Using Go Conditional Compilation

Compile All at once:

```
GOOS=windows GOARCH=amd64 go build -o bin/x64_windows_tunalog.exe &&
GOOS=darwin GOARCH=amd64 go build -o bin/x64_macos_tunalog &&
GOOS=linux GOARCH=amd64 go build -o bin/x64_linux_tunalog &&
GOOS=windows GOARCH=arm64 go build -o bin/arm64_windows_tunalog.exe &&
GOOS=darwin GOARCH=arm64 go build -o bin/arm64_macos_tunalog &&
GOOS=linux GOARCH=arm64 go build -o bin/arm64_linux_tunalog
```

or Compile for 64-bit Windows, macOS, Linux:

```
GOOS=windows GOARCH=amd64 go build -o bin/x64_windows_tunalog.exe &&
GOOS=darwin GOARCH=amd64 go build -o bin/x64_macos_tunalog &&
GOOS=linux GOARCH=amd64 go build -o bin/x64_linux_tunalog
```

or for arm64 Windows, macOS, Linux:

```
GOOS=windows GOARCH=arm64 go build -o bin/arm64_windows_tunalog.exe &&
GOOS=darwin GOARCH=arm64 go build -o bin/arm64_macos_tunalog &&
GOOS=linux GOARCH=arm64 go build -o bin/arm64_linux_tunalog
```

### Using Makefile

You can also use `make` command to build the whole project. The compiled binary will be exported to the `bin` folder.

```
cd tunalog
make
```

