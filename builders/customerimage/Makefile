include ../../CommonVars.inc
include ../../CommonTargets.inc

default: dev

# Perform a development (current-platform-only) build of the customer image builder plugin and publish it to ~/.packer.d/plugins.
dev: PLUGIN_NAME = customerimage
dev: PLUGIN_FOLDER = $(BUILDERS_ROOT)/$(PLUGIN_NAME)
dev: EXECUTABLE_NAME = $(EXECUTABLE_PREFIX_BUILDER)-$(PLUGIN_NAME)
dev: _dev
