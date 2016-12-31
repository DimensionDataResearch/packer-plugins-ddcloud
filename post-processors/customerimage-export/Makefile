include ../../CommonVars.inc
include ../../CommonTargets.inc

default: dev

# Perform a development (current-platform-only) build of the customer image export post-processor plugin and publish it to ~/.packer.d/plugins.
dev: PLUGIN_NAME = customerimage-export
dev: PLUGIN_FOLDER = $(POSTPROCESSORS_ROOT)/$(PLUGIN_NAME)
dev: EXECUTABLE_NAME = $(EXECUTABLE_PREFIX_POSTPROCESSOR)-$(PLUGIN_NAME)
dev: _dev
