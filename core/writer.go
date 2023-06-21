package core

import (
	"io/ioutil"
	"os"
	"strings"
)

func WriteToFile(content []string, path string) {

	// 配列の中身を1行ずつout.txtに上書きする
	contentString := strings.Join(content, "\n") // 改行で要素を結合

	if _, err := os.Stat(path); os.IsNotExist(err) {
		// ファイルが存在しない場合は作成する
		file, err := os.Create(path)
		ExitOnError(err, "Failed to create file.")
		defer file.Close()
	}

	err := ioutil.WriteFile(path, []byte(contentString), os.ModePerm)
	ExitOnError(err, "Failed to write rules.")

}
