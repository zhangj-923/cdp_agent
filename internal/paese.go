package internal

import (
	"cdp_agent/common"
)

func parseHostInfo(hostInfo map[string]interface{}, alertMetricDict map[string]string) []*common.Metric {
	if len(hostInfo) == 0 {
		return nil
	}
	metrics := make([]*common.Metric, 0)
	for _, hosts := range hostInfo["server_data"].(map[string]interface{})["Hosts"].([]interface{}) {
		host := hosts.(map[string]interface{})

		cdpHostTags := map[string]interface{}{
			"clientIp": host["IP"].(string),
			"agentId":  host["AgentId"].(string),
		}
		cdpHostFileds := map[string]interface{}{
			"agentVersion": host["AgentVersion"].(string),
			"hostName":     host["Hostname"].(string),
			"serverIp":     host["server_ip"].(string),
		}

		cdpHostMetrics, err := common.Format(cdpHostFileds, cdpHostTags, "cdp.hostinfo", alertMetricDict)
		if err != nil {
			common.Error.Println(err)
			return make([]*common.Metric, 0)
		}
		metrics = append(metrics, cdpHostMetrics...)
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

			cdpDiskTags := map[string]interface{}{
				"clientIp": diskInfo["IP"].(string),
				"agentId":  diskInfo["AgentId"].(string),
				"cdpId":    disk["CdpId"].(string),
			}
			cdpDiskFileds := map[string]interface{}{
				"state":        disk["State"].(int64),
				"mode":         disk["Mode"].(int64),
				"encryptMode":  disk["EncryptMode"].(int64),
				"backupRate":   disk["BackupRate"].(float64),
				"finishedRate": disk["FinishedRate"].(string),
				"beginTime":    disk["BeginTime"].(string),
				"endTime":      disk["EndTime"].(string),
				"capacity":     disk["Capacity"].(float64),
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
