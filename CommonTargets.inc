_dev:
	@mkdir -p ~/.packer.d/plugins
	go build -o $(BIN_DIRECTORY)/$(EXECUTABLE_NAME) $(PLUGIN_FOLDER)
	@cp $(BIN_DIRECTORY)/$(EXECUTABLE_NAME) ~/.packer.d/plugins
	@echo "Published plugin '$(EXECUTABLE_NAME)' to ~/.packer.d/plugins."

_build: _build-windows64 _build-windows32 _build-linux64 _build-mac64

_build-windows64:
	GOOS=windows GOARCH=amd64 go build -o $(BIN_DIRECTORY)/windows-amd64/$(EXECUTABLE_NAME).exe $(PLUGIN_FOLDER)

_build-windows32:
	GOOS=windows GOARCH=386 go build -o $(BIN_DIRECTORY)/windows-386/$(EXECUTABLE_NAME).exe $(PLUGIN_FOLDER)

_build-linux64:
	GOOS=linux GOARCH=amd64 go build -o $(BIN_DIRECTORY)/linux-amd64/$(EXECUTABLE_NAME) $(PLUGIN_FOLDER)

_build-mac64:
	GOOS=darwin GOARCH=amd64 go build -o $(BIN_DIRECTORY)/darwin-amd64/$(EXECUTABLE_NAME) $(PLUGIN_FOLDER)

_dist:
	@echo "Building distribution package for '$(EXECUTABLE_NAME)'..."
	cd $(BIN_DIRECTORY)/windows-386 && \
		zip -9 ../$(DIST_ZIP_PREFIX).windows-386.zip $(EXECUTABLE_NAME).exe
	cd $(BIN_DIRECTORY)/windows-amd64 && \
		zip -9 ../$(DIST_ZIP_PREFIX).windows-amd64.zip $(EXECUTABLE_NAME).exe
	cd $(BIN_DIRECTORY)/linux-amd64 && \
		zip -9 ../$(DIST_ZIP_PREFIX).linux-amd64.zip $(EXECUTABLE_NAME)
	cd $(BIN_DIRECTORY)/darwin-amd64 && \
		zip -9 ../$(DIST_ZIP_PREFIX)-darwin-amd64.zip $(EXECUTABLE_NAME)
