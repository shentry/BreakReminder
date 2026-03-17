APP_NAME := BreakReminder
BINARY := breakreminder
CMD := ./cmd/breakreminder

.PHONY: build run test clean package

build:
	go build -o $(BINARY) $(CMD)

run: build
	./$(BINARY)

test:
	go test ./...

clean:
	rm -f $(BINARY)
	rm -rf $(APP_NAME).app

package: build
	@mkdir -p $(APP_NAME).app/Contents/MacOS
	@mkdir -p $(APP_NAME).app/Contents/Resources
	@cp $(BINARY) $(APP_NAME).app/Contents/MacOS/
	@cp bundle/Info.plist $(APP_NAME).app/Contents/
	@cp assets/icon.png $(APP_NAME).app/Contents/Resources/
	@echo "Built $(APP_NAME).app"

deps:
	go mod tidy
