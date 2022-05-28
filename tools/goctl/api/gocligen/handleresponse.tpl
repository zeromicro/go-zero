package {{.pkgName}}

import (
	{{.imports}}
)

var ErrRequestFail = errors.New("http request fail")

func HandleResponse(resp *http.Response, val interface{}) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		// todo: add your error logic here and delete this line

		return ErrRequestFail
	}

	if err := json.Unmarshal(body, val); err != nil {
		// todo: add your error logic here and delete this line

		return err
	}
	return nil
}
