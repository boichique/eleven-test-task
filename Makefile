run-service:
	go build -o ./client/bin/cli ./client/. &&\
	docker compose up -d --build

stop-service:
	rm -r ./client/bin &&\
	docker compose down

before-push:
	go mod tidy &&\
	gofumpt -l -w . &&\
	go build ./...&&\
	golangci-lint run ./... &&\
	go test -v ./tests/...