CGO_LDFLAGS="-mmacosx-version-min=11.6 -framework Cocoa -framework OpenGL -framework IOKit -framework CoreVideo" \
  CGO_CFLAGS="-mmacosx-version-min=11.6 -D_GLFW_COCOA -D_GLFW_USE_RETINA" \
  GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 \
  go run main.go