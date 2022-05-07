package internal

import (
	"cdp_agent/common"
)

func MockData(alertMetricDict map[string]string) []*common.Metric {
	metrics := make([]*common.Metric, 0)

	cdpHostTags := map[string]interface{}{
		"clientIp": "192.168.1.75",
		"agentId":  1,
	}
	cdpHostFileds := map[string]interface{}{
		"hostname":     "whdata",
		"agentVersion": "5.1.0.1",
		"serverIp":     "192.168.1.83",
	}
	cdpHostMetrics, err := common.Format(cdpHostFileds, cdpHostTags, "cdp.hostinfo", alertMetricDict)
	if err != nil {
		common.Error.Println(err)
		return make([]*common.Metric, 0)
	}
	metrics = append(metrics, cdpHostMetrics...)

	cdpDiskTags := map[string]interface{}{
		"clientIp": "192.168.1.75",
		"agentId":  1,
		"cdpId":    1,
	}
	cdpDiskFileds := map[string]interface{}{
		"state":        3,
		"mode":         2,
		"encryptMode":  0,
		"backupRate":   "0.00",
		"finishedRate": "100.00%",
		"beginTime":    "2022-03-11 15:54:08",
		"endTime":      "2022-03-11 15:54:16",
		"capacity":     104853504.00,
	}
	cdpDiskMetrics, err := common.Format(cdpDiskFileds, cdpDiskTags, "cdp.diskinfo", alertMetricDict)
	if err != nil {
		common.Error.Println(err)
		return make([]*common.Metric, 0)
	}
	metrics = append(metrics, cdpDiskMetrics...)
	return metrics
}
