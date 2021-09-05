set GOARCH=wasm
set GOOS=js
go build -o main.wasm
set GOARCH=amd64
set GOOS=windows

copy main.wasm ..\simulator2\dist
pause