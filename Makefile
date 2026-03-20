APP_NAME := BreakReminder
BINARY := breakreminder
CMD := ./cmd/breakreminder
APP_BUNDLE := $(APP_NAME).app
APP_ID := com.zhangxinyu.breakreminder

.PHONY: build run test clean package install uninstall deps

build:
	go build -o $(BINARY) $(CMD)

run: build
	./$(BINARY)

test:
	go test ./...

clean:
	rm -f $(BINARY)
	rm -rf $(APP_BUNDLE)

package: build
	@mkdir -p $(APP_BUNDLE)/Contents/MacOS
	@mkdir -p $(APP_BUNDLE)/Contents/Resources
	@cp $(BINARY) $(APP_BUNDLE)/Contents/MacOS/
	@cp bundle/Info.plist $(APP_BUNDLE)/Contents/
	@cp assets/icon.png $(APP_BUNDLE)/Contents/Resources/
	@codesign --force --sign - --timestamp=none --identifier $(APP_ID) $(APP_BUNDLE)/Contents/MacOS/$(BINARY)
	@codesign --force --deep --sign - --timestamp=none --identifier $(APP_ID) $(APP_BUNDLE)
	@codesign --verify --deep --strict $(APP_BUNDLE)
	@echo "Built $(APP_BUNDLE)"

install: package
	@echo "Installing to /Applications..."
	@rm -rf /Applications/$(APP_BUNDLE)
	@cp -R $(APP_BUNDLE) /Applications/
	@echo "Installed. You can launch from /Applications or Spotlight."
	@open /Applications/$(APP_BUNDLE)

uninstall:
	@rm -rf /Applications/$(APP_BUNDLE)
	@echo "Uninstalled $(APP_NAME).app"

deps:
	go mod tidy
