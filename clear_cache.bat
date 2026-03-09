@echo off
echo.
echo ** Privacy and Cache Cleaner **
echo This will stop the server and delete all log files.
echo.

echo Stopping File Graph Visualizer to release log files...
rem Use the same logic as stop.bat to ensure the server and loop are killed.
powershell -NoProfile -Command "try { Invoke-WebRequest -Uri http://localhost:8080/api/kill -Method Get -TimeoutSec 1 > $null } catch { }" > nul 2>&1
echo Server stopped.
echo.

echo Deleting log files from the 'logs' directory...
if exist "logs\*.log" (
    del /Q "logs\*.log"
    echo Log files deleted successfully.
) else (
    echo No log files found to delete.
)

echo.
echo ---------------------------------
echo  Privacy cleaning complete.
echo ---------------------------------
echo.
pause
