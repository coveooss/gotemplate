install:
	go install

deploy:
	GOOS=linux go build -o .pkg/gotemplate_linux
	GOOS=darwin go build -o .pkg/gotemplate_darwin
	GOOS=windows go build -o .pkg/gotemplate_windows.exe
