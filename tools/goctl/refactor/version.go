package refactor

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var client = http.Client{
	Timeout: 5 * time.Second,
}

func getLatest(repo string) ([]string, error) {
	resp, err := client.Get(fmt.Sprintf("https://goproxy.cn/%s/@v/list", repo))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", resp.Status)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	versionStr:=string(data)
	versions := strings.Fields(versionStr)
	return versions, nil
}
