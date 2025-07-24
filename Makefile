.PHONY: run convert

run:
	go run main.go

convert:
	go run ./tools/converter.go