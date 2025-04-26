package lark

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

var Client *lark.Client

func init() {
	// 创建 Client
	Client = lark.NewClient("cli_a7025ca1b3f4102f", "2s8CPygb1ulimnZKJlPZkb27s7KCRE5M")

}

func SendText(openId string, Uuid string, text string) error {

	// 构造 JSON 内容
	contentStruct := map[string]string{
		"text": text,
	}
	contentJSON, err := json.Marshal(contentStruct)
	if err != nil {
		return fmt.Errorf("failed to marshal content: %w", err)
	}

	// 创建请求对象
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(`open_id`).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(openId).
			MsgType(`text`).
			Content(string(contentJSON)).
			Uuid(Uuid).
			Build()).
		Build()

	// 发起请求
	resp, err := Client.Im.Message.Create(context.Background(), req)

	// 处理错误
	if err != nil {
		return err
	}

	// 服务端错误处理
	if !resp.Success() {
		err := fmt.Errorf("logId: %s, error response: \n%s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
		return err
	}

	// 业务处理
	fmt.Println(larkcore.Prettify(resp))

	return nil
}
