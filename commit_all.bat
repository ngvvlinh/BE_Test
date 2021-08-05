@echo off

git add *
git commit -m "Auto commit"
git push -v

if NOT ["%errorlevel%"]==["0"] (
    pause
    exit /b %errorlevel%
)
