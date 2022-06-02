package {{.pkgName}}

import (
	{{.imports}}
)

var ErrRequestFail = errors.New("http request fail")

func (cc *ClientContext) HandleResponse(resp *http.Response, val interface{}) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		if cc.errType != nil {
			v := reflect.New(cc.errType)
			err = json.Unmarshal(body, v.Elem().Addr().Interface())
			if err != nil {
				return err
			}
			return v.Elem().Interface().(error)
		}
		return ErrRequestFail
	}

	if err := json.Unmarshal(body, val); err != nil {
		return err
	}
	return nil
}
