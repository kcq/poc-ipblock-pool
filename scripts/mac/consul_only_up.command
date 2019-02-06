here="$(dirname "$BASH_SOURCE")"
cd $here
docker run -it --rm --name="consul_only" -p 8500:8500 consul:1.4.2 agent -dev -ui -client 0.0.0.0

