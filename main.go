package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"k8s.io/client-go/tools/remotecommand"
	"log"
)

var (
	websocket = flag.Bool("websocket", true, "enable/disable websocket protocol")
)

func init() {
	flag.Parse()
}

//func main() {
//	opts := sockjs.DefaultOptions
//	opts.Websocket = *websocket
//	handler := sockjs.NewHandler("/api/sockjs", opts, echoHandler)
//	http.Handle("/api/sockjs/", handler)
//	http.Handle("/", http.FileServer(http.Dir("web/")))
//	log.Println("Server started on port: 8090")
//	log.Fatal(http.ListenAndServe(":8090", nil))
//}

func echoHandler(session sockjs.Session) {
	log.Println("new sockjs session established?")
	//for {
	//	if msg, err := session.Recv(); err == nil {
	//		if err := session.Send(msg); err != nil {
	//			break
	//		}
	//		continue
	//	} else {
	//		fmt.Println(err)
	//	}
	//	break
	//}
	go consumer(session)
	log.Println("sockjs session closed?")
}

func consumer(session sockjs.Session) {
	for {
		if msg, err := session.Recv(); err == nil {
			if err := session.Send(msg); err != nil {
				break
			}
			continue
		} else {
			fmt.Println(err)
		}
		break
	}
}

func main() {
	K8s.Init()
	r := gin.Default()
	r.StaticFS("/web", gin.Dir("./web", false))
	r.GET("/api/sockjs/*path", gin.WrapH(CreateAttachHandler("/api/sockjs")))
	r.GET("/api/v1/pod/:namespace/:podName/shell/:containerName", func(c *gin.Context) {
		sessionID, err := genTerminalSessionId()
		if err != nil {
			c.JSON(500, gin.H{
				"msg":  "err",
				"code": 500,
			})
		}
		terminalSessions.Set(sessionID, TerminalSession{
			id:       sessionID,
			bound:    make(chan error),
			sizeChan: make(chan remotecommand.TerminalSize),
		})
		// {"Op":"bind","SessionID":"db1888b4dd29e3c61540c56a5f7cfc22"}
		// {"Op":"stdin","Data":"ls\r","Cols":164,"Rows":41}
		go WaitForTerminal(c, sessionID)
		c.JSON(200, gin.H{
			"msg":  "successful",
			"code": 200,
			"data": gin.H{
				"id": sessionID,
			},
		})
	})
	r.Run("0.0.0.0:8090")
}
