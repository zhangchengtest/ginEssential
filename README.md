git pull
go build ./cmd/ding
pkill -9 ding
nohup ./ding daemon > log.out 2>&1 &