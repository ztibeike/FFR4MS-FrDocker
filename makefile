INSTALL_PATH = /usr/local/bin
BUILD_PATH = build
APP_NAME = frdocker

.PHONY: all
all: clean build

.PHONY: build
build:
	@echo "Building..."
	@go mod tidy
	@go build -o $(BUILD_PATH)/$(APP_NAME) main.go

.PHONY: install
install:
	@echo "Installing..."
	@cp $(BUILD_PATH)/$(APP_NAME) $(INSTALL_PATH)/$(APP_NAME)

.PHONY: uninstall
uninstall:
	@echo "Uninstalling..."
	@rm -f $(INSTALL_PATH)/$(APP_NAME)

.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_PATH)