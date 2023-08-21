package report

import (
	"encoding/json"
	"fmt"
	"golang.org/x/exp/slices"
	"lance-light/core"
	"lance-light/ip"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

func ReportAbuseIPs(config *core.Config, reportToAbuseIPDB bool) {
	// journalctlコマンドを実行して出力を取得
	journalctlArgs := []string{"-xe", "-o", "json", "--grep", "[LanceLight]", "--no-pager", "_TRANSPORT=kernel", "--since", fmt.Sprintf("%d minute ago", config.Report.ReportInterval)}

	// SRCとDPTを格納するマップ
	srcAndDpt := make(map[string][]int)

	// JSONをパースしてSRCとDPTを抽出し、マップに格納
	lines := core.ExecCommandGetResult("journalctl", journalctlArgs)
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		var data map[string]interface{}
		err := json.Unmarshal([]byte(line), &data)
		if err != nil {
			continue
		}

		if message, ok := data["MESSAGE"].(string); ok {
			matchSrc := regexp.MustCompile(`SRC=([0-9.]+)`).FindStringSubmatch(message)
			matchDpt := regexp.MustCompile(`DPT=(\d+)`).FindStringSubmatch(message)

			if len(matchSrc) > 1 && len(matchDpt) > 1 {
				srcIP := matchSrc[1]
				dptPort, _ := strconv.Atoi(matchDpt[1])

				if _, exists := srcAndDpt[srcIP]; exists {
					srcAndDpt[srcIP] = append(srcAndDpt[srcIP], dptPort)
				} else {
					srcAndDpt[srcIP] = []int{dptPort}
				}
			}
		}
	}

	// 結果を出力
	for srcIp, dportList := range srcAndDpt {

		// 信頼されたグローバルIPか通報できないIPならパス
		if slices.Contains(config.Report.TrustedIPs, srcIp) || !ip.IsReportableAddress(srcIp) {
			continue
		}

		comment := fmt.Sprintf("Blocked by LanceLight (%s -> :%v)", srcIp, dportList)
		core.MsgInfo(comment)

		if reportToAbuseIPDB {
			fmt.Println("reporting！ ;)")
			apiKey := config.Report.AbuseIpDbAPIKey

			if apiKey == "" {
				core.ExitOnError(fmt.Errorf("invalid config file"), "Invalid API Key")
			}

			url := "https://api.abuseipdb.com/api/v2/report"
			categories := "14"

			// HTTPリクエストを作成
			req, err := http.NewRequest("POST", url, nil)
			core.ExitOnError(err, "http.NewRequest error")

			// クエリパラメータを追加
			q := req.URL.Query()
			q.Add("ip", srcIp)
			q.Add("categories", categories)
			q.Add("comment", comment)
			req.URL.RawQuery = q.Encode()

			// HTTPヘッダーを設定
			req.Header.Set("Key", apiKey)
			req.Header.Set("Accept", "application/json")

			// HTTPリクエストを送信
			client := &http.Client{}
			resp, err := client.Do(req)
			core.ExitOnError(err, "http request error")

			if resp.StatusCode == http.StatusTooManyRequests {
				// 同じIPを何回も通報すると429になることがあるけど無視する
				core.MsgDebug("skipped")
				time.Sleep(2 * time.Second)
			} else if resp.StatusCode != http.StatusOK {
				core.MsgErr("Failed to report. resp: " + resp.Status)
			} else {
				time.Sleep(2 * time.Second)
			}

			resp.Body.Close()
		}
	}
}
