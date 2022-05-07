package internal

import (
	"bytes"
	"cdp_agent/common"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	URL string = "/cdp/check"
)

type Resp struct {
	Code string `json:"code"`
	Msg  string `json:"message"`
}

type Auth struct {
	Ipaddr   string `json:"ip"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func verify(writer http.ResponseWriter, request *http.Request) {
	fmt.Printf("===================================================")
	decoder := json.NewDecoder(request.Body)
	fmt.Println(decoder)
	var auth Auth
	if err := json.NewDecoder(request.Body).Decode(&auth); err != nil {
		_ = request.Body.Close()
		common.Error.Println(err)
		fmt.Println(err)
	}
	fmt.Println("接收：", auth)
	common.Info.Println("接收：", auth)

	data := 0
	conf := common.Conf{}
	parmer := conf.GetModelClass("config.json")
	fmt.Println("config :", parmer)
	if !parmer.IsMock {
		data = getData(auth)
		fmt.Println("验证连接-----> :", data)
		common.Info.Println("验证连接-----> :", data)
	} else {
		data = 200
	}

	var result Resp
	if data != 200 {
		result.Code = "0"
		result.Msg = "验证失败"
		common.Info.Println("验证失败")
		fmt.Println("验证失败")
	} else {
		result.Code = "1"
		result.Msg = "验证成功"
		common.Info.Println("验证成功")
		fmt.Println("验证成功")
	}

	if err := json.NewEncoder(writer).Encode(result); err != nil {
		log.Fatal(err)
	}

}

func getData(auth Auth) int {
	fmt.Println("当前验证ip：", auth.Ipaddr)
	common.Info.Println("当前验证ip：", auth.Ipaddr)
	URL := UrltoEncode(fmt.Sprintf("http://%s:%d//api/login", auth.Ipaddr, auth.Port))
	contentType := "application/json"
	str := fmt.Sprintf(`{"username":"%s","password":"%s"}`, auth.Username, auth.Password)
	key := []byte("silvanware123456")
	result := common.AesEncryptECB([]byte(str), key)

	data := make(map[string]string)
	data["data"] = base64.StdEncoding.EncodeToString(result)
	bytesData, _ := json.Marshal(data)
	response, err := http.Post(URL, contentType, bytes.NewBuffer(bytesData))
	fmt.Println(response)
	defer response.Body.Close()
	if err != nil {
		fmt.Println(err)
		common.Error.Println(err)
		return 0
	}
	return response.StatusCode
}

func Service() {
	http.HandleFunc(URL, verify)
	conf := common.Conf{}
	parmer := conf.GetModelClass("config.json")
	port := parmer.ServicePort
	fmt.Println("资产验证连接服务端口：", port)                          //监听端口输出
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil) //设置监听服务
	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	}

}
