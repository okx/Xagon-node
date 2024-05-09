package actions

import (
	"reflect"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
)

// ProcessorBase is the base struct for all the processors, if reduces the boilerplate
// implementing the Name, SupportedEvents and SupportedForkIds functions
type ProcessorBase[T any] struct {
	supportedEvent   []etherman.EventOrder
	supportedForkIds []ForkIdType
}

// NewProcessorBase creates and initializes internal fields of an new instance of ProcessorBase
func NewProcessorBase[T any](supportedEvent []etherman.EventOrder, supportedForkIds []ForkIdType) *ProcessorBase[T] {
	p := &ProcessorBase[T]{
		supportedEvent:   supportedEvent,
		supportedForkIds: supportedForkIds,
	}

	return p
}

// Name returns the name of the struct T
func (g *ProcessorBase[T]) Name() string {
	var value T
	a := reflect.TypeOf(value)
	b := a.Name()
	return b
}

// SupportedEvents returns the supported events in the struct
func (p *ProcessorBase[T]) SupportedEvents() []etherman.EventOrder {
	return p.supportedEvent
}

// SupportedForkIds returns the supported forkIds in the struct or the default till incaberry forkId
func (p *ProcessorBase[T]) SupportedForkIds() []ForkIdType {
	if len(p.supportedForkIds) != 0 {
		return p.supportedForkIds
	}
	// returns none
	return []ForkIdType{}
}
