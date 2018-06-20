##
## Avalanche build script
##

##### CONFIG BUILD
DEBUG=1

##### TARGET FILES
SERVER_ENTRY=app/server.go
CLI_ENTRY=app/commands/cli.go
LOGGER_MODULES=$(wildcard ./app/modules/core/logger/*/*.go)
CLI_MODULES=$(wildcard ./app/commands/*/*.go)
APP_MODULES=$(wildcard ./app/modules/*/main/*.go)

##### DEPENDENCIES
DEPENDENCIES=\
 github.com/joho/godotenv \
 github.com/hjson/hjson-go \
 github.com/rivo/tview \
 github.com/jinzhu/gorm \
 github.com/sirupsen/logrus \
 github.com/uniplaces/carbon \
 github.com/nicksnyder/go-i18n/goi18n \
 github.com/CloudyKit/jet \
 github.com/peyman-abdi/fasthttp-routing \
 github.com/peyman-abdi/avest \
 github.com/json-iterator/go \
 github.com/graphql-go/graphql \
 github.com/stretchr/testify/assert \


##### BUILD COMMANDS
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOTOOL=$(GOCMD) tool
GOFMT=$(GOCMD) fmt

#### EXPORT PATHES
BINARY_PATH=bin/platforms/$(PLATFORM)/$(VARIANT)
MODULES_PATH=bin/platforms/$(PLATFORM)/$(VARIANT)/modules

#### BUILD VARIABLES
VERSION=1.0.0-Alpha1
PLATFORM=$(call getPlatform)
VARIANT=$(call getVariant)
BUILD_TIME=$(date ”%Y.%m.%d.%H%M%S”)
BUILD_CODE=$(shell git rev-parse HEAD)
VARS_PACKAGE=github.com/peyman-abdi/avalanche/app/modules/core/app

LDFLAGS=-ldflags "-X $(VARS_PACKAGE).Version=$(VERSION) -X $(VARS_PACKAGE).Code=$(BUILD_CODE) -X $(VARS_PACKAGE).Variant=$(VARIANT) -X $(VARS_PACKAGE).Platform=$(PLATFORM) -X $(VARS_PACKAGE).BuildTime=$(BUILD_TIME)"

#### SCRIPTS
all: build test
build_n_serve: build serve
build: modules build_cli_only build_server_only
server_n_serve: build_server_only serve
build_cli: build_cli_only cli_modules app_modules logger_modules
build_server: build_server_only logger_modules app_modules
modules: logger_modules cli_modules app_modules

build_server_only:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_PATH)/$(basename $(notdir $(SERVER_ENTRY))) -v $(SERVER_ENTRY)
build_cli_only:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_PATH)/$(basename $(notdir $(CLI_ENTRY))) -v $(CLI_ENTRY)
test:
	$(GOTEST) $(LDFLAGS) -v -covermode=count -coverprofile=coverage.out ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_PATH)
run:
	$(GORUN) $(ENTRY)
serve:
	./$(BINARY_PATH)/$(basename $(notdir $(SERVER_ENTRY)))
logger_modules:
	$(foreach file, $(LOGGER_MODULES), $(GOBUILD) -v -buildmode=plugin -o $(call getLoggerModulesPath, $(file)) $(file);)
cli_modules:
	$(foreach file, $(CLI_MODULES), $(GOBUILD) -v -buildmode=plugin -o $(call getCliModulesPath, $(file)) $(file);)
app_modules:
	$(foreach file, $(APP_MODULES), $(GOBUILD) -v -buildmode=plugin -o $(call getAppModulesPath, $(file)) $(file);)
go_coverage:
	$(GOTOOL) cover -html=coverage.out
go_format:
	$(GOFMT) ./app/...
go_get:
	@($(foreach dep, $(DEPENDENCIES), $(GOGET) $(dep);))


##### METHODS
define getAppModulesPath
$(MODULES_PATH)/app/$(basename $(notdir $(1))).so
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