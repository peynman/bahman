##
## Avalanche build script
##

##### CONFIG BUILD
BINARY_NAME=avalanche_server
MODULE_FILES=$(wildcard ./app/libs/logger/*/*.go)
DEBUG=1

##### METHODS
define getModulePath
$(MODULES_PATH)/$(if ($(findstring $(1),logger),''),channels,$(lastword $(subst /, ,$(dir $(1)))))/$(basename $(notdir $(1))).so
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
ENTRY=app/main.go

#### EXPORT PATHES
BINARY_PATH=bin/platforms/$(PLATFORM)/$(VARIANT)/$(BINARY_NAME)
MODULES_PATH=bin/platforms/$(PLATFORM)/$(VARIANT)/modules

#### GO VARIABLES
VERSION=1.0.0-Alpha1
PLATFORM=$(call getPlatform)
VARIANT=$(call getVariant)
BUILD_TIME=$(date ”%Y.%m.%d.%H%M%S”)
BUILD_CODE=$(shell git rev-parse HEAD)
PACKAGE=avalanche/app/libs

LDFLAGS=-ldflags "-X $(PACKAGE).Version=$(VERSION) -X $(PACKAGE).Code=$(BUILD_CODE) -X $(PACKAGE).Variant=$(VARIANT) -X $(PACKAGE).Platform=$(PLATFORM) -X $(PACKAGE).BuildTime=$(BUILD_TIME)"

#### SCRIPTS
all: test modules build
build_n_run: modules build dev
build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_PATH) -v $(ENTRY)
test:
	$(GOTEST) ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_PATH)
run:
	$(GORUN) $(ENTRY)
dev:
	./$(BINARY_PATH)
modules:
	$(foreach file, $(MODULE_FILES), $(GOBUILD) -buildmode=plugin -o $(call getModulePath, $(file)) -v $(file);)
