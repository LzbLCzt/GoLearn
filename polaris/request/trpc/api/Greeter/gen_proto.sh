#!/usr/bin/env bash

cd ../../../../..
trpc create --protofile=polaris/request/trpc/api/Greeter/helloworld.proto --rpconly --mock=false -o polaris/request/trpc/api/Greeter