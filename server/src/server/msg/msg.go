package msg

import (
	"errors"
	"fmt"
	"github.com/name5566/leaf/network/json"
	"github.com/name5566/leaf/network/protobuf"
	"os"
	"reflect"
)

var (
	JSONProcessor     = json.NewProcessor()
	ProtobufProcessor = protobuf.NewProcessor()
)

func init() {
	//每加一条消息，都需要在这里将消息注册进处理器内

	//Login.proto
	ProtobufProcessor.Register(&RegisterReq{})
	ProtobufProcessor.Register(&RegisterRes{})
	ProtobufProcessor.Register(&LoginReq{})
	ProtobufProcessor.Register(&LoginRes{})
	ProtobufProcessor.Register(&RoomListReq{})
	ProtobufProcessor.Register(&RoomListRes{})
	ProtobufProcessor.Register(&EnterRoomReq{})
	ProtobufProcessor.Register(&ExitRoomReq{})
	//最后生成MsgID.lua文件
	genMsgID()
}

func genMsgID() {
	go_path := os.Getenv("GOPATH")
	if go_path == "" {
		panic(errors.New("GOPATH is not set"))
	}
	file, err := os.OpenFile(go_path+"/../client/QiPai/src/MsgID.lua", os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	// _, err = file.WriteString("MsgName = {\n")
	// if err != nil {
	// 	panic(err)
	// }
	// ProtobufProcessor.Range(func(id uint16, t reflect.Type) {
	// 	str := fmt.Sprintf("\t[%d] = \"%s\"\n", id, t.Elem().Name())
	// 	_, err = file.WriteString(str)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// })
	// _, err = file.WriteString("}\n")
	// if err != nil {
	// 	panic(err)
	// }

	_, err = file.WriteString("msgID = {\n")
	if err != nil {
		panic(err)
	}
	ProtobufProcessor.Range(func(id uint16, t reflect.Type) {
		str := fmt.Sprintf("\t%s = %d\n,", t.Elem().Name(), id)
		_, err = file.WriteString(str)
		if err != nil {
			panic(err)
		}
	})
	_, err = file.WriteString("}\n")
	if err != nil {
		panic(err)
	}
}
