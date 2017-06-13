VERSION = 0.1.3-beta2
VERSION_INFO_FILE = ./version-info.go

BUILDER_PLUGIN_NAMES = customerimage customerimage-import
POSTPROCESSOR_PLUGIN_NAMES = customerimage-export customerimage-import

REPO_BASE           = github.com/DimensionDataResearch
REPO_ROOT           = $(REPO_BASE)/packer-plugins-ddcloud
BUILDERS_ROOT       = $(REPO_ROOT)/builders
POSTPROCESSORS_ROOT = $(REPO_ROOT)/postprocessors
VENDOR_ROOT         = $(REPO_ROOT)/vendor

BIN_DIRECTORY                   = $(BUILD_ROOT)/_bin
EXECUTABLE_PREFIX_SUFFIX        = ddcloud
EXECUTABLE_PREFIX_BUILDER       = packer-builder-$(EXECUTABLE_PREFIX_SUFFIX)
EXECUTABLE_PREFIX_POSTPROCESSOR = packer-post-processor-$(EXECUTABLE_PREFIX_SUFFIX)
DIST_ZIP_PREFIX                 = packer-plugins-ddcloud-v$(VERSION)
