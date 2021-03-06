package internal

import (
	"bytes"
	"cdp_agent/common"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type collect struct {
	mysqlClinet     *common.MysqlClient
	influxClient    *common.InfluxDBClient
	entityList      *[]common.Endpoint
	entity          *common.Endpoint
	kafka           *common.Kafka
	startTime       string
	endTime         string
	conf            *common.AutoGenerated
	alertMetricDict map[string]string
}

func (c *collect) init() {
	client := common.Conf{}
	conf := client.GetModelClass("config.json")
	c.conf = conf
	c.mysqlClinet = &common.MysqlClient{
		Host:     conf.Mysql.Host,
		Port:     conf.Mysql.Port,
		Dbname:   conf.Mysql.Dbname,
		Username: conf.Mysql.Username,
		Password: conf.Mysql.Password,
	}
	c.influxClient = &common.InfluxDBClient{
		Cli:      nil,
		Addr:     conf.Influxdb.Host,
		UserName: conf.Influxdb.Username,
		Password: conf.Influxdb.Password,
		Port:     conf.Influxdb.Port,
		DB:       conf.Influxdb.Database,
	}
	c.kafka = &common.Kafka{
		Brokers: conf.Kafka.Brokers,
		Topic:   conf.Kafka.Topic,
	}
	m, _ := time.ParseDuration(fmt.Sprintf("-%dm", conf.IntervalTime))
	c.startTime = time.Now().Add(m).Format("2006-01-02 15:04:05")
	c.endTime = time.Now().Format("2006-01-02 15:04:05")
	c.alertMetricDict = map[string]string{
		"state":   "ent.job.state",
		"offline": "silvanware.conn.status",
	}
}

func Run() {
	c := collect{}
	c.init()
	if c.conf.IsMock {
		c.entityList = &c.conf.Endpoint
	} else {
		c.getEntity()
	}
	for _, entity := range *c.entityList {
		c.entity = &entity
		var metrics []*common.Metric
		if !c.conf.IsMock {
			metrics = c.toCollect()
		} else {
			metrics = MockData(c.alertMetricDict)
		}

		//采集的数据不为空 才进行发送数据至kafka
		if metrics != nil {
			c.sendDataToKafka(metrics)
		}
	}
}

// 采集数据
func (c *collect) toCollect() []*common.Metric {
	// 注册登录 获取 cookie
	cookie, e := c.getCookieInfo()
	if e != nil || cookie == "" {
		fmt.Printf("Get cookie info error :%s \n", e)
		common.Error.Printf("Get cookie info error :%s \n", e)
		return nil
	}
	fmt.Printf("CDP登录获取cookie：%s\n", cookie)
	// 获取服务ip和任务id
	hostInfo, e := c.getHostInfo(cookie)
	if e != nil || hostInfo == nil {
		fmt.Printf("Get host info error :%s \n", e)
		common.Error.Printf("Get host info error :%s \n", e)
		return nil
	}
	fmt.Printf("CDP获取hostinfo：%s \n", hostInfo)
	// 获取disk信息
	diskInfo, e := c.getDiskInfo(hostInfo, cookie)
	if e != nil || diskInfo == nil {
		fmt.Printf("Get disk info error :%s \n", e)
		common.Error.Printf("Get disk info error :%s \n", e)
	}
	fmt.Printf("CDP获取diskinfo：%s \n", diskInfo)

	// 数据添加
	metrics := make([]*common.Metric, 0)
	metrics = append(metrics, parseHostInfo(hostInfo, c.alertMetricDict)...)
	metrics = append(metrics, parseDiskInfo(diskInfo, c.alertMetricDict)...)
	return metrics
}

/**
获取资产信息
*/
func (c *collect) getEntity() {
	c.mysqlClinet.GetConn()
	defer c.mysqlClinet.CloesConn()

	sql := fmt.Sprintf(`SELECT e.EntityID AS entityId, 
			max( CASE props.propertyKey WHEN 'ent.monitorstatus' THEN prop.PropertyValue ELSE NULL END ) AS monitorstatus, 
			max( CASE props.propertyKey WHEN 'ent.monitorip' THEN prop.PropertyValue ELSE NULL END ) AS ip, 
			max( CASE props.propertyKey WHEN 'ent.monitor.port' THEN prop.PropertyValue ELSE NULL END ) AS port, 
			max( CASE props.propertyKey WHEN 'ent.username' THEN prop.PropertyValue ELSE NULL END ) AS username, 
			max( CASE props.propertyKey WHEN 'ent.password' THEN prop.PropertyValue ELSE NULL END ) AS password 
			FROM cp_ci_entity e  
			JOIN cp_ci_model m ON e.modelID = m.ModelID  
			AND e.isdeleted = 0  
			AND m.MonitorType = 3486784401  
			AND m.BussType = 1 
			JOIN cp_ci_entity_prop prop ON e.EntityID = prop.EntityID 
			JOIN cp_ci_props props ON prop.PropertyID = props.PropertyID  
			GROUP BY e.EntityID  
			HAVING monitorstatus = 1`)
	result := c.mysqlClinet.Query(sql)
	if len(result) <= 0 {
		common.Info.Println("get entity is null")
		fmt.Println("get entity is null")
	}
	entityList := make([]common.Endpoint, 0)
	for e := range result {
		entityMap := result[e]
		entity := common.Endpoint{}
		entity.IPAddress = entityMap["ip"].(string)
		entity.Port, _ = strconv.Atoi(entityMap["port"].(string))
		entity.Username = entityMap["username"].(string)
		entity.Password = entityMap["password"].(string)
		entityId, _ := strconv.Atoi(entityMap["entityId"].(string))
		entity.EntityID = entityId
		entityList = append(entityList, entity)
	}
	c.entityList = &entityList
}

func (c *collect) sendDataToKafka(metrics []*common.Metric) {
	if len(metrics) != 0 {
		kafkaClient := common.KafkaClient{
			Brokers: strings.Split(c.kafka.Brokers, ","),
			Topic:   c.kafka.Topic,
		}
		if err := kafkaClient.Connect(); err != nil {
			fmt.Println(err)
			common.Error.Println(err)
			panic(err)
		}
		defer kafkaClient.Close()

		timeNow := time.Now().UnixNano() / 1e6
		timeUnixStr := strconv.FormatInt(timeNow, 10)
		// 拆分数据 按20条进行拆分
		data := cutOffList(metrics, 20)
		for i := range data {
			jsonData, err := json.Marshal(data[i])
			if err != nil {
				common.Error.Println(err)
				fmt.Println(err)
				panic(err)
				continue
			}
			icopPkg := common.IcopPackage{}
			icopPkg.SetCollectionKey(fmt.Sprintf("%d_%s", c.entity.EntityID, c.entity.IPAddress))
			icopPkg.SetNodeKey(c.entity.IPAddress)
			icopPkg.SetNodeSite("")
			icopPkg.SetCommandType("9100")
			icopPkg.SetCommandCode("")
			icopPkg.SetIsUpSend(true)
			icopPkg.SetData(string(jsonData))
			icopPkg.SetRecordTime(timeUnixStr)
			//kafka send data
			kafkaClient.Send(icopPkg.EncodePackage())
			fmt.Println(icopPkg.EncodePackage())
			common.Info.Println(icopPkg.EncodePackage())
		}
	}
	fmt.Printf("本次采集结束，共采集发送kafka数据 %d 条\n", len(metrics))
}

func (c *collect) getDiskInfo(hostInfo map[string]interface{}, cookie string) ([]map[string]interface{}, error) {
	diskInfo := make([]map[string]interface{}, 0)
	URL := UrltoEncode(fmt.Sprintf("http://%s:%d/api/disasterTask/list_hosts", c.entity.IPAddress, c.entity.Port))
	for _, serverDatas := range hostInfo["server_data"].([]interface{}) {
		server := serverDatas.(map[string]interface{})
		for _, hosts := range server["Hosts"].([]interface{}) {
			host := hosts.(map[string]interface{})
			data := make(map[string]interface{})
			data["AgentId"] = host["AgentId"].(float64)
			data["serverIp"] = host["server_ip"].(string)
			bytesData, err := json.Marshal(data)
			if err != nil {
				fmt.Println(err)
				common.Error.Println(err)
				return nil, err
			}
			u, _ := url.ParseRequestURI(URL)
			urlStr := u.String()
			httpClient := &http.Client{}
			request, err := http.NewRequest("POST", urlStr, bytes.NewReader(bytesData))
			if err != nil {
				fmt.Println(err)
				common.Error.Println(err)
				return nil, err
			}

			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Cookie", cookie)
			response, err := httpClient.Do(request)
			if err != nil {
				fmt.Println(err)
				common.Error.Println(err)
				return nil, err
			}
			if response.StatusCode != 200 {
				common.Error.Printf("get diskinfo error, response code: %d \n", response.StatusCode)
				fmt.Printf("get diskinfo error, response code: %d \n", response.StatusCode)
				return nil, nil
			}

			body, _ := ioutil.ReadAll(response.Body)
			a := make(map[string]interface{})
			_ = json.Unmarshal(body, &a)
			if len(a) != 0 {
				bytes, _ := json.Marshal(a)
				saveToFile(fmt.Sprintf("%d_DiskInfo", c.entity.EntityID), string(bytes))
			}
			diskInfo = append(diskInfo, a)
			// 结束在关闭
			request.Body.Close()
			response.Body.Close()
		}
	}
	return diskInfo, nil
}

func (c *collect) getHostInfo(cookie string) (map[string]interface{}, error) {
	URL := UrltoEncode(fmt.Sprintf("http://%s:%d/api/hostInfos/all_client_list", c.entity.IPAddress, c.entity.Port))
	data := make(map[string]string)
	data["key"] = ""
	bytesData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		common.Error.Println(err)
		return nil, err
	}
	u, _ := url.ParseRequestURI(URL)
	urlStr := u.String()
	httpClient := &http.Client{}
	request, err := http.NewRequest("POST", urlStr, bytes.NewReader(bytesData))
	defer request.Body.Close()
	if err != nil {
		fmt.Println(err)
		common.Error.Println(err)
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Cookie", cookie)
	response, err := httpClient.Do(request)
	if err != nil {
		fmt.Println(err)
		common.Error.Println(err)
		return nil, err
	}
	if response.StatusCode != 200 {
		common.Error.Printf("get hostinfo error, response code:%d \n", response.StatusCode)
		fmt.Printf("get hostinfo error, response code:%d \n", response.StatusCode)
		return nil, nil
	}
	body, _ := ioutil.ReadAll(response.Body)
	a := make(map[string]interface{})
	_ = json.Unmarshal(body, &a)
	if len(a) != 0 {
		bytes, _ := json.Marshal(a)
		saveToFile(fmt.Sprintf("%d_HostInfo", c.entity.EntityID), string(bytes))
	}
	return a, nil
}

func (c *collect) getCookieInfo() (string, error) {
	url := UrltoEncode(fmt.Sprintf("http://%s:%d/api/login", c.entity.IPAddress, c.entity.Port))
	contentType := "application/json"
	str := fmt.Sprintf(`{"username":"%s","password":"%s"}`, c.entity.Username, c.entity.Password)
	key := []byte("silvanware123456")
	result := common.AesEncryptECB([]byte(str), key)

	data := make(map[string]string)
	data["data"] = base64.StdEncoding.EncodeToString(result)
	bytesData, _ := json.Marshal(data)
	response, err := http.Post(url, contentType, bytes.NewBuffer(bytesData))
	if err != nil {
		fmt.Println(err)
		common.Error.Println(err)
		return "", err
	}
	defer response.Body.Close()
	fmt.Println(response)

	if response.StatusCode != 200 {
		common.Error.Printf("get cookie error, response code : %d \n", response.StatusCode)
		fmt.Printf("get cookie error, response code : %d \n", response.StatusCode)
		return "", nil
	}

	var cookie string
	cookies := response.Cookies()
	for i := 0; i < len(cookies); i++ {
		cookie = cookie + cookies[i].Name + "=" + cookies[i].Value
		if i < len(cookies)-1 {
			cookie = cookie + ";"
		}
	}
	common.Info.Printf("cookie : %s \n", cookie)
	return cookie, nil
}

// 保存版本号
func (c *collect) saveVersion(version string) {
	c.mysqlClinet.GetConn()
	defer c.mysqlClinet.CloesConn()

	sql := fmt.Sprintf(`select PropertyValue from cp_ci_entity_prop where EntityID = %d AND PropertyID = %d`, c.entity.EntityID, 116)
	result := c.mysqlClinet.Query(sql)
	// 长度大于0 则 做更新操作 否则做新增
	if len(result) > 0 {
		sql = fmt.Sprintf(`update cp_ci_entity_prop set PropertyValue = '%s' where EntityID = %d and PropertyID = %d and IsDeleted =0`, version, c.entity.EntityID, 116)
	} else {
		sql = fmt.Sprintf(`INSERT INTO cp_ci_entity_prop (EntityID, PropertyID, PropertyValue, CreateTime, CreateUser, UpdateUser, UpdateTime, IsDeleted, MonitorPropID) VALUES (%d, %d, "%s", now(), 0, 0, NULL, 0, NULL)`, c.entity.EntityID, 116, version)
	}
	c.mysqlClinet.Exec(sql)
	common.Info.Printf("版本号更新成功！ version:%s \n", version)
	fmt.Printf("版本号更新成功！ version:%s \n", version)
}

// 保存到文件
func saveToFile(fileName string, info string) {
	pwd := fmt.Sprintf("./data/%s.data", fileName)
	path := strings.Split(pwd, ".")[0]
	_, e := os.Stat(path)
	if e != nil {
		_ = os.Mkdir(path, os.ModePerm)
	}
	file, e := os.OpenFile("./data/"+fileName+".data", os.O_CREATE|os.O_WRONLY, 0600)
	if e != nil {
		common.Error.Println(fmt.Sprintf("Save To File :%s Error！", fileName))
		return
	}
	defer file.Close()

	_, e = file.WriteString(info)
	if e != nil {
		common.Error.Println(fmt.Sprintf("Save To File :%s Error！", fileName))
		return
	}
}

func UrltoEncode(url string) string {
	new1 := strings.Replace(url, " ", "%20", -1)
	//new2 := strings.Replace(new1, "&", "\\&", -1)
	return new1
}

//strList 要切分的数组
//listSize 切分后每个数组的size
func cutOffList(strList []*common.Metric, listSize int) (ss [][]*common.Metric) {
	// 对listSize取模
	mod := len(strList) % listSize
	// 对listSize取余
	k := len(strList) / listSize

	// 计算循环的截止数
	var end int
	if mod == 0 {
		end = k
	} else {
		end = k + 1
	}

	for i := 0; i < end; i++ {
		if i != k {
			ss = append(ss, strList[i*listSize:(i+1)*listSize])
		} else {
			ss = append(ss, strList[i*listSize:])
		}
	}
	return
}
