NODE=node
NPM=npm
BROWSERIFY=node_modules/browserify/bin/cmd.js
BUILD=build
SRC=src/*
TARGET=kaiju-client.js

build_client: npm
	$(NODE) $(BROWSERIFY) -o $(BUILD)/$(TARGET) $(SRC)

npm:
	$(NPM) install
