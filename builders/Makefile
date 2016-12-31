BUILD_ROOT = ..

include ../CommonVars.inc

default: dev

# Perform a development (current-platform-only) build of all post-processor plugins and publish them to ~/.packer.d/plugins.
dev: $(foreach PLUGIN_NAME,$(BUILDER_PLUGIN_NAMES),dev-builder-$(PLUGIN_NAME))

# Perform a development (current-platform-only) build a specific builder plugin and publish it to ~/.packer.d/plugins.
dev-builder-%:
	cd $* && make dev

# Perform an all-platforms build of all builder plugins.
build: $(foreach PLUGIN_NAME,$(BUILDER_PLUGIN_NAMES),build-builder-$(PLUGIN_NAME))

# Perform an all-platforms build of a specific plugin.
build-builder-%:
	cd $* && make build

# Produce archives containing all builder plugins for a GitHub release.
dist: $(foreach PLUGIN_NAME,$(BUILDER_PLUGIN_NAMES),dist-builder-$(PLUGIN_NAME))

# Produce an archive containing a specific builder plugin for a GitHub release.
dist-builder-%: build-builder-%
	cd $* && make dist
