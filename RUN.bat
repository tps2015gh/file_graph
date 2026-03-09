@echo off
:loop
cls
echo Starting File Graph Visualizer...
go run main.go
echo Server exited. Restarting in 2 seconds...
timeout /t 2 > nul
goto loop
