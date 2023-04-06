
build:
	go build -o bin/ddosarmy main.go

run:
	./bin/ddosarmy

test :
	go test -v ./...
