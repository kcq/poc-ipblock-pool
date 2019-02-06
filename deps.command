here="$(dirname "$BASH_SOURCE")"
cd $here
export GOPATH=$HOME/go
go get ./...