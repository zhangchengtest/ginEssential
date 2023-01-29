cd /ROOT/ginEssential
git pull
go build ./cmd/ding
pkill -9 ding
nohup ./ding daemon > log.out 2>&1 &


安装protoc
go get -u github.com/golang/protobuf

protoc -I . --go_out=plugins=grpc:. proto/helloworld.proto

