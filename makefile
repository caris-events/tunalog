# Define the Go binary name
BINARY_NAME=tunalog

# Define the output directory
OUTPUT_DIR=bin

# Entry point
all: clean build

# Remove bin directory if it exists
clean:
	@if [ -d $(OUTPUT_DIR) ]; then rm -rf $(OUTPUT_DIR); fi

# Build all binaries
build:
	@echo "Building Windows x64..."
	GOOS=windows GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(BINARY_NAME)_windows_x64.exe && \
	\
	echo "Building macOS x64..." && \
	GOOS=darwin GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(BINARY_NAME)_macos_x64 && \
	\
	echo "Building Linux x64..." && \
	GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(BINARY_NAME)_linux_x64 && \
	\
	echo "Building Windows ARM64..." && \
	GOOS=windows GOARCH=arm64 go build -o $(OUTPUT_DIR)/$(BINARY_NAME)_windows_arm64.exe && \
	\
	echo "Building macOS ARM64..." && \
	GOOS=darwin GOARCH=arm64 go build -o $(OUTPUT_DIR)/$(BINARY_NAME)_macos_arm64 && \
	\
	echo "Building Linux ARM64..." && \
	GOOS=linux GOARCH=arm64 go build -o $(OUTPUT_DIR)/$(BINARY_NAME)_linux_arm64


