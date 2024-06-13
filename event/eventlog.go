package event

import (
	"context"
	"encoding/json"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
)

// EventLog is the main struct for the event log
type EventLog struct {
	cfg     Config
	storage Storage
}

// NewEventLog creates and initializes an instance of EventLog
func NewEventLog(cfg Config, storage Storage) *EventLog {
	return &EventLog{
		cfg:     cfg,
		storage: storage,
	}
}

// LogEvent is used to store an event for runtime debugging
func (e *EventLog) LogEvent(ctx context.Context, event *Event) error {
	return e.storage.LogEvent(ctx, event)
}

// LogExecutorError is used to store Executor error for runtime debugging
func (e *EventLog) LogExecutorError(ctx context.Context, responseError executor.ExecutorError, processBatchRequest interface{}) {
	timestamp := time.Now()

	// if it's a user related error, ignore it
	if responseError == executor.ExecutorError_EXECUTOR_ERROR_SM_MAIN_COUNTERS_OVERFLOW_STEPS ||
		responseError == executor.ExecutorError_EXECUTOR_ERROR_SM_MAIN_COUNTERS_OVERFLOW_KECCAK ||
		responseError == executor.ExecutorError_EXECUTOR_ERROR_SM_MAIN_COUNTERS_OVERFLOW_BINARY ||
		responseError == executor.ExecutorError_EXECUTOR_ERROR_SM_MAIN_COUNTERS_OVERFLOW_MEM ||
		responseError == executor.ExecutorError_EXECUTOR_ERROR_SM_MAIN_COUNTERS_OVERFLOW_ARITH ||
		responseError == executor.ExecutorError_EXECUTOR_ERROR_SM_MAIN_COUNTERS_OVERFLOW_PADDING ||
		responseError == executor.ExecutorError_EXECUTOR_ERROR_SM_MAIN_COUNTERS_OVERFLOW_POSEIDON ||
		responseError == executor.ExecutorError_EXECUTOR_ERROR_INVALID_BATCH_L2_DATA {
		return
	}

	log.Errorf("error found in the executor: %v at %v", responseError, timestamp)
	payload, err := json.Marshal(processBatchRequest)
	if err != nil {
		log.Errorf("error marshaling payload: %v", err)
	} else {
		event := &Event{
			ReceivedAt:  timestamp,
			Source:      Source_Node,
			Component:   Component_Executor,
			Level:       Level_Error,
			EventID:     EventID_ExecutorError,
			Description: responseError.String(),
			Json:        string(payload),
		}
		err = e.storage.LogEvent(ctx, event)
		if err != nil {
			log.Errorf("error storing event: %v", err)
		}
	}
}
