cd /ROOT/ginEssential
git pull
pkill -9 ding
go build ./cmd/ding
nohup ./ding daemon > log.out 2>&1 &

git commit -a -m "sss"
git push origin master

go env -w GOOS=linux
go build ./cmd/ding
go env -w GOOS=windows

ssh root@101.43.116.91

pkill -9 ding

拷贝 D:\cheng\ginEssential /ROOT/ginEssential

cd /ROOT/ginEssential
chmod 777 ding
nohup ./ding daemon > log.out 2>&1 &

go env -w GOOS=linux
go env -w GOOS=windows

安装protoc
go get -u github.com/golang/protobuf

protoc -I . --go_out=plugins=grpc:. proto/helloworld.proto

