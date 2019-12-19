# QuenC 昆客

目前QuenC前端是使用Firestore來存取資料, 但Firestore有許多的限制, 希望能用 Golang + MongoDB 來取代 Firestore.

# 待討論
- [ ] 該使用WebSocket的Schema


# 待實現

|完成|作物|優先級|描述|耕種人|完成時間|
|:---:|:---:|:---:|:---:|:---:|:---:|
|<ul><li>- [x] </li></ul>|Flutter <-> Golang 對接測試|Must Have| 測試WebSocket 和 Flutter的對接 | Richard | 18 Dec 2019 |
|<ul><li>- [x] </li></ul>|Insert|Must Have| 測試注入Data至MongoDB | Richard | 20 Dec 2019 |
|<ul><li>- [x] </li></ul>|Update|Must Have| 測試Update MongoDB Document| Richard | 20 Dec 2019 |
|<ul><li>- [x] </li></ul>|Listen to Change Stream|Must Have| 測試MongoDB ChangeStream | Richard | 20 Dec 2019 |
|<ul><li>- [ ] </li></ul>|寫MongoDB對接方法 在ReadMe|Must Have| 描述MongoDB對接 | Richard | |
|<ul><li>- [ ] </li></ul>|遷移User|Must Have| 建造User | Richard | |




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


# Flutter <-> Golang <-> MongoDB 對接

## 為什麼需要這組對接?
我們要用這組對接來取代**FireStore**, 將多數的業務邏輯轉移到後端, 同時我們需要後端能提供一個Stream來告訴前端, "Database 中的資料改變了, 改變後的結果是 xxx", 這樣的實時通訊我們會用來使用在User的Schema和聊天室的功能上, Post和Comment則會繼續使用一般的Http Request。

## 什麼是Stream?
Stream可以看成是一個通道, 而我們這個使用狀況下的Stream, 則是我們追蹤(Watch/Subscribe)的資料在**MongoDB**中有被Updated的時候, 才會將改變的細節丟入通道中, 前端即可接收到這筆資料, 並重新建置Flutter的Widget。這篇[StreamBuilder的介紹](https://www.youtube.com/watch?v=MkKEWHfy99Y)和[Stream的介紹](https://youtu.be/nQBpOIHE4eE)也簡短的解釋了在Flutter中怎麼使用Stream。

在做State Management的時候, 也有一個叫做 **Rx** 系列的方法, [RxDart](https://pub.dev/packages/rxdart)中的 Observable 和 StreamController 都提供很多語法糖讓你輕鬆地應對Stream, 但是在我們的專案中Provider已經可以升任大部分的工作所以我們沒有採用這個比較重型的RxDart

## 對接

### Retrieve & Update

首先我們需要兩個 RESTful API 的基本操作, Retreive 和 Update 所以我們看一下這兩個操作要怎麼在和端和前端執行

#### 先創立一個Test Schema在後端

```go
type Testing struct {
	ID    primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Email string             `json:"email" bson:"email"`
}
```


