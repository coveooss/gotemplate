install:
	go install

deploy:
	GOOS=linux go build -o .pkg/goremote_linux
	GOOS=darwin go build -o .pkg/goremote_darwin
	GOOS=windows go build -o .pkg/goremote_windows.exe
