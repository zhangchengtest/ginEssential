cd /ROOT/ginEssential
git pull
pkill -9 ding
go build ./cmd/ding
nohup ./ding daemon > log.out 2>&1 &


go env -w GOOS=linux
go env -w GOOS=windows

安装protoc
go get -u github.com/golang/protobuf

protoc -I . --go_out=plugins=grpc:. proto/helloworld.proto

