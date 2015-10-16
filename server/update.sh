SERVER_PATH=$GOPATH/src/server
PROTO_PATH=$GOPATH/../share/proto
cd $PROTO_PATH
protoc --go_out=$SERVER_PATH/msg *.proto
go install server

