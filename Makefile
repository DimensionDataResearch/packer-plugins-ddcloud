VERSION = 0.1.0
VERSION_INFO_FILE = ./version-info.go

BUILDER_PLUGIN_NAMES = customerimage

REPO_BASE     = github.com/DimensionDataResearch
REPO_ROOT     = $(REPO_BASE)/packer-plugins-ddcloud
BUILDERS_ROOT = $(REPO_ROOT)/builders
VENDOR_ROOT   = $(REPO_ROOT)/vendor

BIN_DIRECTORY = _bin
EXECUTABLE_PREFIX_SUFFIX = ddcloud
EXECUTABLE_PREFIX_BUILDER = packer-builder-$(EXECUTABLE_PREFIX_SUFFIX)
DIST_ZIP_PREFIX = packer-plugins-ddcloud-v$(VERSION)

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
dev: version fmt dev-builders

# Perform a development (current-platform-only) build of all builder plugins and publish them to ~/.packer.d/plugins.
dev-builders: version fmt $(foreach PLUGIN_NAME,$(BUILDER_PLUGIN_NAMES),dev-builder-$(PLUGIN_NAME))

# Perform a development (current-platform-only) build of the customer image builder plugin and publish it to ~/.packer.d/plugins.
dev-builder-customerimage: PLUGIN_NAME = customerimage
dev-builder-customerimage: -dev-builder

# Perform a full (all-platforms) build of all plugins.
build: version fmt build-builders

# Perform a full (all-platforms) build of all builder plugins.
build-builders: version fmt $(foreach PLUGIN_NAME,$(BUILDER_PLUGIN_NAMES),build-builder-$(PLUGIN_NAME))

# Perform a full (all-platforms) build of the customer image builder plugin.
build-builder-customerimage: EXECUTABLE_NAME = $(EXECUTABLE_PREFIX_BUILDER)-customerimage
build-builder-customerimage: -build

# Produce archives containing all plugins for a GitHub release.
dist: dist-builders

# Produce archives containing builder plugins for a GitHub release.
dist-builders: $(foreach PLUGIN_NAME,$(BUILDER_PLUGIN_NAMES),dist-builder-$(PLUGIN_NAME))

# Add / update the customer image builder plugin in archives for a GitHub release.
dist-builder-customerimage: PLUGIN_NAME = customerimage
dist-builder-customerimage: -dist-builder

# Run most tests.
test: fmt # TODO: Add test targets

# Run all tests.
testall: 
	go test -v $(REPO_ROOT)/...

#################
# Private targets
#################

-dev-builder: PLUGIN_FOLDER = $(BUILDERS_ROOT)/$(PLUGIN_NAME)
-dev-builder: EXECUTABLE_NAME = $(EXECUTABLE_PREFIX_BUILDER)-$(PLUGIN_NAME)
-dev-builder: -dev

-dev:
	@mkdir -p ~/.packer.d/plugins
	go build -o $(BIN_DIRECTORY)/$(EXECUTABLE_NAME) $(PLUGIN_FOLDER)
	@cp $(BIN_DIRECTORY)/$(EXECUTABLE_NAME) ~/.packer.d/plugins
	@echo "Published plugin '$(EXECUTABLE_NAME)' to ~/.packer.d/plugins."

-build-builder: PLUGIN_FOLDER = $(BUILDERS_ROOT)/$(PLUGIN_NAME)
-build-builder: EXECUTABLE_NAME = $(EXECUTABLE_PREFIX_BUILDER)-$(PLUGIN_NAME)
-build-builder: -build

-build: -build-windows64 -build-windows32 -build-linux64 -build-mac64

-build-windows64: version
	GOOS=windows GOARCH=amd64 go build -o $(BIN_DIRECTORY)/windows-amd64/$(EXECUTABLE_NAME).exe

-build-windows32: version
	GOOS=windows GOARCH=386 go build -o $(BIN_DIRECTORY)/windows-386/$(EXECUTABLE_NAME).exe

-build-linux64: version
	GOOS=linux GOARCH=amd64 go build -o $(BIN_DIRECTORY)/linux-amd64/$(EXECUTABLE_NAME)

-build-mac64: version
	GOOS=darwin GOARCH=amd64 go build -o $(BIN_DIRECTORY)/darwin-amd64/$(EXECUTABLE_NAME)

-dist-builder: PLUGIN_FOLDER = $(BUILDERS_ROOT)/$(PLUGIN_NAME)
-dist-builder: EXECUTABLE_NAME = $(EXECUTABLE_PREFIX_BUILDER)-$(PLUGIN_NAME)
-dist-builder: -dist

-dist: -build
	cd $(BIN_DIRECTORY)/windows-386 && \
		zip -9 ../$(DIST_ZIP_PREFIX).windows-386.zip $(EXECUTABLE_NAME).exe
	cd $(BIN_DIRECTORY)/windows-amd64 && \
		zip -9 ../$(DIST_ZIP_PREFIX).windows-amd64.zip $(EXECUTABLE_NAME).exe
	cd $(BIN_DIRECTORY)/linux-amd64 && \
		zip -9 ../$(DIST_ZIP_PREFIX).linux-amd64.zip $(EXECUTABLE_NAME)
	cd $(BIN_DIRECTORY)/darwin-amd64 && \
		zip -9 ../$(DIST_ZIP_PREFIX)-darwin-amd64.zip $(EXECUTABLE_NAME)

version: $(VERSION_INFO_FILE)

$(VERSION_INFO_FILE): Makefile
	@echo "Update version info: v$(VERSION)"
	@echo "package plugins\n\n// ProviderVersion is the current version of the CloudControl plugins for Packer.\nconst ProviderVersion = \"v$(VERSION) (`git rev-parse HEAD`)\"" > $(VERSION_INFO_FILE)
