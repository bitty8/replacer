GC=go

BUILD_DIR=build

all: clean build_dir rpc

.PHONY: rpc
rpc: build_dir
	$(GC) build -o $(BUILD_DIR)/rpc cmd/rpc/rpc.go

.PHONY: build_dir
build_dir:
	@mkdir $(BUILD_DIR)

.PHONY: clean
clean:
	@rm -rf $(BUILD_DIR)