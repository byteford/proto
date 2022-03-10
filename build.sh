echo "building protobuf"
docker run --rm -u $(id -u) -v${PWD}:${PWD} -w${PWD} jaegertracing/protobuf:latest --proto_path=${PWD}/proto \
    --go_out=${PWD}/server --python_out=${PWD}/client -I/usr/include/github.com/gogo/protobuf ${PWD}/proto/click.proto
docker run --rm -u $(id -u) -v${PWD}:${PWD} -w${PWD} jaegertracing/protobuf:latest --proto_path=${PWD}/proto \
    --python_out=${PWD} -I/usr/include/github.com/gogo/protobuf ${PWD}/proto/click.proto
echo "building go"

