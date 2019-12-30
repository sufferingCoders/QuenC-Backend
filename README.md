# QuenC 昆客

目前QuenC前端是使用Firestore來存取資料, 但Firestore有許多的限制, 希望能用 Golang + MongoDB 來取代 Firestore.


- [待討論](#待討論)
- [待實現](#待實現)
- [Golang資源](#Golang-資源)
	- [依賴處理](#依賴處理)
- [WebSocket](#WebSocket)
- [Flutter Golang MongoDB 對接](#Flutter-Golang-MongoDB-對接)
	- [為什麼需要這組對接?](#為什麼需要這組對接)
	- [什麼是Stream?](#什麼是Stream)
	- [對接](#對接)
		- [Retrieve and Update](#Retrieve-and-Update)
		- [數據監控](#數據監控)
- [如何取得你的IP](#如何取得你的IP)
	- [Windows](#Windows)
	- [MacOS](#MacOS)
- [滿滿的坑](#滿滿的坑)
	
	


# 待討論
- [ ] 該使用WebSocket的Schema


# 待實現

|完成|作物|優先級|描述|耕種人|完成時間|
|:---:|:---:|:---:|:---:|:---:|:---:|
|<ul><li>- [x] </li></ul>|Flutter <-> Golang 對接測試|Must Have| 測試WebSocket 和 Flutter的對接 | Richard | 18 Dec 2019 |
|<ul><li>- [x] </li></ul>|Insert|Must Have| 測試注入Data至MongoDB | Richard | 20 Dec 2019 |
|<ul><li>- [x] </li></ul>|Update|Must Have| 測試Update MongoDB Document| Richard | 20 Dec 2019 |
|<ul><li>- [x] </li></ul>|Listen to Change Stream|Must Have| 測試MongoDB ChangeStream | Richard | 20 Dec 2019 |
|<ul><li>- [x] </li></ul>|寫MongoDB對接方法 在ReadMe|Must Have| 描述MongoDB對接 | Richard | 20 Dec 2019 |
|<ul><li>- [x] </li></ul>|User Schema|Must Have| 設計User Schema | Richard | 21 Dec 2019 |
|<ul><li>- [x] </li></ul>|Post Schema|Must Have| 設計Post Schema | Richard | 21 Dec 2019 |
|<ul><li>- [x] </li></ul>|PostCategory Schema|Must Have| 設計PostCategory Schema | Richard | 21 Dec 2019 |
|<ul><li>- [x] </li></ul>|Comment Schema|Must Have| 設計Comment Schema | Richard | 21 Dec 2019 |
|<ul><li>- [x] </li></ul>|Report Schema|Must Have| 設計Report Schema | Richard | 21 Dec 2019 |
|<ul><li>- [x] </li></ul>|User Api|Must Have| 完成 User Api 搭建 | Richard | 21 Dec 2019 |
|<ul><li>- [x] </li></ul>|Post Api|Must Have| 完成 Post Api 搭建 | Richard | 21 Dec 2019 |
|<ul><li>- [x] </li></ul>|PostCategory Api|Must Have| 完成 PostCategory Api 搭建 | Richard | 21 Dec 2019 |
|<ul><li>- [x] </li></ul>|Comment Api|Must Have| 完成 Comment Api 搭建 | Richard | 21 Dec 2019 |
|<ul><li>- [x] </li></ul>|Report Api|Must Have| 完成 Report Api 搭建 | Richard | 21 Dec 2019 |
|<ul><li>- [x] </li></ul>|User Router|Must Have| 完成 User Router 規劃 | Richard | 21 Dec 2019 |
|<ul><li>- [x] </li></ul>|Post Router|Must Have| 完成 Post Router 規劃 | Richard | 21 Dec 2019 |
|<ul><li>- [x] </li></ul>|PostCategory Router|Must Have| 完成 PostCategory Router 規劃 | Richard | 21 Dec 2019 |
|<ul><li>- [x] </li></ul>|Comment Router|Must Have| 完成 Comment Router 規劃 | Richard | 21 Dec 2019 |
|<ul><li>- [x] </li></ul>|Report Router|Must Have| 完成 Report Router 規劃 | Richard | 21 Dec 2019 |
|<ul><li>- [x] </li></ul>|User Stream|Must Have| 搭建WebSocket&ChangedStream給前端Subscribe MongoDBw特定User | Richard | 21 Dec 2019 |
|<ul><li>- [x] </li></ul>|將User下的LikePosts和LikeComments取消|Must Have| 在Comment&Post中增加 Likers field| Richard | 26 Dec 2019 |
|<ul><li>- [x] </li></ul>|刪除AuthorGender&AuthorDomain|Must Have| 用Population的方式來取代 | Richard | 26 Dec 2019 |
|<ul><li>- [x] </li></ul>|User migrate to Aggregation|Must Have| 使用Aggregation和Pipeline來做排序和Pupolation | Richard | 26 Dec 2019 |
|<ul><li>- [x] </li></ul>|Report migrate to Aggregation|Must Have| 使用Aggregation和Pipeline來做排序和Pupolation | Richard | 26 Dec 2019 |
|<ul><li>- [x] </li></ul>|Post migrate to Aggregation|Must Have| 使用Aggregation和Pipeline來做排序和Pupolation | Richard | 26 Dec 2019 |
|<ul><li>- [x] </li></ul>|Comment migrate to Aggregation|Must Have| 使用Aggregation和Pipeline來做排序和Pupolation | Richard | 26 Dec 2019 |
|<ul><li>- [x] </li></ul>|Aggregation Complete |Must Have| 為規劃的Aggregation Function加上Hander&Router| Richard | 26 Dec 2019 |
|<ul><li>- [x] </li></ul>| 測試 User 登入&註冊 |Must Have| 測試Flutter前端是否能順利和後端API連接候登入, 並拿到JWT | Richard | 30 Dec 2019 |
|<ul><li>- [x] </li></ul>| 測試 UserStream |Must Have| 測試Flutter前端是否能順利經由後端監聽MongodDB中的UserStream| Richard | 30 Dec 2019 |
|<ul><li>- [x] </li></ul>| 測試 PostCategory 新增 |Must Have| 測試PostCategory的新增功能 | Richard | 30 Dec 2019 |
|<ul><li>- [ ] </li></ul>| 測試 PostCategory 刪除 |Must Have| 測試PostCategory的刪除功能 | Richard | |
|<ul><li>- [x] </li></ul>| 測試 Post |Must Have| 測試Post的新增和編輯功能 | Richard | 30 Dec 2019 |
|<ul><li>- [x] </li></ul>| 測試 Post GET |Must Have| 測試前端Flutter能請求到Post | Richard | 30 Dec 2019 |
|<ul><li>- [ ] </li></ul>| 測試 Post 刪除 |Must Have| 測試Post的刪除功能 | Richard |  |
|<ul><li>- [x] </li></ul>| Post 字數限制 |Must Have| Post的字數限制, 顯示在編輯區 | Richard | 30 Dec 2019  |
|<ul><li>- [x] </li></ul>| 壓縮上傳的圖片大小 |Must Have| 壓縮Firestore儲存的圖片大小 | Richard | 30 Dec 2019  |
|<ul><li>- [x] </li></ul>| Post 載入失敗時, 一樣顯示排序 |Must Have| Post 載入失敗時, 一樣顯示排序 | Richard | 30 Dec 2019  |
|<ul><li>- [ ] </li></ul>| 測試 Comment 新增 |Must Have| 測試Comment的新增功能 | Richard | |
|<ul><li>- [ ] </li></ul>| 測試 Comment 刪除 |Must Have| 測試Comment的刪除功能 | Richard | |
|<ul><li>- [ ] </li></ul>| 測試 Report 新增 |Must Have| 測試Report的新增功能 | Richard | |
|<ul><li>- [ ] </li></ul>| 測試 Report 編輯 |Must Have| 測試Report的編輯功能 | Richard | |
|<ul><li>- [ ] </li></ul>| 測試 Report 刪除 |Must Have| 測試Report的刪除功能 | Richard | |
|<ul><li>- [ ] </li></ul>|前端API路徑修正|Must Have| 連接上新的Backend | Richard |  |
|<ul><li>- [ ] </li></ul>|將CompleteField簡化成toAddingMap|Must Have| CompleteField中止需要對特定的Field做更改 | Richard |  |
|<ul><li>- [ ] </li></ul>|兩個API節點對應Unsolved 跟 Solved Reports|Must Have| 兩個API節點對應Unsolved 跟 Solved | Richard |  |


|<ul><li>- [x] </li></ul>|測試Primitive.ObjectID對接時候是可否可以轉成String|Must Have| 若無法則須新增另一Struct (可以轉為S) | Richard | 30 Dec 2019 |


# Image Compression 測試
|Quality|Size| Image |
|:---:|:---:|:---:|
| 100 | 5092852 |![圖片載入中...](https://firebasestorage.googleapis.com/v0/b/quenc-hlc.appspot.com/o/images%2F2019-12-30%2016%3A28%3A16.231209.png?alt=media&token=20ce35f7-740a-4aaa-971b-39fe5414b41b)|
| 85 | 1309107 |![圖片載入中...](https://firebasestorage.googleapis.com/v0/b/quenc-hlc.appspot.com/o/images%2F2019-12-30%2016%3A29%3A48.089381.png?alt=media&token=d201dbff-64bd-4f83-94de-c9b59df81535)|
| 70 | 684145 |![圖片載入中...](https://firebasestorage.googleapis.com/v0/b/quenc-hlc.appspot.com/o/images%2F2019-12-30%2016%3A30%3A38.745708.png?alt=media&token=8b6928b6-e717-4bae-a17f-ce92a3f97a9f)|
| 60 | 500781 |![圖片載入中...](https://firebasestorage.googleapis.com/v0/b/quenc-hlc.appspot.com/o/images%2F2019-12-30%2016%3A32%3A01.220338.png?alt=media&token=e0dea274-1738-4cc9-820a-769a0b90e139)|
| 50 | 407540 |![圖片載入中...](https://firebasestorage.googleapis.com/v0/b/quenc-hlc.appspot.com/o/images%2F2019-12-30%2016%3A32%3A41.956861.png?alt=media&token=4bb39d06-229b-4a76-8c9e-27e96a4a0cd0)|
| 40 | 327527 |![圖片載入中...](https://firebasestorage.googleapis.com/v0/b/quenc-hlc.appspot.com/o/images%2F2019-12-30%2016%3A33%3A06.446733.png?alt=media&token=48a10bab-ad4e-4e0c-9914-3a4ba138f1a8)|
| 30 | 260946 |![圖片載入中...](https://firebasestorage.googleapis.com/v0/b/quenc-hlc.appspot.com/o/images%2F2019-12-30%2016%3A33%3A42.170038.png?alt=media&token=86db26a0-eb93-4564-9e4b-213b3d0eb600)|
| 20 | 199159 |![圖片載入中...](https://firebasestorage.googleapis.com/v0/b/quenc-hlc.appspot.com/o/images%2F2019-12-30%2016%3A34%3A09.249090.png?alt=media&token=4025c6e9-53b9-4cf8-9e61-2fe077ec2d1a)|


# Golang 資源

### 依賴處理
[[Go Module]](https://openhome.cc/Gossip/Go/Module.html)
[[Go Websocket]](https://zhuanlan.zhihu.com/p/35167916)

# WebSocket
在後端我們使用 WebSocket 來實現實時通信, 主要用來傳送任何有關User的更改.

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


# Flutter Golang MongoDB 對接

## 為什麼需要這組對接?
我們要用這組對接來取代 **FireStore** , 將多數的業務邏輯轉移到後端, 同時我們需要後端能提供一個 Stream 來告訴前端, " Database 中的資料改變了, 改變後的結果是 xxx ", 這樣的實時通訊我們會用來使用在 User 的 Schema 和聊天室的功能上, Post 和 Comment 則會繼續使用一般的 Http Request。

## 什麼是Stream?
Stream 可以看成是一個通道, 而我們這個使用狀況下的 Stream, 則是我們追蹤 (Watch/Subscribe) 的資料在 **MongoDB** 中有被 Updated 的時候, 才會將改變的細節丟入通道中, 前端即可接收到這筆資料, 並重新建置 Flutter 的 Widget。這篇 [StreamBuilder的介紹](https://www.youtube.com/watch?v=MkKEWHfy99Y) 和 [Stream的介紹](https://youtu.be/nQBpOIHE4eE) 也簡短的解釋了在 Flutter 中怎麼使用 Stream。

在做 State Management 的時候, 也有一個叫做 **Rx** 系列的方法, [RxDart](https://pub.dev/packages/rxdart) 中的 Observable 和 StreamController 都提供很多語法糖讓你輕鬆地應對 Stream, 但是在我們的專案中 Provider 已經可以升任大部分的工作所以我們沒有採用這個比較重型的 RxDart。

## 對接

### Retrieve and Update

首先我們需要兩個 RESTful API 的基本操作, Retreive 和 Update 所以我們看一下這兩個操作要怎麼在和端和前端執行。

#### 先創立一個 Test Schema 在後端 (使用 Schema 我們定義資料剛怎麼存儲在 MongoDB )

它長的這樣, 有兩個 Fields, 一個 ID 是當我們將資料加入 MongoDB 中時會自動生成的, 另外一個則是我們可以自由決定的 "Email" 。Golang 中的 struct 可以看待成傳統 Java 或 C# 這種 OOP 中的 Class, 但它又沒有一些 Class 擁有的功能 (ex: Inheritance)。

```go
type Testing struct {
	ID    primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Email string             `json:"email" bson:"email"`
}
```

#### 在後端創立兩個 Handlers ( Creating & Updating )

我們使用 [Gin](https://github.com/gin-gonic/gin) 框架, 來創造一個 POST 和一個 PUT Handler.


```go

type TestInfo struct {
	Email string `json:"email" bson:"email"`
}

// Create Handler 用來接受前端傳來的資料後在MongoDB中加入這筆資料
router.POST("/test", func(c *gin.Context) {

	var testAdding TestInfo
	err := c.ShouldBindJSON(&testAdding) // 預期前端的ContentTypeHeader使用的是application/json

	if err != nil {
		errStr := fmt.Sprintf("Cannot bind the given info : %+v \n", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": errStr,
			"msg": "Cannot bind the given info",
		})
		return
	}

	testingClient := Testing{
		Email: testAdding.Email,
	}

	result, err := database.DB.Collection("test").InsertOne(context.TODO(), testingClient)

	if err != nil {
		errStr := fmt.Sprintf("Can't insert a test due to the error : %+v \n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errStr,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": result.InsertedID,
	})

})


// Update Handler 用來更新MongoDB已有的數據
router.PUT("/test/:id", func(c *gin.Context) {

	id := c.Param("id")

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		errStr := fmt.Sprintf("The given id cannot be transform to oid : %+v \n", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"err": errStr,
		})
		return
	}

	var testInfo TestInfo
	err = c.ShouldBindJSON(&testInfo)

	if err != nil {
		errStr := fmt.Sprintf("Cannot bind the given info: %+v \n", err)

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": errStr,
		})
		return
	}

	result, err := database.DB.Collection("test").UpdateOne(
		context.TODO(),
		bson.M{"_id": oid},
		bson.M{"$set": bson.M{"email": testInfo.Email}},
	)

	if err != nil {
		errStr := fmt.Sprintf("Cannot update a test due to the error: %+v \n", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errStr,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": result,
	})

})

```


#### 在前端中加入一個 Provider 來傳輸數據到後端


```dart

class WebScoketService with ChangeNotifier {

  static IOWebSocketChannel channel;
  static String insertedID;
  
  IOWebSocketChannel get currentChannel {
    return channel;
  }

  String get currentId {
    return insertedID;
  }
  
  
  /// 加入數據至後端
  Future<void> addTestDocument(String email) async {
    if (email == null || email.isEmpty) {
      return null;
    }
    final String url = "http://192.168.1.135:8080/test"; // 如果要在實機上測試則用Wifi下的IP, 下面會介紹怎麼取得
    
    final res = await http.post(
      url,
      headers: {
        HttpHeaders.contentTypeHeader: "application/json",
      },
      body: json.encode(
        {
          "email": email,
        },
      ),
    );

    if (res.body == null || res.body.isEmpty) {
      return;
    }

    final resData = json.decode(res.body);
    
   	//  若在後端沒有錯誤的話, 我們在此處拿到的res.body, 應該就會對應 Create Handler的:
	//	c.JSON(http.StatusOK, gin.H{
	//		"id": result.InsertedID,
	//	})

    return resData["id"];
  }
  
  
  
  /// 更新在後端的數據
  Future<dynamic> updateTestDocument(String id, String email) async {
    if (id == null || id.isEmpty) {
      if (insertedID == null) {
        return null;
      } else {
        id = insertedID;
      }
    }

    final String url = "http://192.168.1.135:8080/test/$id";
    final res = await http.put(
      url,
      headers: {
        HttpHeaders.contentTypeHeader: "application/json",
      },
      body: json.encode(
        {
          "email": email,
        },
      ),
    );

    if (res.body == null || res.body.isEmpty) {
      return;
    }

    final resData = json.decode(res.body);
    
    return resData;
  }
  

```

#### 收集數據的UI

```dart
	// 用來輸入創建或是更新的Email
 	 Padding(
            padding: const EdgeInsets.all(8.0),
            child: TextField(
              decoration: InputDecoration(
                hintText: "Email",
                hintStyle: TextStyle(
                  fontSize: 16,
                ),
              ),
              controller: emailController,
            ),
          ),
	  // 用來輸入需要跟新的Document Id
          Padding(
            padding: const EdgeInsets.all(8.0),
            child: TextField( 
              decoration: InputDecoration(
                hintText: "ID",
                hintStyle: TextStyle(
                  fontSize: 16,
                ),
              ),
              controller: idController,
            ),
          ),
	  // 按下即請求後端增加一個Test Document
          FlatButton(
            child: Text("Create"),
            onPressed: () async {
              print("Create Pressed");
              await Provider.of<WebScoketService>(context, listen: false)
                  .addTestDocument(emailController.text);
            },
          ),
	  // 按下即請求後端Update特定的Test Docuemnt
          FlatButton(
            child: Text("Update"),
            onPressed: () {
              print("update Pressed");
              Provider.of<WebScoketService>(context, listen: false)
                  .updateTestDocument(idController.text, emailController.text);
            },
          ),
```
我們用這幾個 UI 來創建和更新後端的 Test Document


### 數據監控

#### 後端 

後端要提供一個 GET 接口, 這個接口可以被 Upgrade 成 Websocket 接手。當這個 WebSocket 接口開啟時, 隨即監聽 MongoDB 中的 **ChangeStream**。

[Change Stream](https://docs.mongodb.com/manual/changeStreams/)
[Change Events](https://docs.mongodb.com/manual/reference/change-events/)

```go
router.GET("/test/subscribe/:id", func(c *gin.Context) {

	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil) // 將此 GET REQUEST 升級成 WebSocket

	defer ws.Close()

	if err != nil {
		errStr := fmt.Sprintf("The websocket is not working due to the error: %+v \n", err)	
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errStr,
		})
		return
	}

	id := c.Param("id") // 取出要監聽的Document ID

	oid, err := primitive.ObjectIDFromHex(id) 

	if err != nil {
		errStr := fmt.Sprintf("The given id cannot be transform to oid: %+v \n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": errStr,
		})
		return
	}

	pipeline := mongo.Pipeline{bson.D{{"$match", bson.D{{"fullDocument._id", oid}}}}} // 選擇監聽Output的條件

	collectionStream, err := database.DB.Collection("test").Watch(context.TODO(), pipeline, 
		options.ChangeStream().SetFullDocument(options.UpdateLookup)) // 拿到監聽的 Stream


	if err != nil {
		errStr := fmt.Sprintf("Cannot get the stream: %+v \n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errStr,
		})
		return
	}

	defer collectionStream.Close(context.TODO())

	for {   // 開始監聽MongoDB特定ID的Document是否有被Changed
		ok := collectionStream.Next(context.TODO()) // 若有新的Change Event
		if ok {
			next := collectionStream.Current // 取出現在的Change Event

			log.Printf("Next: %+v", next)
				
			// 使用WebSocket 將Change Event 的資料傳送至前端
			err = ws.WriteMessage(websocket.TextMessage, []byte(next.String())) 

			if err != nil {
				break
			}
		}
	}
})

```




# 如何取得你的IP
當我們需要用實體機來測試時, 我們可以透過Wifi實現前後端的對連, 此時我們要將封包由前端發向後端, 則需要Wifi下區網內的IP.

## Windows
再Command Prompt中輸入 `ipconfig`, 紅框內則為所需的IP Address
![](https://github.com/sufferingCoders/QuenC-Backend/blob/master/github_images/windows_ip.jpg?raw=true)

## MacOS
偏好設定 > Network > Advanced > TCP/IP 中就有
![](https://github.com/sufferingCoders/QuenC-Backend/blob/master/github_images/MacOS_IP.png?raw=true)



# 滿滿的坑

## ( 讀取HTML Templates 時 ) invalid memory address or nil pointer dereference

當跑

```go
	c.HTML(http.StatusOK, "EmailVerificationSuccessful.tmpl", gin.H{
		"email": user.Email,
	})
```

的時候出現 `invalid memory address or nil pointer dereference` 的訊息

問題是在於:

沒有導入HTML

在router/index.go中
```go
	router := gin.Default()
	router.LoadHTMLGlob("templates/*") // 由此項導入 HTML
```


## response.Write on hijacked connection from github.com/gin-gonic/gin.(*responseWriter).Write
這是因為我們已經將 Get Request 升級為 Websocket, 不能再使用c.JSON去傳送HTTP訊息



