.PHONY: win64
win64:
	GOOS=windows GOARCH=amd64 go build -o lkr-meldung-export.exe -ldflags '-w -s' .