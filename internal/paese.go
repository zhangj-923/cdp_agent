package internal

import (
	"cdp_agent/common"
	"strconv"
)

func parseHostInfo(hostInfo map[string]interface{}, alertMetricDict map[string]string) []*common.Metric {
	if len(hostInfo) == 0 {
		return nil
	}
	metrics := make([]*common.Metric, 0)
	for _, serverDatas := range hostInfo["server_data"].([]interface{}) {
		server := serverDatas.(map[string]interface{})
		for _, hosts := range server["Hosts"].([]interface{}) {
			host := hosts.(map[string]interface{})
			cdpHostTags := map[string]interface{}{
				"clientIp": host["IP"].(string),
				"agentId":  int(host["AgentId"].(float64)),
				"sub_key":  host["Hostname"].(string),
			}
			cdpHostFileds := map[string]interface{}{
				"agentVersion": host["AgentVersion"].(string),
				"hostname":     host["Hostname"].(string),
				"serverIp":     host["server_ip"].(string),
				"Delaytime":    int(host["Offline"].(float64)),
			}
			if host["Offline"].(float64) > 0 {
				cdpHostFileds["offline"] = 1
			} else {
				cdpHostFileds["offline"] = 0
			}

			cdpHostMetrics, err := common.Format(cdpHostFileds, cdpHostTags, "cdp.hostinfo", alertMetricDict)
			if err != nil {
				common.Error.Println(err)
				return make([]*common.Metric, 0)
			}
			metrics = append(metrics, cdpHostMetrics...)
		}
	}
	return metrics
}

func parseDiskInfo(infos []map[string]interface{}, alertMetricDict map[string]string) []*common.Metric {
	if len(infos) == 0 {
		return nil
	}
	metrics := make([]*common.Metric, 0)
	for _, diskInfo := range infos {
		for _, disks := range diskInfo["Disks"].([]interface{}) {
			disk := disks.(map[string]interface{})
			flag := false
			for _, partitions := range disk["Partitions"].([]interface{}) {
				partition := partitions.(map[string]interface{})
				if int(partition["Backup"].(float64)) == 1 {
					flag = true
					break
				}
			}
			if !flag && disk["EndTime"] != nil {
				flag = true
			}
			if !flag {
				continue
			}
			cdpDiskTags := map[string]interface{}{
				"clientIp": diskInfo["IP"].(string),
				"agentId":  int(diskInfo["AgentId"].(float64)),
				"diskId":   disk["DiskGuid"].(string),
				"sub_key":  diskInfo["IP"].(string) + ":" + strconv.FormatFloat(disk["CdpId"].(float64), 'E', -1, 64),
			}
			cdpDiskFileds := map[string]interface{}{
				"cdpId":        int(disk["CdpId"].(float64)),
				"state":        int(disk["State"].(float64)),
				"mode":         int(disk["Mode"].(float64)),
				"encryptMode":  int(disk["EncryptMode"].(float64)),
				"backupRate":   disk["BackupRate"].(string),
				"finishedRate": disk["FinishedRate"].(string),
				"beginTime":    disk["BeginTime"].(string),
				"capacity":     disk["Capacity"].(float64),
			}
			if disk["EndTime"] != nil {
				cdpDiskFileds["endTime"] = disk["EndTime"].(string)
			} else {
				cdpDiskFileds["endTime"] = ""
			}
			cdpDiskMetrics, err := common.Format(cdpDiskFileds, cdpDiskTags, "cdp.diskinfo", alertMetricDict)
			if err != nil {
				common.Error.Println(err)
				return make([]*common.Metric, 0)
			}
			metrics = append(metrics, cdpDiskMetrics...)
		}
	}
	return metrics
}
