##
# Boring Calender Server
#
# @file
# @version 0.1
EXEC      = boring-calendar-cli
BUILD_DIR = ./build
CC        = go
BLD       = build
SRC       = ./main.go

all: init-db build run

build:
	@mkdir -p $(BUILD_DIR)
	$(CC) $(BLD) -o $(BUILD_DIR)/$(EXEC) $(SRC)

run:
	@$(BUILD_DIR)/$(EXEC)

.PHONY: build

# end
