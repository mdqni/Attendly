GO      ?= go
GOOS    ?= linux
GOARCH  ?= amd64
TAGS    ?= netgo
OUTPUT_DIR := bin

SERVICES := auth user group qr api-gateway

build: $(SERVICES)

$(OUTPUT_DIR):
	mkdir -p $(OUTPUT_DIR)

$(SERVICES): %: $(OUTPUT_DIR)
	$(GO) build -tags '$(TAGS)' -ldflags="-s -w" \
		-o $(OUTPUT_DIR)/$@ ./services/$@/cmd/main.go
	@echo "Built $@ → $(OUTPUT_DIR)/$@"

clean:
	rm -rf $(OUTPUT_DIR)
	@echo "Cleaned"

.PHONY: build clean $(SERVICES)
