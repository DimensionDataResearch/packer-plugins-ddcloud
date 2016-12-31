BUILD_ROOT = ..

include ../CommonVars.inc

default: dev

# Perform a development (current-platform-only) build of all post-processor plugins and publish them to ~/.packer.d/plugins.
dev: $(foreach PLUGIN_NAME,$(POSTPROCESSOR_PLUGIN_NAMES),dev-postprocessor-$(PLUGIN_NAME))

# Perform a development (current-platform-only) build of a specific post-processor plugin and publish it to ~/.packer.d/plugins.
dev-postprocessor-%:
	cd $* && make dev

# Perform an all-platforms build of all post-processor plugins.
build: $(foreach PLUGIN_NAME,$(POSTPROCESSOR_PLUGIN_NAMES),build-postprocessor-$(PLUGIN_NAME))

# Perform an all-platforms build of a specific post-processor plugin.
build-postprocessor-%:
	cd $* && make build

# Produce archives containing all post-processor plugins for a GitHub release.
dist: $(foreach PLUGIN_NAME,$(POSTPROCESSOR_PLUGIN_NAMES),dist-postprocessor-$(PLUGIN_NAME))

# Produce an archive containing a specific plugin for a GitHub release.
dist-postprocessor-%: build-postprocessor-%
	cd $* && make dist
