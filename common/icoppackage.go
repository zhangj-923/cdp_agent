package common

import (
	"fmt"
)

type IcopPackage struct {
	CollectionKey string
	NodeKey       string
	NodeSite      string
	RecordTime    string
	CommandType   string
	CommandCode   string
	CmdSubCode    string
	IsUpSend      bool
	Data          string
}

type Metric struct {
	DType       string      `json:"dType"`
	Metric      string      `json:"metric"`
	Table       string      `json:"table"`
	Tags        interface{} `json:"tags,omitempty"`
	Value       interface{} `json:"value"`
	AlertMetric string      `json:"alertMetric"`
}

func Format(result, tags map[string]interface{}, measurement string, alertMetricDict map[string]string) ([]*Metric, error) {
	var data []*Metric

	for k, v := range result {
		m_ := Metric{}
		if v2, ok := v.(int); ok {
			m_.DType = "int"
			m_.Value = v2
		} else if v2, ok := v.(int32); ok {
			m_.DType = "int"
			m_.Value = v2
		} else if v2, ok := v.(int64); ok {
			m_.DType = "int"
			m_.Value = v2
		} else if v2, ok := v.(float32); ok {
			m_.DType = "double"
			m_.Value = v2
		} else if v2, ok := v.(float64); ok {
			m_.DType = "double"
			m_.Value = v2
		} else {
			m_.DType = "string"
			m_.Value = v
		}
		for metricName, alertAttribute := range alertMetricDict {
			if k != metricName {
				continue
			}
			m_.AlertMetric = alertAttribute
		}

		m_.Metric = k
		m_.Table = measurement
		m_.Tags = tags

		if len(tags) == 0 {
			m_.Tags = nil
		}
		data = append(data, &m_)
	}

	//JsonData, err := json.Marshal(data)

	//if err != nil {
	//	return nil, err
	//}
	return data, nil
}

func (pkg *IcopPackage) EncodePackage() string {
	var (
		full_node_key string
		dirct         int
		command_pre   string
		cmdKey        string
	)

	if pkg.NodeSite == "" {
		full_node_key = pkg.NodeKey
	} else {
		full_node_key = pkg.NodeKey + ":" + pkg.NodeSite
	}

	if full_node_key == "" {
		full_node_key = "null"
	}

	if pkg.IsUpSend {
		dirct = 1
	} else {
		dirct = 0
	}

	if pkg.CmdSubCode == "" {
		command_pre = pkg.CommandType
	} else {
		command_pre = pkg.CommandType + "-" + pkg.CmdSubCode
	}

	if pkg.CommandCode == "" {
		cmdKey = command_pre
	} else {
		cmdKey = command_pre + ":" + pkg.CommandCode
	}

	return fmt.Sprintf("^^%s^%d%s^%s^%s^%s$$\\n", full_node_key, dirct, cmdKey, pkg.CollectionKey, pkg.RecordTime, pkg.Data)

}

func (pkg *IcopPackage) GetCollectionKey() string {
	return pkg.CollectionKey
}

func (pkg *IcopPackage) SetCollectionKey(CollectionKey string) {
	pkg.CollectionKey = CollectionKey
}

func (pkg *IcopPackage) GetNodeSite() string {
	return pkg.NodeSite
}

func (pkg *IcopPackage) SetNodeSite(NodeSite string) {
	pkg.NodeSite = NodeSite
}

func (pkg *IcopPackage) GetNodeKey() string {
	return pkg.NodeKey
}

func (pkg *IcopPackage) SetNodeKey(NodeKey string) {
	pkg.NodeKey = NodeKey
}

func (pkg *IcopPackage) GetRecordTime() string {
	return pkg.RecordTime
}

func (pkg *IcopPackage) SetRecordTime(RecordTime string) {
	pkg.RecordTime = RecordTime
}

func (pkg *IcopPackage) GetCommandType() string {
	return pkg.CommandType
}

func (pkg *IcopPackage) SetCommandType(CommandType string) {
	pkg.CommandType = CommandType
}

func (pkg *IcopPackage) GetCommandCode() string {
	return pkg.CommandCode
}

func (pkg *IcopPackage) SetCommandCode(CommandCode string) {
	pkg.CommandCode = CommandCode
}

func (pkg *IcopPackage) GetCmdSubCode() string {
	return pkg.CmdSubCode
}

func (pkg *IcopPackage) SetCmdSubCode(CmdSubCode string) {
	pkg.CmdSubCode = CmdSubCode
}

func (pkg *IcopPackage) GetIsUpSend() bool {
	return pkg.IsUpSend
}

func (pkg *IcopPackage) SetIsUpSend(IsUpSend bool) {
	pkg.IsUpSend = IsUpSend
}

func (pkg *IcopPackage) GetData() string {
	return pkg.Data
}

func (pkg *IcopPackage) SetData(Data string) {
	pkg.Data = Data
}
