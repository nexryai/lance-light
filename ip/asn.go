package ip

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"lance-light/core"
	"net/http"
)

type PrefixData struct {
	Prefixes []struct {
		Prefix    string `json:"prefix"`
		Timelines []struct {
			StartTime string `json:"starttime"`
			EndTime   string `json:"endtime"`
		} `json:"timelines"`
	} `json:"prefixes"`
}

func GetIpRangeFromASN(asn string) []string {
	url := "https://stat.ripe.net/data/announced-prefixes/data.json?resource=AS" + asn
	ipCidr := []string{}

	core.MsgDebug("send request: " + url)

	response, err := http.Get(url)
	core.ExitOnError(err, "Failed to convert ASN to IP CIDR. Request did not succeed.")
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(response.Body)
		core.ExitOnError(err, "Failed to convert ASN to IP CIDR. The request was successful, but parsing failed.")

		var data struct {
			Data PrefixData `json:"data"`
		}

		err = json.Unmarshal(body, &data)
		core.ExitOnError(err, "Failed to convert ASN to IP CIDR. The request was successful, but parsing failed.")

		for _, prefix := range data.Data.Prefixes {
			ipCidr = append(ipCidr, prefix.Prefix)
		}
	} else {
		core.ExitOnError(errors.New("request failed"), "Failed to convert ASN to IP CIDR. An error code was returned from the server.")
	}

	//念のため確認
	if !CheckIPAddresses(ipCidr) {
		core.ExitOnError(errors.New("invalid IP from api"), "Failed to convert ASN to IP CIDR. The request was successful, but an invalid IP address was detected.")
	}

	return ipCidr
}
