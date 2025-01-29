package migrate

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/tools/goctl/util/console"
)

var client = http.Client{
	Timeout: 5 * time.Second,
}

func getLatest(repo string, verbose bool) ([]string, error) {
	log := func(err error) {
		console.Warning("get latest versions failed, error: %+v", err)
	}
	resp, err := client.Get(fmt.Sprintf("%s/@v/list", repo))
	if err != nil {
		log(err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("%s", resp.Status)
		log(err)
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log(err)
		return nil, err
	}

	versionStr := string(data)
	versions := strings.Fields(versionStr)
	return versions, nil
}
