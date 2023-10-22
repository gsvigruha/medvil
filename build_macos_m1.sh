CGO_LDFLAGS="-mmacosx-version-min=11.6 -framework Cocoa -framework OpenGL -framework IOKit -framework CoreVideo" \
  CGO_CFLAGS="-mmacosx-version-min=11.6 -D_GLFW_COCOA -D_GLFW_USE_RETINA" \
  GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 \
  go build
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
cp settings.json.default Medville.app/Contents/Resources/settings.json
