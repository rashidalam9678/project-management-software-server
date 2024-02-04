build:
	go build -o bin/app cmd/web/main.go

run:
	go run cmd/web/*.go


.PHONEY: build run


