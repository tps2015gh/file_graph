@echo off
echo Building...

for /f "tokens=2 delims==" %%a in ('wmic OS Get localdatetime /value') do set "dt=%%a"
set "YY=%dt:~2,2%"
set "MM=%dt:~4,2%"
set "DD=%dt:~6,2%"
set "HH=%dt:~8,2%"
set "NN=%dt:~10,2%"

set TIMESTAMP=%YY%%MM%%DD%_%HH%%NN%

go build -o "file_graph_%TIMESTAMP%.exe" .

if exist "file_graph_%TIMESTAMP%.exe" (
    echo Done! Output: file_graph_%TIMESTAMP%.exe
) else (
    echo Build failed!
)
