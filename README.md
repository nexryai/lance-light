## LanceLight Firewall

### これは何
nftablesのラッパーとして機能する、軽量でシンプルかつセキュアでエッセンシャルなファイアウォール。  
開発中です。


### 設計原則

#### シンプルで使いやすい
設定は全てyamlで管理します。もう`ufw status`を連打する必要はありません！  
さらにシングルバイナリで動作するためインストールも簡単です。

#### セキュア
不適切な設定には警告を表示します。またメモリセーフなGoで書かれています。

#### 数々の便利機能
Cloudflareからのアクセスのみを許可したり、特定のインターフェースでのみアクセスを許可する、TorやAbuseIPDBでマークされてるIPからのアクセスを禁止するなど、よりサーバーをセキュアにするのに役立つ機能が搭載されています。  
ルーターを自作したいですか？ LanceLightには簡易的なホームルーターの構築を支援する機能が搭載されています！  


### 使い方
注意: 使う前に既存のファイアウォールを無効化してください

```
wget https://raw.githubusercontent.com/nexryai/lance-light/main/install.sh
sudo bash install.sh

# ルールを編集
sudo nano /etc/lance.yml

# 適用
sudo llfctl enable

# 起動時に適用されるようにする
sudo systemctl enable lance
```
#### ソースからビルド
```
# ビルド
make

# インストール
sudo make install
```

### 既知の問題
#### IPv6のIPを指定しても許可されない / IPv6アドレスを入れるとエラーになることがある / AllowIP: "cloudflarev6"が効かない
IPv4を使ってください  
詳細は[issue](https://github.com/nexryai/lance-light/issues/14)参照

#### RHEL系に入れたらsystemd経由で扱おうとするとPermissionで文句言われる
SELinuxのせいです。  
`sudo restorecon -R /usr/local/llfctl`で解決します。

#### Dockerと競合する
`allowAllFwd`をtrueにすれば回避できます。

### Special Thanks
 - ChatGPT
