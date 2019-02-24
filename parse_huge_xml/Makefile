SHELL = /bin/bash -o pipefail

APP_NAME = parse_huge_xml

ROOT_DIR = $(CURDIR)
BIN_DIR = $(ROOT_DIR)/bin
APP = $(BIN_DIR)/$(APP_NAME)

LOG_DIR = $(ROOT_DIR)/log
LOG_FILE = $(LOG_DIR)/$(APP_NAME).log
LOG_FILE_BUILD = $(LOG_DIR)/build-$(APP_NAME).log

go-build:
	cd $(ROOT_DIR) && go mod vendor -v 2>&1 | tee --append $(LOG_FILE_BUILD)
	cd $(ROOT_DIR) && go build -o $(APP) -v -x -mod vendor 2>&1 | tee --append $(LOG_FILE_BUILD)

go-run:
	cd $(ROOT_DIR) && go mod vendor -v 2>&1 | tee --append $(LOG_FILE_BUILD)
	cd $(ROOT_DIR) && go run -v -mod vendor main.go 2>&1 | tee --append $(LOG_FILE_BUILD)
