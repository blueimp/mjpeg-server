# --- Variables ---

# The cross-compiled program binaries:
BIN_LINUX_AMD64=bin/mjpeg-server-linux
BIN_DARWIN_AMD64=bin/mjpeg-server-macos
BIN_WINDOWS_AMD64=bin/MJPEGServer.exe
BINS=$(BIN_LINUX_AMD64) $(BIN_DARWIN_AMD64) $(BIN_WINDOWS_AMD64)

# Dependencies:
DEP_MULTI = internal/multi/multi.go
DEP_RECORDING = internal/recording/recording.go
DEP_REQUEST = internal/request/request.go
DEPS = $(DEP_MULTI) $(DEP_RECORDING) $(DEP_REQUEST) main.go

# Use the git tag for the current commit as version or "dev" as fallback:
GET_VERSION=git describe --exact-match --tags 2> /dev/null || echo dev

# Set the program version and disable symbol table and DWARF generation:
LD_FLAGS=-X main.Version=$$($(GET_VERSION)) -s -w


# --- Main targets ---

# The default target builds binaries for all platforms:
all: $(BINS)

# Runs the unit tests for all components:
test:
	@go test ./...

# Releases the binaries on GitHub:
release: all
	@bin/github-release.sh $(BINS)

# Removes all build artifacts:
clean:
	@rm -f $(BINS)


# --- Helper targets ---

# Defines phony targets (targets without a corresponding target file):
.PHONY: \
	all \
	test \
	release \
	clean

# Builds the Linux binary:
$(BIN_LINUX_AMD64): $(DEPS)
	GOOS=linux GOARCH=amd64 go build -ldflags="$(LD_FLAGS)" -o $@

# Builds the MacOS binary:
$(BIN_DARWIN_AMD64): $(DEPS)
	GOOS=darwin GOARCH=amd64 go build -ldflags="$(LD_FLAGS)" -o $@

# Builds the Windows binary:
$(BIN_WINDOWS_AMD64): $(DEPS)
	GOOS=windows GOARCH=amd64 go build -ldflags="$(LD_FLAGS)" -o $@
