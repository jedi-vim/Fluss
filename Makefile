BINARY_NAME := fluss
DIST_DIR := dist

build:
	mkdir -p $(DIST_DIR)
	go build -o $(DIST_DIR)/$(BINARY_NAME)

clean:
	rm -rf $(DIST_DIR)

