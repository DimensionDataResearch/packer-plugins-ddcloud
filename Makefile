BUILD_ROOT = .

include CommonVars.inc

################
# Public targets
################

default: fmt build test

fmt:
	go fmt $(REPO_ROOT)/...

clean:
	rm -rf $(BIN_DIRECTORY) $(VERSION_INFO_FILE)
	go clean $(REPO_ROOT)/...

# Perform a development (current-platform-only) build of all plugins and publish them to ~/.packer.d/plugins.
dev: version fmt dev-builders dev-postprocessors

# Perform a development (current-platform-only) build of all plugins of a specific type and publish them to ~/.packer.d/plugins.
dev-%: version fmt
	cd $* && make dev

# Perform a full (all-platforms) build of all plugins.
build: version fmt build-builders build-postprocessors

# Perform a full (all-platforms) build of all plugins of a specific type.
build-%: version fmt
	cd $* && make build

# Produce archives containing all plugins for a GitHub release.
dist: dist-builders dist-postprocessors

# Produce archives containing all plugins of a specific type for a GitHub release.
dist-%:
	cd $* && make dist

# Run most tests.
test: fmt # TODO: Add test targets

# Run all tests.
testall: 
	go test -v $(REPO_ROOT)/...


version: $(VERSION_INFO_FILE)
$(VERSION_INFO_FILE): Makefile CommonTargets.inc
	@echo "Update version info: v$(VERSION)"
	@echo "package plugins\n\n// ProviderVersion is the current version of the CloudControl plugins for Packer.\nconst ProviderVersion = \"v$(VERSION) (`git rev-parse HEAD`)\"" > $(VERSION_INFO_FILE)
