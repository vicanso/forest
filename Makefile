.PHONY: default test test-cover dev generate hooks lint-web doc

# for dev
dev:
	air -c .air.toml	
dev-debug:
	LOG_LEVEL=0 make dev
doc:
	CGO_ENABLED=0 swagger generate spec -t swagger -o ./asset/api.yml && swagger validate ./asset/api.yml 

test:
	go test -race -cover ./...

install:
	go get -d entgo.io/ent/cmd/entc@v0.11.8

generate: 
	rm -rf ./ent
	go run entgo.io/ent/cmd/ent generate --feature privacy --feature intercept ./schema --template ./template --target ./ent

describe:
	entc describe ./schema

test-cover:
	go test -race -coverprofile=test.out ./... && go tool cover --html=test.out

list-mod:
	go list -m -u all

tidy:
	go mod tidy

build:
	go build -ldflags "-X main.Version=0.0.1 -X 'main.BuildedAt=`date`'" -o forest 


lint:
	golangci-lint run

lint-web:
	cd web && yarn lint 

hooks:
	cp hooks/* .git/hooks/