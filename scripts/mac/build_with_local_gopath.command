here="$(dirname "$BASH_SOURCE")"
cd $here/..
./clean.sh
./local_gopath.sh
. ./env.sh
./build.sh
