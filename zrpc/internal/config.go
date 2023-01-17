package internal

type (
	// ClientMiddlewaresConf defines whether to use client middlewares.
	ClientMiddlewaresConf struct {
		Trace      bool `json:",default=true"`
		Duration   bool `json:",default=true"`
		Prometheus bool `json:",default=true"`
		Breaker    bool `json:",default=true"`
		Timeout    bool `json:",default=true"`
	}

	// ServerMiddlewaresConf defines whether to use server middlewares.
	ServerMiddlewaresConf struct {
		Trace      bool `json:",default=true"`
		Recover    bool `json:",default=true"`
		Stat       bool `json:",default=true"`
		Prometheus bool `json:",default=true"`
		Breaker    bool `json:",default=true"`
	}
)
