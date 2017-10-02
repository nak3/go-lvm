all:
	go build -o example example.go liblvm.go
run:
	sudo ./example
clean:
	rm example

