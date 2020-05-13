build:
	# Linux 32-bit
	GOOS=linux GOARCH=386 go build -o bin/detect-gps src/detect-gps.go

clean:
	@rm -rf bin/

all: build