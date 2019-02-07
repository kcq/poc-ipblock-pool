here="$(dirname "$BASH_SOURCE")"
cd $here/../..
export GOPATH=$HOME/go
go run cmd/ipblock-pool-server/main.go