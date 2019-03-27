# --- Variables ---

# The cross-compiled program binaries:
BIN_LINUX_AMD64=build/linux-amd64/mjpeg-server
BIN_DARWIN_AMD64=build/darwin-amd64/mjpeg-server
BIN_WINDOWS_AMD64=build/windows-amd64/MJPEGServer.exe

# The release artifacts:
RELEASE_LINUX_AMD64=build/mjpeg-server-linux-amd64.tar.gz
RELEASE_DARWIN_AMD64=build/mjpeg-server-darwin-amd64.zip
RELEASE_WINDOWS_AMD64=build/mjpeg-server-windows-amd64.zip
RELEASES=$(RELEASE_LINUX_AMD64) $(RELEASE_DARWIN_AMD64) $(RELEASE_WINDOWS_AMD64)

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

# The default target builds release artifacts for all platforms:
all: $(RELEASES)

# Runs the unit tests for all components:
test:
	@go test ./...

# Releases the packaged binaries on GitHub:
release: $(RELEASES)
	@bin/github-release.sh $(RELEASES)

# Removes all build artifacts:
clean:
	@rm -rf build


# --- Helper targets ---

# Defines phony targets (targets without a corresponding target file):
.PHONY: \
	all \
	test \
	release \
	clean

# Builds the Linux binary:
$(BIN_LINUX_AMD64): $(DEPS)
	mkdir -p $(dir $(BIN_LINUX_AMD64))
	GOOS=linux GOARCH=amd64 go build -ldflags="$(LD_FLAGS)" -o $@

# Builds the MacOS binary:
$(BIN_DARWIN_AMD64): $(DEPS)
	mkdir -p $(dir $(BIN_DARWIN_AMD64))
	GOOS=darwin GOARCH=amd64 go build -ldflags="$(LD_FLAGS)" -o $@

# Builds the Windows binary:
$(BIN_WINDOWS_AMD64): $(DEPS)
	mkdir -p $(dir $(BIN_WINDOWS_AMD64))
	GOOS=windows GOARCH=amd64 go build -ldflags="$(LD_FLAGS)" -o $@

# Builds the Linux release artifact:
$(RELEASE_LINUX_AMD64): $(BIN_LINUX_AMD64)
	tar -czf $@ -C $(dir $(BIN_LINUX_AMD64)) $(notdir $(BIN_LINUX_AMD64))

# Builds the MacOS release artifact:
$(RELEASE_DARWIN_AMD64): $(BIN_DARWIN_AMD64)
	zip -jq $@ $(BIN_DARWIN_AMD64)

# Builds the Windows release artifact:
$(RELEASE_WINDOWS_AMD64): $(BIN_WINDOWS_AMD64)
	zip -jq $@ $(BIN_WINDOWS_AMD64)
