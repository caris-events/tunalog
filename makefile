# Define the Go binary name
BINARY_NAME=tunalog

# Define the output directory
OUTPUT_DIR=bin

# Build targets
build_windows_amd64:
	GOOS=windows GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(BINARY_NAME)_windows_x64.exe

build_macos_amd64:
	GOOS=darwin GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(BINARY_NAME)_macos_x64

build_linux_amd64:
	GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(BINARY_NAME)_linux_x64

build_windows_arm64:
	GOOS=windows GOARCH=arm64 go build -o $(OUTPUT_DIR)/$(BINARY_NAME)_windows_arm64.exe

build_macos_arm64:
	GOOS=darwin GOARCH=arm64 go build -o $(OUTPUT_DIR)/$(BINARY_NAME)_macos_arm64

build_linux_arm64:
	GOOS=linux GOARCH=arm64 go build -o $(OUTPUT_DIR)/$(BINARY_NAME)_linux_arm64

# Remove bin directory if it exists
clean:
	@if [ -d $(OUTPUT_DIR) ]; then rm -rf $(OUTPUT_DIR); fi
	
# Default target (build all)
all: build_windows_amd64 build_macos_amd64 build_linux_amd64 build_windows_arm64 build_macos_arm64 build_linux_arm64