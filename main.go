package main

import (
    "subscriber/database"
    "subscriber/message"
    "strconv"
    "net/http"
    "github.com/gin-gonic/gin"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    "os"
    "os/signal"
    "log"
)

const (
    PORT = "8080"
)

func main() {

    go func(){
        router := gin.Default()
        router.LoadHTMLGlob("templates/*.html")
        database.Init()

        // 一覧画面表示
        router.GET("/", func(ctx *gin.Context) {
            replicas := database.SelectAll()
            ctx.HTML(http.StatusOK, "index.html", gin.H{
                "replicas": replicas,
            })
        })
        // メッセージを取得してDBを更新する
        router.POST("/sync", func(ctx *gin.Context) {
            messages := message.ReceiveSqsMessages()
            for _, message := range messages{
                var replica database.Replica
                replica.OriginalId = convertId(message.OriginalId)
                replica.ItemCode = message.ItemCode
                replica.Text = message.Text

                switch(message.OperationType){
                case "INSERT" :
                    database.Insert(replica.OriginalId, replica.ItemCode, replica.Text)
                case "DELETE" :
                    database.Delete(replica.OriginalId)
                case "UPDATE" :
                    database.Update(replica.OriginalId, replica.ItemCode, replica.Text)
                }
            }
            ctx.Redirect(http.StatusFound, "/")
        })

        router.Run(":" + PORT)
    }()

    // Ctrl + c で停止した際にDBをクローズする
    quit := make(chan os.Signal)
    signal.Notify(quit, os.Interrupt)
    <-quit
    database.Close()
}

// IDを数値に変換する
func convertId(idStr string) int {
    id, err := strconv.Atoi(idStr)
    if err != nil {
        log.Fatal("IDの変換に失敗しました")
    }
    return id
}


