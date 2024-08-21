# Define the Go build command
GOBUILD=go build -v

# Define the Go install command
GOINSTALL=go install -v

# Define the target binary name
BINARY_NAME=trm

# Define the target OS and architectures
TARGET_OSARCHES=darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 linux/arm windows/amd64 windows/arm64

# Build for all OS and architectures
all:
	@for osarch in $(TARGET_OSARCHES); do \
		echo "Building for $$osarch"; \
		GOOS=$$(echo $$osarch | cut -d'/' -f1) GOARCH=$$(echo $$osarch | cut -d'/' -f2) $(GOBUILD) -o $(BINARY_NAME)-$$osarch ; \
	done

# Install for all OS and architectures
install:
	@for osarch in $(TARGET_OSARCHES); do \
		echo "Installing for $$osarch"; \
		GOOS=$$(echo $$osarch | cut -d'/' -f1) GOARCH=$$(echo $$osarch | cut -d'/' -f2) $(GOINSTALL) ; \
	done

# Clean up
clean:
	rm -rf $(BINARY_NAME)-*