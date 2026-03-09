@echo off
echo Sending Kill signal to File Graph Visualizer...
powershell -NoProfile -Command "try { Invoke-WebRequest -Uri http://localhost:8080/api/kill -Method Get -TimeoutSec 2 > $null } catch { }"
echo.
echo If the server was running, it has been stopped.
echo Loop in RUN.bat should also be terminated.
pause
