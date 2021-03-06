package chat

import (
	"context"
	"deercoder-chat/chat-srv/proto"
	"encoding/json"
	uuid "github.com/dreamlu/go.uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

//客户端
type Client struct {
	GroupID string // 标识客户端
	UID     int64  // 唯一标识用户id
	Conn    *websocket.Conn
}

var clients []*Client //客户端队列,指针同步同一个client data
//var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan *proto.Message) // broadcast channel

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 消息读取
// 开启不同进程代表对应的客户端通信
func WsHander(cli proto.StreamerService, ws *websocket.Conn) {

	defer ws.Close()

	// Register our new client
	//注册客户端连接
	var ct Client
	ct.Conn = ws
	//放入连接队列
	clients = append(clients, &ct)

	// go micro stream
	// there is error
	// always say can not find chat micro
	// var stream 	proto.Streamer_ServerStreamService

	//消息读取,每个客户端数据
	for {
		var req proto.Request
		// var msg chat.Message
		// Read in a new message as JSON and map it to a Message object

		_, data, err := ws.ReadMessage()
		//log.Printf("[消息内容]: %v", data)
		if err != nil {
			log.Printf("[错误-read]: %v", err)
		}

		if string(data) == "ping" {
			log.Printf("[心跳检测]: %v", string(data))
			continue
		}
		err = json.Unmarshal(data, &req.Message) //ws.ReadJSON(&req.Message)
		//log.Println("[消息内容]: ", req.Message)
		if err != nil {
			log.Printf("[错误-read]: %v", err)
			//delete(clients, ws) //删除对应连接
			for _, v := range clients { //删除对应连接,emm...暂时先遍历删除～
				//fmt.Println(v)
				if v.Conn == ws {
					break
				}
			}
			break
		}
		log.Println("[聊天测试]: ", req.Message)

		// 消息体以及个人身份绑定
		// 这部分有待抽取出去
		// 后续完善
		ct.GroupID = req.Message.GroupId //客户端唯一标识
		ct.UID = req.Message.FromUid
		// 个人身份绑定事件
		if req.Message.Content == "" {
			continue
		}

		// 消息唯一id
		req.Message.Uuid = uuid.NewV1().String()
		// Send the newly received message to the broadcast channel
		broadcast <- req.Message

		// send broadcast, then save the message
		//_ = chat.CreateGroupMsg(req.Message.Uuid, req.Message.GroupId, req.Message.FromUid, req.Message.Content, req.Message.ContentType)
		// use go-micro stream deal with the emessage
		// Send request to stream server
		// rpc service
		// 异步消息写入
		go CreateGroupMsg(req)

	}

	//defer stream.Close()
}

// 消息写入
// 消息推送(不通进程代表各自客户端的写入进程)
func handleMessages() {
	for {
		msg := <-broadcast //广播
		//获得广播内容,筛选获得group_id,并将内容分发给指定的客户
		for k, client := range clients {
			// send message to every specified client, hehe~
			if client.GroupID != msg.GroupId { // must same group_id
				continue
			}
			// next have same group_id
			err := client.Conn.WriteJSON(msg)
			if err != nil { //当连接挂了
				//fmt.Println("客户:",client,"聊天记录写入失败")
				log.Printf("[错误-write]: %v", err)
				client.Conn.Close()
				clients = append(clients[:k], clients[k+1:]...)
				////记录该用户最后读的消息id,广播中处理,待优化
				//chat.CreateGroupLastMsg(msg.GroupId,msg.FromUid,msg.Flag,msg.ID)
				continue
			}
		}
		//连接该断的也断了
		//进行用户在线检测
		//	gusers := chat.GetChatUsers(msg.GroupId)
		//into:
		//	for _, v2 := range clients {
		//		if v2.GroupID == msg.GroupId { //在线用户
		//			for k, v := range gusers {
		//				if v2.UID == v.Uid {
		//					gusers = append(gusers[:k], gusers[k+1:]...) //去除在线用户
		//					goto into
		//				}
		//			}
		//		}
		//	}
		// 剩下的为群聊离线用户
		// 记录离线消息
		//for _, v := range gusers{
		//	_ = chat.CreateGroupLastMsg(msg.GroupId, v.Uid, msg.Uuid)
		//}
	}
}

// 聊天记录创建
func CreateGroupMsg(req proto.Request) {
	_, err := ChatClient.CreateGroupMsg(context.TODO(), &req)
	if err != nil {
		log.Println("[错误]：", err)
	}
}
