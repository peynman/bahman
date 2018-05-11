##
## Avalanche build script
##

##### CONFIG BUILD
DEBUG=1

##### METHODS
define getAppModulesPath
$(MODULES_PATH)/plugins/$(basename $(notdir $(1))).so
endef
define getCliModulesPath
$(MODULES_PATH)/console/$(basename $(notdir $(1))).so
endef
define getLoggerModulesPath
$(MODULES_PATH)/channels/$(basename $(notdir $(1))).so
endef
define getVariant
$(if ($(DEBUG),1),DEBUG,PRODUCTION)
endef
define getPlatform
$(if ($(OS),Windows_NT),$(if ($(shell uname -s),Linux),OSX,LINUX),WIN32)
endef

##### BUILD VARIABLES
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
SERVER_ENTRY=app/server.go
CLI_ENTRY=app/commands/cli.go
LOGGER_MODULES=$(wildcard ./app/core/logger/*/*.go)
CLI_MODULES=$(wildcard ./app/commands/*/*.go)
APP_MODULES=$(wildcard ./app/modules/*/main/*.go)

#### EXPORT PATHES
BINARY_PATH=bin/platforms/$(PLATFORM)/$(VARIANT)
MODULES_PATH=bin/platforms/$(PLATFORM)/$(VARIANT)/modules

#### GO VARIABLES
VERSION=1.0.0-Alpha1
PLATFORM=$(call getPlatform)
VARIANT=$(call getVariant)
BUILD_TIME=$(date ”%Y.%m.%d.%H%M%S”)
BUILD_CODE=$(shell git rev-parse HEAD)
PACKAGE=avalanche/app/core/app

LDFLAGS=-ldflags "-X $(PACKAGE).Version=$(VERSION) -X $(PACKAGE).Code=$(BUILD_CODE) -X $(PACKAGE).Variant=$(VARIANT) -X $(PACKAGE).Platform=$(PLATFORM) -X $(PACKAGE).BuildTime=$(BUILD_TIME)"

#### SCRIPTS
all: test modules build
build_n_serve: modules build serve
build: build_cli build_server
build_cli: build_cli_only cli_modules app_modules
build_server: build_server_only logger_modules app_modules
modules: logger_modules cli_modules app_modules

build_server_only:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_PATH$)/$(basename $(notdir $(SERVER_ENTRY))) -v $(SERVER_ENTRY)
build_cli_only:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_PATH$)/$(basename $(notdir $(CLI_ENTRY))) -v $(CLI_ENTRY)
test:
	$(GOTEST) ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_PATH)
run:
	$(GORUN) $(ENTRY)
serve:
	./$(BINARY_PATH)/$(basename $(notdir $(SERVER_ENTRY)))
logger_modules:
	$(foreach file, $(LOGGER_MODULES), $(GOBUILD) -buildmode=plugin -o $(call getLoggerModulesPath, $(file)) -v $(file);)
cli_modules:
	$(foreach file, $(CLI_MODULES), $(GOBUILD) -buildmode=plugin -o $(call getCliModulesPath, $(file)) -v $(file);)
app_modules:
	$(foreach file, $(APP_MODULES), $(GOBUILD) -buildmode=plugin -o $(call getAppModulesPath, $(file)) -v $(file);)
