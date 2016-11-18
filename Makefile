all: build
get-deps:
	go get -u github.com/Masterminds/glide
	glide install
build:
	CGO_ENABLED=0 go build -o ${BUILD_PREFIX}/nspawnlet ./main.go
test:
	go test ./src/...
clean:
	rm -rf ./**/*.o ./**/*.a ./**/*.so
uninstall:
	rm ${INSTALL_PREFIX}/nspawnlet
install:
	cp ${BUILD_PREFIX}/nspawnlet ${INSTALL_PREFIX}/nspawnlet
