all: run

bin:
	mkdir -p ./bin
	go build -o ./bin .

templ:
	templ generate

run: bin templ
	./bin/$(shell basename $(CURDIR))

clean:
	rm -rf ./bin
