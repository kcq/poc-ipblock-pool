here="$(dirname "$BASH_SOURCE")"
cd $here
docker run -it --rm --name="service_only" -p 5555:5555 poc/ipblock-pool

