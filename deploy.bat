@echo off
setlocal

set "DEST_DIR=%~1"

if "%DEST_DIR%"=="" (
    echo.
    echo Usage: deploy.bat "C:\Path\To\Destination\Folder"
    echo.
    echo This will copy the MOST RECENTly built .exe and index.html to the destination.
    exit /b 1
)

if not exist "%DEST_DIR%" (
    echo Creating destination directory: %DEST_DIR%
    mkdir "%DEST_DIR%"
)

rem Find the most recent built .exe
for /f "delims=" %%i in ('dir /b /a-d /o-d file_graph_*.exe') do (
    set "LATEST_EXE=%%i"
    goto :found_exe
)

:found_exe
if "%LATEST_EXE%"=="" (
    echo No built .exe found starting with 'file_graph_'. Please run BUILD.bat first.
    exit /b 1
)

echo.
echo Deploying to: %DEST_DIR%
echo -------------------------------------------
echo Copying Executable: %LATEST_EXE%
copy /y "%LATEST_EXE%" "%DEST_DIR%\file_graph_app.exe"

echo Copying UI File: index.html (MUST BE PRESENT)
copy /y "index.html" "%DEST_DIR%\index.html"
echo -------------------------------------------
echo Done! Run it from %DEST_DIR%
echo.

endlocal
pause
