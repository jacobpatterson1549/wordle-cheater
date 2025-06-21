.PHONY: all test coverage clean run serve

BUILD_DIR := build
BIN_DIR := $(BUILD_DIR)/bin
COVERAGE_OBJ := coverage.out
WORDS_OBJ := words.txt
SRC := *.go
GO_SRC_FN = find $(1) $(foreach g,$(GENERATE_SRC),-path $g -prune -o) -print 
SRC := $(shell $(call GO_SRC_FN,cmd/ internal/ *.go))
OBJ_DIRS := $(wildcard cmd/*/)
OBJS := $(patsubst cmd/%/,%,$(OBJ_DIRS))
OBJ_BINS := $(addprefix $(BIN_DIR)/,$(OBJS))
GO_ARGS :=

all: $(OBJ_BINS)

test: $(BUILD_DIR)/$(COVERAGE_OBJ)

coverage: $(BUILD_DIR)/$(COVERAGE_OBJ)
	go tool cover -html=$<

clean:
	rm -rf $(BUILD_DIR)

$(BUILD_DIR) $(BIN_DIR):
	mkdir -p $@

$(BIN_DIR)/%: $(BUILD_DIR)/$(COVERAGE_OBJ)
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

