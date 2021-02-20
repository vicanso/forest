.PHONY: default test test-cover dev generate hooks lint-web

# for dev
dev:
	air -c .air.toml	
doc:
	swagger generate spec -o ./api.yml && swagger validate ./api.yml 

test:
	go test -race -cover ./...

generate: 
	go generate ./ent

describe:
	entc describe ./ent/schema

test-cover:
	go test -race -coverprofile=test.out ./... && go tool cover --html=test.out

list-mod:
	go list -m -u all

tidy:
	go mod tidy

build:
	go build -ldflags "-X main.Version=0.0.1 -X 'main.BuildedAt=`date`'" -o forest 


lint:
	golangci-lint run --timeout 2m --skip-dirs web

lint-web:
	cd web && yarn lint 

hooks:
	cp hooks/* .git/hooks/