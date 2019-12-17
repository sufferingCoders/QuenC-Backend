# QuenC 昆客

目前QuenC前端是使用Firestore來存取資料, 但Firestore有許多的限制, 希望能用 Golang + MongoDB 來取代 Firestore.

# 待討論
- [ ] 該使用WebSocket的Schema


# 待實現

|完成|作物|優先級|描述|耕種人|完成時間|
|:---:|:---:|:---:|:---:|:---:|:---:|
|<ul><li>- [x] </li></ul>|Flutter <-> Golang 對接測試|Must Have| 測試WebSocket 和 Flutter的對接 | Richard | 18 Dec 2019 |
|<ul><li>- [ ] </li></ul>|Golang <-> MongoDB 對接測試|Must Have| 測試Subscription 和 MongoDB的對接 | Richard | |



# Golang 資源

### 依賴處理
[[Go Module]](https://openhome.cc/Gossip/Go/Module.html)
[[Go Websocket]](https://zhuanlan.zhihu.com/p/35167916)

# WebSocket
在後端我們使用WebSocket來實現實時通信, 主要用來傳送任何有關User的更改.

我們使用的主要Package是 [Gorilla WebSocket](https://github.com/gorilla/websocket)

使用上我們將WebSocket放進 [Gin](https://github.com/gin-gonic/gin) 的框架中

我們設置一個GET節點後, 將此節點升級為WebSocket

```go
ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
```

之後我們再用Loop去管理接收到的訊息, 因為是雙向通道, 所以伺服器端也能主動傳送訊息

```go

		for {
			mt, message, err := ws.ReadMessage()

			log.Printf("Get message %+v", string(message))

			if err != nil {
				log.Printf("Error occur: %+v\b", err)
				break
			}

			if string(message) == "ping" {
				message = []byte("pong")
			}

			err = ws.WriteMessage(mt, message)

			if err != nil {
				break
			}
		}
    
```
