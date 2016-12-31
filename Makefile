VERSION = 0.1.1
VERSION_INFO_FILE = ./version-info.go

BUILDER_PLUGIN_NAMES = customerimage
POSTPROCESSOR_PLUGIN_NAMES = customerimage-export

REPO_BASE           = github.com/DimensionDataResearch
REPO_ROOT           = $(REPO_BASE)/packer-plugins-ddcloud
BUILDERS_ROOT       = $(REPO_ROOT)/builders
POSTPROCESSORS_ROOT = $(REPO_ROOT)/post-processors
VENDOR_ROOT         = $(REPO_ROOT)/vendor

BIN_DIRECTORY = _bin
EXECUTABLE_PREFIX_SUFFIX = ddcloud
EXECUTABLE_PREFIX_BUILDER = packer-builder-$(EXECUTABLE_PREFIX_SUFFIX)
EXECUTABLE_PREFIX_POSTPROCESSOR = packer-post-processor-$(EXECUTABLE_PREFIX_SUFFIX)
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
dev: version fmt dev-builders dev-postprocessors

# Perform a development (current-platform-only) build of all builder plugins and publish them to ~/.packer.d/plugins.
dev-builders: version fmt $(foreach PLUGIN_NAME,$(BUILDER_PLUGIN_NAMES),dev-builder-$(PLUGIN_NAME))

# Perform a development (current-platform-only) build of the customer image builder plugin and publish it to ~/.packer.d/plugins.
dev-builder-customerimage: PLUGIN_NAME = customerimage
dev-builder-customerimage: _dev-builder

# Perform a development (current-platform-only) build of all builder plugins and publish them to ~/.packer.d/plugins.
dev-postprocessors: version fmt $(foreach PLUGIN_NAME,$(POSTPROCESSOR_PLUGIN_NAMES),dev-postprocessor-$(PLUGIN_NAME))

dev-postprocessor-customerimage-export: PLUGIN_NAME = customerimage-export
dev-postprocessor-customerimage-export: _dev-postprocessor

# Perform a full (all-platforms) build of all plugins.
build: version fmt build-builders build-postprocessors

# Perform a full (all-platforms) build of all builder plugins.
build-builders: version fmt $(foreach PLUGIN_NAME,$(BUILDER_PLUGIN_NAMES),build-builder-$(PLUGIN_NAME))

# Perform a full (all-platforms) build of the customer image builder plugin.
build-builder-customerimage: PLUGIN_NAME = customerimage
build-builder-customerimage: _build-builder

# Perform a full (all-platforms) build of all post-processor plugins.
build-postprocessors: version fmt $(foreach PLUGIN_NAME,$(POSTPROCESSOR_PLUGIN_NAMES),build-postprocessor-$(PLUGIN_NAME))

# Perform a full (all-platforms) build of the customer image export post-processor plugin.
build-postprocessor-customerimage-export: PLUGIN_NAME = customerimage-export
build-postprocessor-customerimage-export: _build-postprocessor

# Produce archives containing all plugins for a GitHub release.
dist: dist-builders dist-postprocessors

# Produce archives containing builder plugins for a GitHub release.
dist-builders: $(foreach PLUGIN_NAME,$(BUILDER_PLUGIN_NAMES),dist-builder-$(PLUGIN_NAME))

# Add / update the customer image builder plugin in archives for a GitHub release.
dist-builder-customerimage: PLUGIN_NAME = customerimage
dist-builder-customerimage: _dist-builder

# Produce archives containing post-processor plugins for a GitHub release.
dist-postprocessors: $(foreach PLUGIN_NAME,$(POSTPROCESSOR_PLUGIN_NAMES),dist-postprocessor-$(PLUGIN_NAME))

# Add / update the customer image post-processor plugin in archives for a GitHub release.
dist-postprocessor-customerimage-export: PLUGIN_NAME = customerimage-export
dist-postprocessor-customerimage-export: _dist-postprocessor

# Run most tests.
test: fmt # TODO: Add test targets

# Run all tests.
testall: 
	go test -v $(REPO_ROOT)/...

#################
# Private targets
#################

_dev-builder:
	@ \
	PLUGIN_NAME="$(PLUGIN_NAME)" \
	PLUGIN_FOLDER="$(BUILDERS_ROOT)/$(PLUGIN_NAME)" \
	EXECUTABLE_NAME="$(EXECUTABLE_PREFIX_BUILDER)-$(PLUGIN_NAME)" \
	make _dev

_dev-postprocessor:
	@ \
	PLUGIN_NAME="$(PLUGIN_NAME)" \
	PLUGIN_FOLDER="$(POSTPROCESSORS_ROOT)/$(PLUGIN_NAME)" \
	EXECUTABLE_NAME="$(EXECUTABLE_PREFIX_POSTPROCESSOR)-$(PLUGIN_NAME)" \
	make _dev

_dev:
	@mkdir -p ~/.packer.d/plugins
	go build -o $(BIN_DIRECTORY)/$(EXECUTABLE_NAME) $(PLUGIN_FOLDER)
	@cp $(BIN_DIRECTORY)/$(EXECUTABLE_NAME) ~/.packer.d/plugins
	@echo "Published plugin '$(EXECUTABLE_NAME)' to ~/.packer.d/plugins."

_build-builder:
	@ \
	PLUGIN_NAME="$(PLUGIN_NAME)" \
	PLUGIN_FOLDER="$(BUILDERS_ROOT)/$(PLUGIN_NAME)" \
	EXECUTABLE_NAME="$(EXECUTABLE_PREFIX_BUILDER)-$(PLUGIN_NAME)" \
	make _build

_build-postprocessor:
	@ \
	PLUGIN_NAME="$(PLUGIN_NAME)" \
	PLUGIN_FOLDER="$(POSTPROCESSORS_ROOT)/$(PLUGIN_NAME)" \
	EXECUTABLE_NAME="$(EXECUTABLE_PREFIX_POSTPROCESSOR)-$(PLUGIN_NAME)" \
	make _build

_build: _build-windows64 _build-windows32 _build-linux64 _build-mac64

_build-windows64: version
	GOOS=windows GOARCH=amd64 go build -o $(BIN_DIRECTORY)/windows-amd64/$(EXECUTABLE_NAME).exe $(PLUGIN_FOLDER)

_build-windows32: version
	GOOS=windows GOARCH=386 go build -o $(BIN_DIRECTORY)/windows-386/$(EXECUTABLE_NAME).exe $(PLUGIN_FOLDER)

_build-linux64: version
	GOOS=linux GOARCH=amd64 go build -o $(BIN_DIRECTORY)/linux-amd64/$(EXECUTABLE_NAME) $(PLUGIN_FOLDER)

_build-mac64: version
	GOOS=darwin GOARCH=amd64 go build -o $(BIN_DIRECTORY)/darwin-amd64/$(EXECUTABLE_NAME) $(PLUGIN_FOLDER)

_dist-builder:
	@ \
	PLUGIN_NAME="$(PLUGIN_NAME)" \
	PLUGIN_FOLDER="$(BUILDERS_ROOT)/$(PLUGIN_NAME)" \
	EXECUTABLE_NAME="$(EXECUTABLE_PREFIX_BUILDER)-$(PLUGIN_NAME)" \
	make _dist

_dist-postprocessor:
	@ \
	PLUGIN_NAME="$(PLUGIN_NAME)" \
	PLUGIN_FOLDER="$(POSTPROCESSORS_ROOT)/$(PLUGIN_NAME)" \
	EXECUTABLE_NAME="$(EXECUTABLE_PREFIX_POSTPROCESSOR)-$(PLUGIN_NAME)" \
	make _dist

_dist: _build
	@echo "Building distribution package for '$(EXECUTABLE_NAME)'..."
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
