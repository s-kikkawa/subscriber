package database

import (
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
    DBMS = "mysql"
    USER = "test_user" // DBのユーザ名
    PASS = "test_pass" // DBのパスワード
    HOST = "localhost" // DBのホスト名もしくはIP
    PORT = "3306"      // DBのポート
    DBNAME = "exampledb" // DB名
    CONNECT = USER + ":" + PASS + "@tcp(" + HOST + ":" + PORT + ")/" + DBNAME + "?parseTime=true"
)

var db *gorm.DB

type Replica struct {
    gorm.Model
    OriginalId int
    ItemCode string
    Text   string
}

// DB初期化
func Init() {
    var err error
    db, err = gorm.Open(DBMS, CONNECT)
    if err != nil {
        panic("データベースの接続に失敗しました")
    }
    db.AutoMigrate(&Replica{})
}


// 登録
func Insert(originalId int, itemCode string, text string) {
    db.Create(&Replica{OriginalId: originalId, ItemCode: itemCode, Text: text})
}

// 更新
func Update(original_id int, itemCode string, text string) {
    id := SelectIdByOriginalId(original_id)
    var replica Replica
    db.First(&replica, id)
    replica.ItemCode = itemCode
    replica.Text = text
    db.Save(&replica)
}

// 削除
func Delete(original_id int) {
    id := SelectIdByOriginalId(original_id)
    var replica Replica
    db.First(&replica, id)
    db.Delete(&replica)
}

// 全取得
func SelectAll() []Replica {
    var replicas []Replica
    db.Order("created_at desc").Find(&replicas)
    return replicas
}

// 1行取得
func SelectRow(id int) Replica {
    var replica Replica
    db.First(&replica, id)
    return replica
}

// コピー元のIDを元にレプリカのidを検索する
func SelectIdByOriginalId(original_id int) uint {
    var replica Replica
    db.Where("original_id = ?", original_id).First(&replica)
    return replica.ID
}

// DBクローズ
func Close() {
    db.Close()
}
