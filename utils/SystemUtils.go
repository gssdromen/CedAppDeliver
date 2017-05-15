package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
)

// RunShellCommand 运行Shell命令
func RunShellCommand(cmd *exec.Cmd) (string, error) {
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	outString := out.String()
	return outString, nil
}

// SendTextMessageToDingDing 给钉钉发消息
func SendTextMessageToDingDing(accessToken string, content string, ats []string) {
	var url = "https://oapi.dingtalk.com/robot/send?access_token=" + accessToken

	var atObj = make(map[string]interface{})
	atObj["isAtAll"] = false
	atObj["atMobiles"] = ats

	var cententObj = map[string]interface{}{"content": content}
	var strObj = map[string]interface{}{"msgtype": "text", "text": cententObj, "at": atObj}
	str, _ := json.Marshal(strObj)
	var jsonStr = []byte(str)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
