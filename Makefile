.PHONY: deps
deps:
	go mod download
	go mod tidy

.PHONY: build
build:
	go build -o bin/api ./cmd/api

.PHONY: build-cache
build-cache:
	@echo "Using Go build cache for faster builds"
	go build -i -o bin/api ./cmd/api

.PHONY: test
test:
	go test -v ./... -race -timeout 300s -cover

.PHONY: swagger
swagger:
	swag init -g cmd/api/main.go -o docs

lint:
	golangci-lint run

clean:
	go clean -testcache
	go clean -modcache
	go clean -cache

.PHONY: migrate-up
migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" -verbose up

.PHONY: migrate-down
migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" -verbose down

.PHONY: migrate-create
migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

.PHONY: migrate-force
migrate-force:
	migrate -path migrations -database "$(DATABASE_URL)" force $(version)

.PHONY: migrate-version
migrate-version:
	migrate -path migrations -database "$(DATABASE_URL)" version

aws-check:
	aws sts get-caller-identity

cert:
	openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout mykey.key -out mycert.crt \
  -subj "/CN=express-nodejs-demo-alb-1541480660.us-east-1.elb.amazonaws.com"

	aws iam upload-server-certificate \
    --server-certificate-name express-nodejs-demo-cert \
    --certificate-body file://mycert.crt \
    --private-key file://mykey.key

up:
	docker compose -f Docker-compose.yml up -d --build

down:
	docker compose down

restart:
	make stop
	make run