GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build
rm -rf Medville.app
mkdir -p Medville.app/Contents/MacOS
mkdir -p Medville.app/Contents/Resources
cp res/MacOS/Info.plist Medville.app/Contents/
cp res/MacOS/medville.icns Medville.app/Contents/Resources/
cp medvil Medville.app/Contents/MacOS
cp -r icon Medville.app/Contents/Resources
cp -r texture Medville.app/Contents/Resources
cp -r samples Medville.app/Contents/Resources
mkdir -p Medville.app/Contents/Resources/saved
cp saved/example* Medville.app/Contents/Resources/saved
