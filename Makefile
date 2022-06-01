# Build and test automation

APP_NAME=pacman

build:
	go build -o ${APP_NAME}.o ${APP_NAME}.go
test: build
	@echo Test 1
	./${APP_NAME}.o  || true
run:
	go run ${APP_NAME}.go
clean:
	rm -rf *.o