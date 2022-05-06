package common

import (
	"fmt"
	"strconv"
	"time"

	"github.com/influxdata/influxdb1-client/v2"
)

const (
	MyDB          = "test"
	username      = "admin"
	password      = ""
	MyMeasurement = "cpu_usage"
)

func main() {

	influxClient := InfluxDBClient{
		Addr:     "192.168.1.164",
		UserName: "admin",
		Password: "123456",
		DB:       "icoptest",
	}
	influxClient.Conn()
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  influxClient.DB,
		Precision: "s",
	})
	if err!= nil{fmt.Println("error")}
	point, _ := client.NewPoint("testIns", map[string]string{"aa": "aa"}, map[string]interface{}{"bb": "bb"}, time.Now())
	bp.AddPoint(point)

	//influxClient.WritesPoints()
	//insert
	//influxClient.WritesPoints()

	////获取10条数据并展示
	//qs := fmt.Sprintf("SELECT * FROM %s LIMIT %d", MyMeasurement, 10)
	//res, err := influxClient.QueryDB(qs)
	//if err != nil {
	//	Error.Println(err)
	//}
	//
	//for i, row := range res[0].Series[0].Values {
	//	t, err := time.Parse(time.RFC3339, row[0].(string))
	//	if err != nil {
	//		Error.Println(err)
	//	}
	//	//fmt.Println(reflect.TypeOf(row[1]))
	//	valu := row[2].(json.Number)
	//	Info.Printf("[%2d] %s: %s\n", i, t.Format(time.Stamp), valu)
	//}
}

type InfluxDBClient struct {
	Cli 		client.Client
	Addr 		string
	UserName 	string
	Password 	string
	Port        int
	DB 			string
}

func (i *InfluxDBClient) Conn() {
	port 	 := strconv.Itoa(i.Port)
	cli, err := client.NewHTTPClient(client.HTTPConfig{
		Addr	:     "http://"+i.Addr+":"+ port,
		Username: 	  i.UserName,
		Password:  	  i.Password,
	})
	if err != nil {
		Error.Println(err)
	}
	_,_, err = cli.Ping(5)
	if err != nil {
		Error.Println(err)
	}
	i.Cli = cli

}

//query
func (i *InfluxDBClient) QueryDB( cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: i.DB,
	}
	if response, err := i.Cli.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

//Insert
func (i *InfluxDBClient) WritesPoints(points []*client.Point) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  i.DB,
		Precision: "s",
	})
	if err != nil {
		Error.Println(err)
	}
	//pt, err := client.NewPoint(
	//	measurement,
	//	tags,
	//	fields,
	//	time.Now(),
	//)
	//if err != nil {
	//	Error.Println(err)
	//}
	//bp.AddPoint(pt)
	if len(points) != 0{
		bp.AddPoints(points)
	}

	if err := i.Cli.Write(bp); err != nil {
		Error.Println(err)
	}
}

func (i *InfluxDBClient) CloseConn() {
	_ = i.Cli.Close()
}