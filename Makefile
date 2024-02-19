all: clean
	go build -o ./bin/arc ./cmd/main.go

clean: 
	rm -rf ./bin/*