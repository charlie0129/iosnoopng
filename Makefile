build:
	npm run build
	CGO_ENABLED=0 go build -o iosnoopng .

run: build
	sudo ./iosnoopng
