all: run

templgen:
	templ generate

gobuild:
	mkdir -p ./bin
	go build -o ./bin .

bin: templgen gobuild

run: bin 
	./bin/$(shell basename $(CURDIR))

airbuild: templgen
	mkdir -p ./tmp
	go build -o ./tmp/main .

air: 
	air

clean:
	rm -rf ./bin
	rm -rf ./tmp
