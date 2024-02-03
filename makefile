BINARY_NAME=build
TARGET=./
ARGS=

dep:
	go mod tidy
	cd ./tests_playwright && npm i
	cd ./server && npm i

dep-browsers:
	cd ./tests_playwright && npx playwright install

test:
	go test ./...

test-v:
	go test ./... -v

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

test-e2e:
	cd ./tests_playwright && npx playwright test

test-e2e-report:
	cd ./tests_playwright && npx playwright show-report

test-e2e-ui:
	cd ./tests_playwright && npx playwright test --ui

assets-clear:
	cd ./server && npm run cache:clear

assets-build:
	cd ./server && npm run cpx:build

assets-watch:
	cd ./server && npm run cpx:watch

parcel-build:
	cd ./server && npm run parcel:build

parcel-watch:
	cd ./server && npm run parcel:watch

dev:
	go run . -dev

build:
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-windows ${TARGET} ${ARGS}
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin ${TARGET} ${ARGS}
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux ${TARGET} ${ARGS}

clean:
	go clean
	rm ${BINARY_NAME}-windows
	rm ${BINARY_NAME}-darwin
	rm ${BINARY_NAME}-linux