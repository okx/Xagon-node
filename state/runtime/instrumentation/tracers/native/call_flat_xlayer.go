package native

import (
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/instrumentation/tracers"
)

// SetFlatCallTracerLimit set the limit for flatCallFrame.
func SetFlatCallTracerLimit(t tracers.Tracer, l int) tracers.Tracer {
	if flatTracer, ok := t.(*flatCallTracer); ok {
		flatTracer.limit = l
		return flatTracer
	}
	return t
}
