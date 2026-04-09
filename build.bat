@echo off
setlocal

pushd web
echo Current directory: %cd%
echo Building web...

call npm run build

popd
pushd teamweb
echo Building teamweb...
call npm run build
popd

echo Syncing teamweb dist to web\\team...
if exist web\\team\\assets rmdir /s /q web\\team\\assets
copy /Y teamweb\\dist\\index.html web\\team\\data.html >nul
copy /Y teamweb\\dist\\favicon.ico web\\team\\favicon.ico >nul
xcopy /Y /E /I teamweb\\dist\\assets web\\team\\assets >nul

echo Current directory: %cd%
echo Building Go executable...

mkdir dist 2>nul

echo Target: windows amd64
set GOOS=windows
set GOARCH=amd64
go build -tags="nomsgpack" -ldflags="-s -w" -o dist\stzbHelper-windows-amd64.exe stzbHelper

echo Done. Output: dist\\stzbHelper-windows-amd64.exe
pause
