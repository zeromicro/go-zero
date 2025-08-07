package sqlx

import "errors"

var (
	errEmptyDatasource = errors.New("empty datasource")
	errEmptyDriverName = errors.New("empty driver name")
)

// SqlConf defines the configuration for sqlx.
type SqlConf struct {
	DataSource string
	DriverName string   `json:",default=mysql"`
	Replicas   []string `json:",optional"`
	Policy     string   `json:",default=round-robin,options=round-robin|random"`
}

// Validate validates the SqlxConf.
func (sc SqlConf) Validate() error {
	if len(sc.DataSource) == 0 {
		return errEmptyDatasource
	}

	if len(sc.DriverName) == 0 {
		return errEmptyDriverName
	}

	return nil
}
