@echo off
setlocal

set MYDIR=%~dp0
set MYDIR=%MYDIR:~0,-1%

:loop
cls
echo Starting File Graph Visualizer...
"%MYDIR%\file_graph_server.exe" -startpath="%MYDIR%"
if %ERRORLEVEL% NEQ 0 (
    echo.
    echo Server stopped by user request (Kill).
    echo Close this window or press any key to exit.
    pause > nul
    exit /b
)
echo Server exited (Restart). Restarting in 2 seconds...
timeout /t 2 > nul
goto loop
