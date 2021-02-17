#/usr/bin/make
build:
	go build -o ${GOPATH}/bin/goserve -ldflags "-s -w" main.go
