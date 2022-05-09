package internal

import (
	"cdp_agent/common"
)

func MockData(alertMetricDict map[string]string) []*common.Metric {
	metrics := make([]*common.Metric, 0)

	cdpHostTags11 := map[string]interface{}{
		"clientIp": "192.168.1.75",
		"agentId":  1,
	}
	cdpHostFileds11 := map[string]interface{}{
		"hostname":     "whdata",
		"agentVersion": "5.1.0.1",
		"serverIp":     "192.168.1.83",
	}
	cdpHostMetrics11, err := common.Format(cdpHostFileds11, cdpHostTags11, "cdp.hostinfo", alertMetricDict)
	if err != nil {
		common.Error.Println(err)
		return make([]*common.Metric, 0)
	}
	metrics = append(metrics, cdpHostMetrics11...)

	cdpDiskTags11 := map[string]interface{}{
		"clientIp": "192.168.1.75",
		"agentId":  1,
		"cdpId":    1,
	}
	cdpDiskFileds11 := map[string]interface{}{
		"state":        3,
		"mode":         2,
		"encryptMode":  0,
		"backupRate":   "0.00",
		"finishedRate": "100.00%",
		"beginTime":    "2022-03-11 15:54:08",
		"endTime":      "2022-03-11 15:54:16",
		"capacity":     104853504.00,
	}
	cdpDiskMetrics11, err := common.Format(cdpDiskFileds11, cdpDiskTags11, "cdp.diskinfo", alertMetricDict)
	if err != nil {
		common.Error.Println(err)
		return make([]*common.Metric, 0)
	}
	metrics = append(metrics, cdpDiskMetrics11...)

	cdpdisktags12 := map[string]interface{}{
		"clientIp": "192.168.1.75",
		"agentId":  1,
		"cdpId":    2,
	}
	cdpdiskfileds12 := map[string]interface{}{
		"state":        2,
		"mode":         1,
		"encryptMode":  1,
		"backupRate":   "65.00",
		"finishedRate": "85.00%",
		"beginTime":    "2022-05-09 15:54:08",
		"endTime":      "2022-05-11 15:54:16",
		"capacity":     92274688.00,
	}
	cdpdiskmetrics12, err := common.Format(cdpdiskfileds12, cdpdisktags12, "cdp.diskinfo", alertMetricDict)
	if err != nil {
		common.Error.Println(err)
		return make([]*common.Metric, 0)
	}
	metrics = append(metrics, cdpdiskmetrics12...)

	cdpHostTags21 := map[string]interface{}{
		"clientIp": "192.168.1.53",
		"agentId":  1,
	}
	cdpHostFileds21 := map[string]interface{}{
		"hostname":     "USER-M650SN7FHL",
		"agentVersion": "5.2.0.0",
		"serverIp":     "192.168.1.78",
	}
	cdpHostMetrics21, err := common.Format(cdpHostFileds21, cdpHostTags21, "cdp.hostinfo", alertMetricDict)
	if err != nil {
		common.Error.Println(err)
		return make([]*common.Metric, 0)
	}
	metrics = append(metrics, cdpHostMetrics21...)

	cdpDiskTags21 := map[string]interface{}{
		"clientIp": "192.168.1.53",
		"agentId":  1,
		"cdpId":    4,
	}
	cdpDiskFileds21 := map[string]interface{}{
		"state":        3,
		"mode":         2,
		"encryptMode":  0,
		"backupRate":   "0.00",
		"finishedRate": "100.00%",
		"beginTime":    "2022-03-11 15:54:11",
		"endTime":      "2022-03-11 15:54:16",
		"capacity":     209709056.00,
	}
	cdpDiskMetrics21, err := common.Format(cdpDiskFileds21, cdpDiskTags21, "cdp.diskinfo", alertMetricDict)
	if err != nil {
		common.Error.Println(err)
		return make([]*common.Metric, 0)
	}
	metrics = append(metrics, cdpDiskMetrics21...)

	cdpHostTags22 := map[string]interface{}{
		"clientIp": "192.168.1.54",
		"agentId":  2,
	}
	cdpHostFileds22 := map[string]interface{}{
		"hostname":     "backup",
		"agentVersion": "5.2.0.0",
		"serverIp":     "192.168.1.78",
	}
	cdpHostMetrics22, err := common.Format(cdpHostFileds22, cdpHostTags22, "cdp.hostinfo", alertMetricDict)
	if err != nil {
		common.Error.Println(err)
		return make([]*common.Metric, 0)
	}
	metrics = append(metrics, cdpHostMetrics22...)
	return metrics
}
