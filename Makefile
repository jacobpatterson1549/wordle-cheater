.PHONY: all test coverage clean run serve

OBJ := wordle-cheater
SERVER_OBJ := server
BUILD_DIR := build
COVERAGE_OBJ := coverage.out
WORDS_OBJ := words.txt
SRC := *.go
GO_SRC_FN = find $(1) $(foreach g,$(GENERATE_SRC),-path $g -prune -o) -print 
SRC := $(shell $(call GO_SRC_FN,cmd/ internal/ *.go))
GO_ARGS :=

all: $(BUILD_DIR)/$(OBJ) $(BUILD_DIR)/$(SERVER_OBJ)

test: $(BUILD_DIR)/$(COVERAGE_OBJ)

coverage: $(BUILD_DIR)/$(COVERAGE_OBJ)
	go tool cover -html=$<

clean:
	rm -rf $(BUILD_DIR)

run: $(BUILD_DIR)/$(OBJ)
	$<

serve: $(BUILD_DIR)/$(SERVER_OBJ)
	$<

$(BUILD_DIR):
	mkdir -p $@

$(BUILD_DIR)/$(SERVER_OBJ): $(BUILD_DIR)/$(COVERAGE_OBJ) | $(BUILD_DIR)
	go list ./... | grep -E cmd/$(@F)$$ \
		| $(GO_ARGS) xargs go build \
			-o $@

$(BUILD_DIR)/$(OBJ): $(BUILD_DIR)/$(COVERAGE_OBJ) | $(BUILD_DIR)
	go list ./... | grep -E cmd/$(@F)$$ \
		| $(GO_ARGS) xargs go build \
			-o $@

$(BUILD_DIR)/$(COVERAGE_OBJ): $(SRC) $(BUILD_DIR)/$(WORDS_OBJ) | $(BUILD_DIR)
	go test ./... -coverprofile=$@

$(BUILD_DIR)/$(WORDS_OBJ): | $(BUILD_DIR)
	aspell -d en_US dump master \
		| sort \
		| uniq \
		| grep -E ^[a-z]{5}$$ \
		> $@

