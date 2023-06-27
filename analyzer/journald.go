package analyzer

import (
	"lance-light/core"
	"lance-light/memory"
	"time"
)

// 最後に取得した時点から現在までのログを取得する
func getLogsFromJournald() {

	// レジストリから最後にログを取得したタイムスタンプを取得する
	timestampString := memory.GetRegistryValue("analyzer.journald.lastRecordTimestamp")

	var args []string

	if timestampString == "" {
		// ログを取得するのが初めての場合、analyzer.journald.lastRecordTimestampが空になる。全てのログを取得する。
		args = append(args, "-o", "json", "--no-pager")
	} else {
		// StringからUnixタイムスタンプに変換
		timestamp, err := time.Parse("2006-01-02 15:04:05", timestampString)
		core.ExitOnError(err, core.GenBugCodeMessage("ea14b6c8-aba2-4e2e-8fda-faa448e5771d"))

		// 3秒前の時間を計算
		// 重複排除を後に行うので余裕をもって取得する
		prevTime := timestamp.Add(-3 * time.Second)

		// journaldでも理解できる形式に変換する (unix時刻くらい理解しろや)
		layout := "2006-01-02 15:04:05"
		timestampStringForJournald := prevTime.Format(layout)

		// since以降のログを取得する
		args = append(args, "--since="+timestampStringForJournald, "-o", "json", "--no-pager")
	}

	result := core.ExecCommandGetResult("journalctl", args)
	core.MsgDebug(result[0])

	// lastRecordTimestampを直接書き換えると上位の関数で問題が発生して処理が中断されるなどしたときに不整合になるのでpendingRecordTimestampに記録する。
	// 処理が成功したら呼び出し元でpendingRecordTimestampをlastRecordTimestampに代入する。
	memory.SetRegistryValue("analyzer.journald.pendingRecordTimestamp", core.GetUnixTimestampString())
}

func GetJournaldLog() {
	getLogsFromJournald()
	memory.SetRegistryValue("analyzer.journald.lastRecordTimestamp", memory.GetRegistryValue("analyzer.journald.pendingRecordTimestamp"))
}
