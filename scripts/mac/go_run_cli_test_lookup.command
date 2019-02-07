here="$(dirname "$BASH_SOURCE")"
cd $here/../..
export GOPATH=$HOME/go
go run cmd/ipblock-pool/main.go lookup --key main1.sub1