include ../CommonVars.inc

default: dev

# Perform a development (current-platform-only) build of all post-processor plugins and publish them to ~/.packer.d/plugins.
dev: $(foreach PLUGIN_NAME,$(BUILDER_PLUGIN_NAMES),dev-builder-$(PLUGIN_NAME))

# Perform a development (current-platform-only) build of the customer image builder plugin and publish it to ~/.packer.d/plugins.
dev-builder-customerimage:
	cd customerimage && make dev
