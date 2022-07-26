go build ./cmd/ding
nohup ./ding daemon > log.out 2>&1 &