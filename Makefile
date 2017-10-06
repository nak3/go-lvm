all:
	go build -o example cmd/example.go
run:
	sudo ./example
clean:
	rm example

