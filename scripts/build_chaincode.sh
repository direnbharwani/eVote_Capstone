#!/bin/bash
cd chaincode

if [ ! -d ./bin/ ]; then
    mkdir bin
fi

echo "Building evote_poc chaincode to bin/"

start=$(date +%s%N) # %N for elapsed time in seconds

GOOS=linux GOARCH=amd64 go build -o ./bin/evote_poc.bin

end=$(date +%s%N)

echo "Build complete!"
#Convert time diff to milliseconds
echo "Time taken: $(($(($end-$start)) / 1000000)) ms"