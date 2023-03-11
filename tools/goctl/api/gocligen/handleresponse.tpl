package {{.pkgName}}

import (
	{{.imports}}
)

type CodeError struct {
	Status int
	Desc   string
}

func (e *CodeError) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("status: %d, desc: %s", e.Status, e.Desc)
}

func (cc *ClientContext) HandleResponse(resp *http.Response, val interface{}) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

    if resp.StatusCode != http.StatusOK {
		if cc.errType != nil {
			v := reflect.New(cc.errType.Elem())
			err = json.Unmarshal(body, v.Elem().Addr().Interface())
			if err != nil {
				return err
			}
			e, ok := v.Elem().Addr().Interface().(error)
			if !ok {
				return &CodeError{
					Status: resp.StatusCode,
					Desc:   string(body),
				}
			}
			return e
		}
        return &CodeError{
			Status: resp.StatusCode,
			Desc:   string(body),
		}
	}

	if err := json.Unmarshal(body, val); err != nil {
		return err
	}
	return nil
}
