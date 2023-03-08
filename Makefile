# ==============================================================================
# 定義全局 Makefile 變量

COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
# 項目根目錄
ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/ && pwd -P))
# 構建產物、臨時文件存放資料夾
OUTPUT_DIR := $(ROOT_DIR)/_output
# Protobuf 文件存放資料夾
APIROOT=$(ROOT_DIR)/pkg/proto

# ==============================================================================
# version variable

## version package for application used， `-ldflags -X` cmd can define valiable's value
VERSION_PACKAGE=microblog/pkg/version

## version NO.
ifeq ($(origin VERSION), undefined)
VERSION := $(shell git describe --tags --always --match='v*')
endif

## checak repo state, default:dirty
GIT_TREE_STATE:="dirty"
ifeq (, $(shell git status --porcelain 2>/dev/null))
	GIT_TREE_STATE="clean"
endif
GIT_COMMIT:=$(shell git rev-parse HEAD)

GO_LDFLAGS += \
	-X $(VERSION_PACKAGE).GitVersion=$(VERSION) \
	-X $(VERSION_PACKAGE).GitCommit=$(GIT_COMMIT) \
	-X $(VERSION_PACKAGE).GitTreeState=$(GIT_TREE_STATE) \
	-X $(VERSION_PACKAGE).BuildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

# ==============================================================================
# 定義 Makefile all 偽目標，執行 `make` 時，會默認執行 all 偽目標
.PHONY: all
all: add-copyright format build

# ==============================================================================
# 定義其他需要的偽目標

.PHONY: build
build: tidy # 編譯源碼，依賴 tidy 目標自動添加/移除依賴包.
	@go build -v -ldflags "$(GO_LDFLAGS)" -o $(OUTPUT_DIR)/microblog/microblog.exe $(ROOT_DIR)/cmd/microblog/main.go

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

.PHONY: ca
ca: ## 生成 CA 文件
	@mkdir -p $(OUTPUT_DIR)/cert
	@openssl genrsa -out $(OUTPUT_DIR)/cert/ca.key 1024 # 生成根證書私鑰
	@openssl req -new -key $(OUTPUT_DIR)/cert/ca.key -out $(OUTPUT_DIR)/cert/ca.csr \
		-subj "/C=CN/ST=Guangdong/L=Shenzhen/O=devops/OU=it/CN=127.0.0.1/emailAddress=@test@gmail.com" # 2. 生成請求文件
	@openssl x509 -req -in $(OUTPUT_DIR)/cert/ca.csr -signkey $(OUTPUT_DIR)/cert/ca.key -out $(OUTPUT_DIR)/cert/ca.crt # 3. 生成根證書
	@openssl genrsa -out $(OUTPUT_DIR)/cert/server.key 1024 # 4. 生成服務端私鑰
	@openssl rsa -in $(OUTPUT_DIR)/cert/server.key -pubout -out $(OUTPUT_DIR)/cert/server.pem # 5. 生成服務端公鑰
	@openssl req -new -key $(OUTPUT_DIR)/cert/server.key -out $(OUTPUT_DIR)/cert/server.csr \
		-subj "/C=CN/ST=Guangdong/L=Shenzhen/O=serverdevops/OU=serverit/CN=127.0.0.1/emailAddress=@test@gmail.com" # 6. 生成服務端向 CA 申請簽名的 CSR
	@openssl x509 -req -CA $(OUTPUT_DIR)/cert/ca.crt -CAkey $(OUTPUT_DIR)/cert/ca.key \
		-CAcreateserial -in $(OUTPUT_DIR)/cert/server.csr -out $(OUTPUT_DIR)/cert/server.crt # 7. 生成服務端帶有 CA 簽名的證書

protoc: ## 編譯 protobuf 文件
	@echo "===========> Generate protobuf files"
	@protoc                                            \
		--proto_path=$(APIROOT)                          \
		--proto_path=$(ROOT_DIR)/third_party             \
		--go_out=paths=source_relative:$(APIROOT)        \
		--go-grpc_out=paths=source_relative:$(APIROOT)   \
		$(shell find $(APIROOT) -name *.proto)