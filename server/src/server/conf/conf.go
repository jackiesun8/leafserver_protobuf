package conf

var (
	// gate conf
	Encoding               = "protobuf" // "json" or "protobuf"
	PendingWriteNum        = 2000
	LenMsgLen              = 4 //消息头长度4个字节
	MinMsgLen       uint32 = 2
	MaxMsgLen       uint32 = 4096
	LittleEndian           = false

	// skeleton conf
	GoLen              = 10000
	TimerDispatcherLen = 10000
	ChanRPCLen         = 10000
)
