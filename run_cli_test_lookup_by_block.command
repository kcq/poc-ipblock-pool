here="$(dirname "$BASH_SOURCE")"
cd $here
export GOPATH=$HOME/go
go run cmd/ipblock-pool/main.go lookup --block 169.254.51.0