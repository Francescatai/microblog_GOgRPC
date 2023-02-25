# ==============================================================================
# 定義全局 Makefile 變量

COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
# 項目根目錄
ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/ && pwd -P))
# 構建產物、臨時文件存放資料夾
OUTPUT_DIR := $(ROOT_DIR)/_output

# ==============================================================================
# 定義 Makefile all 偽目標，執行 `make` 時，會默認執行 all 偽目標
.PHONY: all
all: add-copyright format build

# ==============================================================================
# 定義其他需要的偽目標

.PHONY: build
build: tidy # 編譯源碼，依賴 tidy 目標自動添加/移除依賴包.
	@go build -v -o $(OUTPUT_DIR)/microblog $(ROOT_DIR)/cmd/microblog/main.go

.PHONY: format
format: # 格式化 Go 源碼
	@gofmt -s -w ./

.PHONY: add-copyright
add-copyright: # 添加版權資訊
	@addlicense -v -f $(ROOT_DIR)/scripts/boilerplate.txt $(ROOT_DIR) --skip-dirs=third_party,vendor,$(OUTPUT_DIR)

.PHONY: swagger
swagger: # 啟動swagger
	@swagger serve -F=swagger --no-open --port 65534 $(ROOT_DIR)/api/openapi/openapi.yaml

.PHONY: tidy
tidy: # 自動添加/移除依賴包
	@go mod tidy

.PHONY: clean
clean: # 清理構建產物、臨時文件
	@-rm -vrf $(OUTPUT_DIR)