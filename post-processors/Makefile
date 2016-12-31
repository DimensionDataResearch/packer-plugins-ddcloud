include ../CommonVars.inc

default: dev

# Perform a development (current-platform-only) build of all post-processor plugins and publish them to ~/.packer.d/plugins.
dev: $(foreach PLUGIN_NAME,$(POSTPROCESSOR_PLUGIN_NAMES),dev-postprocessor-$(PLUGIN_NAME))

# Perform a development (current-platform-only) build of the customer image export post-processor plugin and publish it to ~/.packer.d/plugins.
dev-postprocessor-customerimage-export:
	cd customerimage-export && make dev
