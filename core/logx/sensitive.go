package logx

// Sensitive is an interface that defines a method for masking sensitive information in logs.
// It is typically implemented by types that contain sensitive data,
// such as passwords or personal information.
// Infov, Errorv, Debugv, and Slowv methods will call this method to mask sensitive data.
// The values in LogField will also be masked if they implement the Sensitive interface.
type Sensitive interface {
	// MaskSensitive masks sensitive information in the log.
	MaskSensitive() any
}

// maskSensitive returns the value returned by MaskSensitive method,
// if the value implements Sensitive interface.
func maskSensitive(v any) any {
	if s, ok := v.(Sensitive); ok {
		return s.MaskSensitive()
	}

	return v
}
