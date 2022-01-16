.PHONY: all test coverage clean run

OBJ := wordle-cheater
BUILD_DIR := build
COVERAGE_OBJ := coverage.out
WORDS_OBJ := words.txt
SRC := $(shell find -name '*.go')

all: $(BUILD_DIR)/$(OBJ)

test: $(BUILD_DIR)/$(COVERAGE_OBJ)

coverage: $(BUILD_DIR)/$(COVERAGE_OBJ)
	go tool cover -html=$<

clean:
	rm -rf $(BUILD_DIR) $(WORDS_OBJ)

run: $(BUILD_DIR)/$(OBJ)
	$<

$(BUILD_DIR):
	mkdir -p $@

$(BUILD_DIR)/$(OBJ): $(BUILD_DIR)/$(COVERAGE_OBJ) | $(BUILD_DIR)
	go build -o $@

$(BUILD_DIR)/$(COVERAGE_OBJ): $(SRC) $(WORDS_OBJ) | $(BUILD_DIR)
	go test ./... -coverprofile=$@

$(WORDS_OBJ): 
	aspell -d en_US dump master \
		| sort \
		| uniq \
		| grep -E ^[a-z]{5}$$ \
		> $@
	truncate -s -1 $@
