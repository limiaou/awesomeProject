GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=sudoku

build:
	$(GOBUILD) -o $(BINARY_NAME)
test:
	cd ./src/api && $(GOTEST)
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
start:
	cowsay -f small "GO TEST..."
	cd ./src/api && $(GOTEST)
	cowsay -f small "GO TEST PASS!"
	cowsay -f small "GO BUILD..."
	$(GOBUILD) -o $(BINARY_NAME) -v
	cowsay -f small "GO BUILD SUCCESS!"
	cowsay -f small "GO RUN!"
	./$(BINARY_NAME)
install:
	dep ensure