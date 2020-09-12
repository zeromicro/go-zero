package stat

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/tal-tech/go-zero/core/logx"
)

const httpTimeout = time.Second * 5

var ErrWriteFailed = errors.New("submit failed")

type RemoteWriter struct {
	endpoint string
}

func NewRemoteWriter(endpoint string) Writer {
	return &RemoteWriter{
		endpoint: endpoint,
	}
}

func (rw *RemoteWriter) Write(report *StatReport) error {
	bs, err := json.Marshal(report)
	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout: httpTimeout,
	}
	resp, err := client.Post(rw.endpoint, "application/json", bytes.NewBuffer(bs))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logx.Errorf("write report failed, code: %d, reason: %s", resp.StatusCode, resp.Status)
		return ErrWriteFailed
	}

	return nil
}
