dep:
	go mod tidy
	cd ./tests_playwright && npm i

dep-browsers:
	cd ./tests_playwright && npx playwright install

test:
	go test ./... -v

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

test-e2e:
	cd ./tests_playwright && npx playwright test

test-e2e-report:
	cd ./tests_playwright && npx playwright show-report

test-e2e-ui:
	cd ./tests_playwright &&npx playwright test --ui

tailwind:
	tailwindcss -i ./server/public/src/tailwind.css -o ./server/public/dist/css/tailwind.css

tailwind-watch:
	tailwindcss -i ./server/public/src/tailwind.css -o ./server/public/dist/css/tailwind.css --watch

dev:
	go run . -dev

build:
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin ${TARGET} ${ARGS}
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux ${TARGET} ${ARGS}
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-windows ${TARGET} ${ARGS}

clean:
	go clean
	rm ${BINARY_NAME}-darwin
	rm ${BINARY_NAME}-linux
	rm ${BINARY_NAME}-windows