cd /ROOT/ginEssential
git pull
pkill -9 ding
go build ./cmd/ding
nohup ./ding daemon > log.out 2>&1 &

go env -w GOOS=linux
go build ./cmd/ding
go env -w GOOS=windows
git commit -a -m "sss"
git push origin master


cd /ROOT/ginEssential
git checkout .
git pull
pkill -9 ding
chmod 777 ding
nohup ./ding daemon > log.out 2>&1 &

go env -w GOOS=linux
go env -w GOOS=windows

安装protoc
go get -u github.com/golang/protobuf

protoc -I . --go_out=plugins=grpc:. proto/helloworld.proto

