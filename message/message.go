package message

import (
    "log"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/sqs"
    "github.com/aws/aws-sdk-go/aws/session"
    "encoding/json"
    "fmt"
)

const (
    DO_NOTHING = false // メッセージ受信しない場合はtrueにしてください
    AWS_REGION = "ap-northeast-1"
    QUEUE_URL  = "https://sqs.ap-northeast-1.amazonaws.com/123456789012/TestQueue"
)

type Message struct {
    OperationType string `json:"operationType"`
    OriginalId string `json:"id"`
    ItemCode string `json:"itemCode"`
    Text   string   `json:"text"`
}

// SQSからメッセージを受信します
func ReceiveSqsMessages() []Message {
    var messages []Message
    if DO_NOTHING{
        log.Print("DO_NOTHING = true のためメッセージは受信しません")
        return messages
    }
    sess := session.Must(session.NewSession())
    svc := sqs.New(sess, aws.NewConfig().WithRegion(AWS_REGION))

    params := &sqs.ReceiveMessageInput{
        QueueUrl: aws.String(QUEUE_URL),
        // 1度に取得できるメッセージの数
        MaxNumberOfMessages: aws.Int64(10),
        // メッセージがない場合の待機時間（秒）
        WaitTimeSeconds: aws.Int64(20),
    }

    resp, err := svc.ReceiveMessage(params)
    log.Print(err)
    log.Printf("messages count: %d", len(resp.Messages))
    
    if len(resp.Messages) == 0 {
        log.Print("empty queue.")
    } else {
        for _, m := range resp.Messages {
            fmt.Printf("Found target message:\n%s\n", m.GoString())
            body := *m.Body
            log.Print(body)

            var jsonline Message
            if err := json.Unmarshal([]byte(body), &jsonline); err != nil {
                log.Fatal(err)
            }
            messages = append(messages, jsonline)
            fmt.Printf("%s : %s\n", jsonline.OperationType, jsonline.ItemCode)

            deleteParams := &sqs.DeleteMessageInput{
                QueueUrl:      aws.String(QUEUE_URL),
                ReceiptHandle: aws.String(*m.ReceiptHandle),
            }
            _, err := svc.DeleteMessage(deleteParams)

            if err != nil {
                log.Print(err)
            }
        }
    }

    return messages
}

// json形式でメッセージ本体を作成します
func createMessage(operationType string, itemCode string, text string) string {
    message := new(Message)
    message.OperationType = operationType
    message.ItemCode = itemCode
    message.Text = text
    message_json, _ := json.Marshal(message)
    return string(message_json)
}
