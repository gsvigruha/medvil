export __NV_PRIME_RENDER_OFFLOAD=1
export __GLX_VENDOR_LIBRARY_NAME=nvidia
# https://pkg.go.dev/cmd/cgo
export GODEBUG=cgocheck=2

go build
./medvil
