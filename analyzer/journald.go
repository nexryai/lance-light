package analyzer

import "lance-light/core"

// getLogsFromJournald 特定のタイムスタンプ以降のログを取得する。空の場合は全部取得する。
func getLogsFromJournald(since string) {

	var args []string

	if since == "" {
		args = append(args, "-o", "json", "--no-pager")
	} else {
		args = append(args, "--since="+since, "-o", "json", "--no-pager")
	}

	result := core.ExecCommandGetResult("journalctl", args)
	core.MsgDebug(result[0])
}

// GetJournaldLog 最後に取得した時点から現在までのログを取得する
func GetJournaldLog() {
	since := ""
	getLogsFromJournald(since)
}
