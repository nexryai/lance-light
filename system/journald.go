package system

import "lance-light/core"

func getLogsFromJournald(sinceDate int64) {
	result := core.ExecCommandGetResult("journalctl", []string{"-xe", "-o", "json", "--no-pager"})
	core.MsgDebug(result[0])
}
