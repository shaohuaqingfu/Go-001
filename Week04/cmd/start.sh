# shellcheck disable=SC2164
cd /data/home/wanghan/GoProjects/Go-001/
protoc --go_out=plugins=grpc:. Week04/api/product/product.proto
protoc --grpc-gateway_out=logtostderr=true:. Week04/api/product/product.proto
