BINARY_NAME := fluss
DIST_DIR := dist

clean:
	rm -rf $(DIST_DIR)

build: clean
	mkdir -p $(DIST_DIR)
	go build -o $(DIST_DIR)/$(BINARY_NAME)
