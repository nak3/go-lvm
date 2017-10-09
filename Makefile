all:
	go build -o example cmd/example.go
	go build -o testrun cmd/test.go
run:
	sudo ./example
test:
	sudo ./testrun
clean:
	rm example

