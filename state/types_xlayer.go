package state

// IsFlatCallTracer returns true when should use flatCallTracer
func (t *TraceConfig) IsFlatCallTracer() bool {
	return t.Tracer != nil && *t.Tracer == "flatCallTracer"
}
