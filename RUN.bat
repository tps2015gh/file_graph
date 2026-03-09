@echo off
:loop
cls
echo Starting File Graph Visualizer...
go run main.go
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
