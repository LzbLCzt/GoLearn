#--go_out=plugins=grpc:.: 这部分指定了生成的输出文件的语言和插件。
# --go_out 指定了输出的目标语言是 Go。
# plugins=grpc 表明除了生成标准的 Go 代码外，还要使用 gRPC 插件生成 gRPC 相关的代码。
# 最后的 . 表示生成的文件将被输出到当前目录
protoc --go_out=plugins=grpc:. ./programmer.proto
protoc --go_out=plugins=grpc:. ./hello.proto