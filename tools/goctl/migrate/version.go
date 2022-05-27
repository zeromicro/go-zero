package migrate

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/tools/goctl/util/console"
)

var client = http.Client{
	Timeout: 5 * time.Second,
}

func getLatest(repo string, verbose bool) ([]string, error) {
	proxies := goProxy()
	for _, proxy := range proxies {
		if verbose {
			console.Info("use go proxy %q", proxy)
		}
		log := func(err error) {
			console.Warning("get latest versions failed from proxy %q, error: %+v", proxy, err)
		}
		resp, err := client.Get(fmt.Sprintf("%s/%s/@v/list", proxy, repo))
		if err != nil {
			log(err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			log(fmt.Errorf("%s", resp.Status))
			continue
		}
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log(err)
			continue
		}
		versionStr := string(data)
		versions := strings.Fields(versionStr)
		return versions, nil
	}
	return []string{}, nil
}
