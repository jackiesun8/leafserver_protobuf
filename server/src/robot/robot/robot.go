//add  by bingyu

package robot

import (
	//"fmt"
	"bytes"
	"encoding/binary"
	"net"
	"robot/elog"
	//"os"
	//"os/signal"
	"sync"
	"sync/atomic"
	//"syscall"
	"encoding/hex"
	"errors"
	"time"
)

const defRobotNum = 1 //默认机器人数

const gateAddr = "127.0.0.1:3563" //网关地址

const defDur = time.Duration(1 * time.Millisecond / 60) //默认ticker时间间隔

//机器人状态
const (
	kRobotActive = iota
	kRobotClose
)

//等待组包装
type WaitGroupWrapper struct {
	sync.WaitGroup
}

//执行一个回调函数
func (w *WaitGroupWrapper) Run(cb func()) {
	w.Add(1)
	elog.LogDebug("------ wait gropu add 1")
	go func() {
		cb()
		w.Done()
		elog.LogDebug(" wait gropu del  1")
	}()
}

// //消息
// type Message struct {
// 	id   uint16 //消息ID
// 	data []byte //数据
// }

// func (m *Message) getMsgLen() {
// 	return 2 + len
// }

// //设置消息ID
// func (m *Message) SetId(id int32) {
// 	m.id = id
// }

// //获取消息ID
// func (m *Message) GetId() int32 {
// 	return m.id
// }

//更新回调
type UPDATE_CALLBACK func(r *Robot)

//机器人
type Robot struct {
	id        uint32            //机器人ID
	addr      *net.TCPAddr      //地址
	tcpConn   *net.TCPConn      //连接
	recvBuf   *bytes.Buffer     //接收缓冲
	sendBuf   *bytes.Buffer     //发送缓冲
	delCh     chan uint32       //删除通道
	closeCh   chan bool         //关闭通道
	sendCh    chan Message      //发送通道
	recvCh    chan Message      //接收通道
	state     uint32            //状态
	waitGroup *WaitGroupWrapper //等待组包装
	ticker    *time.Ticker      //时钟
	upCb      UPDATE_CALLBACK   //更新回调
	UserData  interface{}       //用户数据
}

//创建机器人
func NewRobot(chDel chan uint32, chClose chan bool, wg *WaitGroupWrapper) *Robot {
	r := &Robot{
		recvBuf:   &bytes.Buffer{},
		sendBuf:   &bytes.Buffer{},
		delCh:     chDel,
		closeCh:   make(chan bool, 1), //没有用机器人管理器传过来的关闭通道，而是自己创建了一个
		waitGroup: wg,
		state:     kRobotClose,
		sendCh:    make(chan Message, 100),
		recvCh:    make(chan Message, 100),
	}
	return r
}

//关闭机器人
func (r *Robot) Close() {

	if atomic.CompareAndSwapUint32(&r.state, kRobotActive, kRobotClose) { //如果状态为active，则返回true，并且设置状态为close
		r.delCh <- r.id  //这里告诉机器人管理器删除机器人
		close(r.closeCh) //通知关闭机器人
		elog.LogInfo(" robot:%d close ", r.id)
		r.ticker.Stop()   //关闭时钟
		r.tcpConn.Close() //关闭连接
	}
}

//发送消息
func (r *Robot) SendMsg(msg Message) {
	r.sendCh <- msg //发送到通道
	elog.LogInfo(" send  msg : ", msg.id)
}

//发送循环
func (r *Robot) SendLoop() {
	defer r.Close()
	elog.LogDebug(" robot %d  send  loop run  ", r.id)
	for {
		select {
		case msg := <-r.sendCh: //从发送通道取出待发送消息
			elog.LogInfo("  encode  msg :%d ", msg.id)
			binary.Write(r.sendBuf, binary.LittleEndian, msg.id)   //写ID
			binary.Write(r.sendBuf, binary.LittleEndian, msg.Data) //写数据
			byte := r.sendBuf.Bytes()
			n, err := r.tcpConn.Write(byte)
			if err != nil {
				elog.LogSysln(" conn ", r.id, " write data fail :", err)
				return
			}

			elog.LogInfo(" write  msg :%d ", n, hex.Dump(r.sendBuf.Bytes()))
		case <-r.closeCh: //机器人关闭信号
			elog.LogDebug(" send loop begin close ")
			return
		}
	}

	elog.LogDebug("send loop close  ")
}

//接收消息
func (r *Robot) RecvMsg() (msg Message, err error) {
	select {
	case msg = <-r.recvCh: //从接收通道取出待接收消息
		elog.LogInfo("receieve msg :%s ", msg.id)
		err = nil
		return
	default:
	}
	err = errors.New(" not msg ")
	return
}

//接收循环
func (r *Robot) RecvLoop() {
	defer r.Close()
	elog.LogDebug(" robot %d  read loop run  ", r.id)
	buf := make([]byte, 1024*1024)
	for {
		select {
		case <-r.closeCh: //机器人关闭信号
			elog.LogSysln("read loop begin stop ")
			return
		default:
		}

		r.tcpConn.SetDeadline(time.Now().Add(1e9))
		n, err := r.tcpConn.Read(buf)
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			elog.LogErrorln(" read data error ", err)
			return
		}
		elog.LogSys(" ********* read data : %d ", n)
		r.recvBuf.Write(buf[0:n])

		var msg Message
		binary.Read(r.recvBuf, binary.LittleEndian, msg.id)
		msg.Data = make([]byte, r.recvBuf.Len())
		r.recvBuf.Read(msg.Data)

		r.recvCh <- msg
	}

}

//启动机器人
func (r *Robot) Run() {

	elog.LogDebug(" robot %d begin run ", r.id)
	var err error
	r.addr, err = net.ResolveTCPAddr("tcp", gateAddr) //解析服务器地址
	if err != nil {
		elog.LogErrorln(gateAddr, ":resolve tcp addr fail, please try: 127.0.0.1:3563, fail: ", err)
		return
	}
	r.tcpConn, err = net.DialTCP("tcp", nil, r.addr) //连接服务器
	if err != nil {
		elog.LogErrorln("connect server fail , because :", err)
		return
	}

	elog.LogInfoln(" connect  server sucess :", r.id, r.tcpConn.RemoteAddr().String())
	r.state = kRobotActive //设置机器人状态

	r.waitGroup.Run(r.SendLoop) //启动发送循环
	r.waitGroup.Run(r.RecvLoop) //启动接收循环

	r.ticker = time.NewTicker(defDur) //启动机器人时钟

	defer r.Close() //延迟关闭机器人
	for {
		select {
		case <-r.closeCh: //关闭通道
			elog.LogDebugln("update loop begin stop ")
			return
		case <-r.ticker.C: //时钟到时
			//逻辑处理
			r.update()
			//elog.LogInfoln("robot :", r.id, "heart :", t)
		default:
		}
	}
}

//设置更新回调
func (r *Robot) SetUpdateCb(cb UPDATE_CALLBACK) {
	r.upCb = cb
}

//获取机器人ID
func (r *Robot) GetId() uint32 {
	return r.id
}

//更新
func (r *Robot) update() {
	if r.upCb != nil {
		r.upCb(r) //调用更新回调
	}
}

//机器人管理器
type RobotMng struct {
	lastRobotId uint32            //上一个机器人的ID
	robots      map[uint32]*Robot //机器人表
	mapMutex    sync.Mutex        //表锁
	delCh       chan uint32       //删除通道
	ticker      *time.Ticker      //时钟
	waitGroup   *WaitGroupWrapper //等待组包装
	closeCh     chan bool         //关闭通道
	upCb        UPDATE_CALLBACK   //更新回调
}

//设置更新回调
func (rbMng *RobotMng) SetUpdateCb(cb UPDATE_CALLBACK) {
	rbMng.upCb = cb
}

//添加机器人
func (rbMng *RobotMng) AddRobot(r *Robot) {
	rbMng.mapMutex.Lock()         //加锁
	defer rbMng.mapMutex.Unlock() //解锁
	rbMng.robots[r.id] = r        //添加映射
	elog.LogDebugln("add robot ", r.id)
}

//删除机器人
func (rbMng *RobotMng) DelRobot(id uint32) {
	rbMng.mapMutex.Lock()         //加锁
	defer rbMng.mapMutex.Unlock() //解锁
	delete(rbMng.robots, id)      //删除映射
	elog.LogDebugln("del robot ", id)
}

//创建机器人
func (rbMng *RobotMng) NewRobot() {

	r := NewRobot(rbMng.delCh, rbMng.closeCh, rbMng.waitGroup) //创建机器人
	r.id = atomic.AddUint32(&rbMng.lastRobotId, 1)             //计算机器人ID
	rbMng.AddRobot(r)                                          //添加机器人
	elog.LogDebug(" create robot %d ", r.id)
	r.SetUpdateCb(rbMng.upCb)  //设置更新回调
	rbMng.waitGroup.Run(r.Run) //启动机器人
}

//关闭机器人管理器
func (rbMng *RobotMng) Close() {

	elog.LogDebug(" .......................begin close ........................")
	rbMng.ticker.Stop() //先关闭机器人管理器时钟,这样就不会检测缺少的机器人
	elog.LogDebug(" rbMng ticket stop")
	for _, r := range rbMng.robots { //关闭所有机器人
		r.Close() //内部会告诉机器人管理器删除机器人
	}
	close(rbMng.closeCh) //关闭通道，机器人管理器心跳结束
	elog.LogDebug(" rbmng  closech already close ")
	rbMng.waitGroup.Wait() //等待所有机器人关闭
	elog.LogDebug(" rbMng  wait all robot close ")
	//最后关del ch
	close(rbMng.delCh) //关闭删除通道
	elog.LogDebug(" rbMng delCh close ")
	elog.LogDebug(" Everything is ok, i quit ....................................")
}

//机器人管理器心跳
func (rbMng *RobotMng) Heart() {

	//500 ms 检测一次
	rbMng.ticker = time.NewTicker(defDur) //创建时钟
	for {
		// del close
		select {
		case id := <-rbMng.delCh: //删除机器人
			rbMng.DelRobot(id)
		case <-rbMng.closeCh: //关闭机器人管理器心跳
			elog.LogDebugln(" rbMng heart close ")
			return
		case <-rbMng.ticker.C:
			//定时检测
			diff := defRobotNum - len(rbMng.robots) //计算缺少的机器人
			if diff > 0 {
				elog.LogDebugln(" RobotMng heart : ", "diff ", diff)
				for i := 0; i < diff; i++ {
					rbMng.NewRobot() //补足缺少的机器人
				}
			}
		default:
		}
	}
}

//启动机器人管理器
func (rbMng *RobotMng) Run() {
	go rbMng.waitGroup.Run(rbMng.Heart) //启动心跳
}

//创建机器人管理器
func NewRobotMng() *RobotMng {
	robotMng := &RobotMng{
		robots:      make(map[uint32]*Robot),
		lastRobotId: 0,
		delCh:       make(chan uint32, 100),
		closeCh:     make(chan bool),
		waitGroup:   &WaitGroupWrapper{},
	}

	return robotMng
}
