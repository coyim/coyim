mkdir %GOPATH%\src\github.com\twstrike\
xcopy %APPVEYOR_BUILD_FOLDER%\* %GOPATH%\src\github.com\twstrike\coyim /e /i /s /EXCLUDE:%MSYS_PATH% > nul
xcopy %APPVEYOR_BUILD_FOLDER%\.git %GOPATH%\src\github.com\twstrike\coyim\ /e /i /s
dir %GOPATH%\src\github.com\twstrike\coyim
if "%METHOD%"=="ci" SET MSYS_PATH=c:\msys64
if "%METHOD%"=="cross" SET MSYS_PATH=%APPVEYOR_BUILD_FOLDER%\msys%MSYS2_BITS%
SET PATH=%MSYS_PATH%\usr\bin;%PATH%
SET PATH=%MSYS_PATH%\mingw%MSYS2_BITS%\bin;%PATH%
if "%METHOD%"=="cross" appveyor DownloadFile http://kent.dl.sourceforge.net/project/msys2/Base/%MSYS2_ARCH%/msys2-base-%MSYS2_ARCH%-%MSYS2_BASEVER%.tar.xz
if "%METHOD%"=="cross" 7z x msys2-base-%MSYS2_ARCH%-%MSYS2_BASEVER%.tar.xz > nul
if "%METHOD%"=="cross" 7z x msys2-base-%MSYS2_ARCH%-%MSYS2_BASEVER%.tar > nul
%MSYS_PATH%\usr\bin\bash -lc "echo update-core starting..." 2> nul
%MSYS_PATH%\usr\bin\bash -lc "update-core" > nul
%MSYS_PATH%\usr\bin\bash -lc "echo install-deps starting..."
%MSYS_PATH%\usr\bin\bash -lc "pacman --noconfirm --needed -Sy autoconf" > nul
%MSYS_PATH%\usr\bin\bash -lc "pacman --noconfirm --needed -Sy automake" > nul
%MSYS_PATH%\usr\bin\bash -lc "pacman --noconfirm --needed -Sy make" > nul
%MSYS_PATH%\usr\bin\bash -lc "pacman --noconfirm --needed -Sy mingw-w64-%MSYS2_ARCH%-libiconv" > nul
%MSYS_PATH%\usr\bin\bash -lc "pacman --noconfirm --needed -Sy mingw-w64-%MSYS2_ARCH%-gcc" > nul
%MSYS_PATH%\usr\bin\bash -lc "pacman --noconfirm --needed -Sy mingw-w64-%MSYS2_ARCH%-gdb" > nul
%MSYS_PATH%\usr\bin\bash -lc "pacman --noconfirm --needed -Sy mingw-w64-%MSYS2_ARCH%-make" > nul
%MSYS_PATH%\usr\bin\bash -lc "pacman --noconfirm --needed -Sy zlib-devel" > nul
%MSYS_PATH%\usr\bin\bash -lc "pacman --noconfirm --needed -Sy mingw-w64-%MSYS2_ARCH%-pango" > nul
%MSYS_PATH%\usr\bin\bash -lc "pacman --noconfirm --needed -Sy mingw-w64-%MSYS2_ARCH%-gtk3" > nul
%MSYS_PATH%\usr\bin\bash -lc "pacman --noconfirm --needed -Sy mingw-w64-%MSYS2_ARCH%-pkg-config" > nul
%MSYS_PATH%\usr\bin\bash -lc "yes|pacman --noconfirm -Sc" > nul
if "%METHOD%"=="cross" %MSYS_PATH%\autorebase.bat > nul
