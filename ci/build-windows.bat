if "%METHOD%"=="ci" SET MSYS_PATH=c:\msys64
if "%METHOD%"=="cross" SET MSYS_PATH=%APPVEYOR_BUILD_FOLDER%\msys%MSYS2_BITS%
%MSYS_PATH%\usr\bin\bash -lc "cd /c/gopath/src/github.com/coyim/coyim && make win-ci-deps"
%MSYS_PATH%\usr\bin\bash -lc "cd /c/gopath/src/github.com/coyim/coyim && make build-gui-win"
