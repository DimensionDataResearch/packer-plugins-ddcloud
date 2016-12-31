_dev:
	@mkdir -p ~/.packer.d/plugins
	go build -o $(BIN_DIRECTORY)/$(EXECUTABLE_NAME) $(PLUGIN_FOLDER)
	@cp $(BIN_DIRECTORY)/$(EXECUTABLE_NAME) ~/.packer.d/plugins
	@echo "Published plugin '$(EXECUTABLE_NAME)' to ~/.packer.d/plugins."
