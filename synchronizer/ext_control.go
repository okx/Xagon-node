package synchronizer

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/l1_parallel_sync"
)

const (
	externalControlFilename = "/tmp/synchronizer_in"
	externalOutputFilename  = "/tmp/synchronizer_out"
	filePermissions         = 0644
	sleepTimeToReadFile     = 500 * time.Millisecond
)

// This is a local end-point in filesystem to send commands to a running synchronizer
// this is used for debugging purposes, to provide a way to reproduce some situations that are difficult
// to reproduce in a real test.
// It accept next commands:
// l1_producer_stop: stop producer
// l1_orchestrator_reset: reset orchestrator to a given block number
//
// example of usage (first you need to run the service):
// echo "l1_producer_stop" >> /tmp/synchronizer_in
// echo "l1_orchestrator_reset|8577060" >> /tmp/synchronizer_in

// ExtCmdArgs is the type of the arguments of the command
type ExtCmdArgs []string

// ExtControlCmd is the interface of the external  command
type ExtControlCmd interface {
	// FunctionName returns the name of the function to be called example: "l1_producer_stop"
	FunctionName() string
	// ValidateArguments validates the arguments of the command, returns nil if ok, error if not
	ValidateArguments(ExtCmdArgs) error
	// Process the command
	// args: the arguments of the command
	// return: string with the output and an error
	Process(ExtCmdArgs) (string, error)
	// Help returns the help of the command
	Help() string
}

type externalCmdControl struct {
	//producer     *l1_parallel_sync.L1RollupInfoProducer
	//orquestrator *l1_parallel_sync.L1SyncOrchestration
	RegisteredCmds map[string]ExtControlCmd
}

func newExternalCmdControl(producer *l1_parallel_sync.L1RollupInfoProducer, orquestrator *l1_parallel_sync.L1SyncOrchestration) *externalCmdControl {
	res := &externalCmdControl{
		RegisteredCmds: make(map[string]ExtControlCmd),
	}
	res.RegisterCmd(&helpCmd{externalControl: res})
	res.RegisterCmd(&l1OrchestratorResetCmd{orquestrator: orquestrator})
	res.RegisterCmd(&l1ProducerStopCmd{producer: producer})
	return res
}

// RegisterCmd registers a command
func (e *externalCmdControl) RegisterCmd(cmd ExtControlCmd) {
	if e.RegisteredCmds == nil {
		e.RegisteredCmds = make(map[string]ExtControlCmd)
	}
	e.RegisteredCmds[cmd.FunctionName()] = cmd
}

// GetCmd returns a command by its name
func (e *externalCmdControl) GetCmd(functionName string) (ExtControlCmd, error) {
	cmd, ok := e.RegisteredCmds[functionName]
	if !ok {
		return nil, errors.New("command not found")
	}
	return cmd, nil
}

func (e *externalCmdControl) start() {
	log.Infof("EXT:start: starting external control opening %s", externalControlFilename)
	file, err := os.OpenFile(externalControlFilename, os.O_APPEND|os.O_CREATE|os.O_RDONLY, filePermissions)
	if err != nil {
		log.Warnf("EXT:start:error opening file %s: %v", externalControlFilename, err)
		return
	}
	_, err = file.Seek(0, io.SeekEnd)
	if err != nil {
		log.Warnf("EXT:start:error seeking file %s: %v", externalControlFilename, err)
	}
	go e.readFile(file)
}

// https://medium.com/@arunprabhu.1/tailing-a-file-in-golang-72944204f22b
func (e *externalCmdControl) readFile(file *os.File) {
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		for {
			line, err := reader.ReadString('\n')

			if err != nil {
				if err == io.EOF {
					// without this sleep you would hogg the CPU
					time.Sleep(sleepTimeToReadFile)
					continue
				}
				break
			}
			log.Infof("EXT:readFile: new command: %s", line)
			cmd, cmdArgs, err := e.parse(line)
			if err != nil {
				log.Warnf("EXT:readFile: error parsing command %s:err %s", line, err)
				continue
			}
			e.process(cmd, cmdArgs)
		}
	}
}

func (e *externalCmdControl) parse(line string) (ExtControlCmd, ExtCmdArgs, error) {
	cmd := strings.Split(line, "|")
	if len(cmd) < 1 {
		return nil, nil, errors.New("invalid command")
	}
	functionName := strings.TrimSpace(cmd[0])
	args := cmd[1:]
	cmdObj, err := e.GetCmd(functionName)
	if err != nil {
		return nil, nil, err
	}
	err = cmdObj.ValidateArguments(args)
	if err != nil {
		return nil, nil, err
	}
	return cmdObj, args, nil
}

func (e *externalCmdControl) process(cmd ExtControlCmd, args ExtCmdArgs) {
	fullFunc, err := fmt.Printf("%s(%s)", cmd.FunctionName(), strings.Join(args, ","))
	if err != nil {
		log.Warnf("EXT:readFile: error composing cmd %s:err %s", cmd.FunctionName(), err)
		return
	}
	output, err := cmd.Process(args)
	if err != nil {
		log.Warnf("EXT:readFile: error processing command %s:err %s", fullFunc, err)
		return
	}
	log.Warnf("EXT:readFile: command %s processed with output: %s", fullFunc, output)
}

// COMMANDS IMPLEMENTATION
// HELP
type helpCmd struct {
	externalControl *externalCmdControl
}

func (h *helpCmd) FunctionName() string {
	return "help"
}
func (h *helpCmd) ValidateArguments(args ExtCmdArgs) error {
	if len(args) > 0 {
		return errors.New(h.FunctionName() + " command does not accept arguments")
	}
	return nil
}

func (h *helpCmd) Process(args ExtCmdArgs) (string, error) {
	var help string
	for _, cmd := range h.externalControl.RegisteredCmds {
		help += cmd.Help() + "\n"
	}
	return help, nil
}
func (h *helpCmd) Help() string {
	return h.FunctionName() + ": show the help of the commands"
}

// COMMANDS "l1_orchestrator_reset"
type l1OrchestratorResetCmd struct {
	orquestrator *l1_parallel_sync.L1SyncOrchestration
}

func (h *l1OrchestratorResetCmd) FunctionName() string {
	return "l1_orchestrator_reset"
}

func (h *l1OrchestratorResetCmd) ValidateArguments(args ExtCmdArgs) error {
	if len(args) != 1 {
		return errors.New(h.FunctionName() + " needs 1 argument")
	}
	_, err := strconv.ParseUint(strings.TrimSpace(args[0]), 10, 64)
	if err != nil {
		return fmt.Errorf("error parsing block number: %s err:%w", args[0], err)
	}
	return nil
}
func (h *l1OrchestratorResetCmd) Process(args ExtCmdArgs) (string, error) {
	blockNumber, err := strconv.ParseUint(strings.TrimSpace(args[0]), 10, 64)
	if err != nil {
		return "error param", err
	}
	log.Warnf("EXT:"+h.FunctionName()+": calling orchestrator reset(%d)", blockNumber)
	h.orquestrator.Reset(blockNumber)
	res := fmt.Sprintf("EXT: "+h.FunctionName()+": reset to block %d", blockNumber)
	return res, nil
}

func (h *l1OrchestratorResetCmd) Help() string {
	return h.FunctionName() + ": reset L1 parallel sync orchestrator to a given block number"
}

// COMMANDS l1_producer_stop
type l1ProducerStopCmd struct {
	producer *l1_parallel_sync.L1RollupInfoProducer
}

func (h *l1ProducerStopCmd) FunctionName() string {
	return "l1_producer_stop"
}

func (h *l1ProducerStopCmd) ValidateArguments(args ExtCmdArgs) error {
	if len(args) > 0 {
		return errors.New(h.FunctionName() + " command does not accept arguments")
	}
	return nil
}
func (h *l1ProducerStopCmd) Process(args ExtCmdArgs) (string, error) {
	log.Warnf("EXT:" + h.FunctionName() + ": calling producer stop")
	h.producer.Stop()
	res := "EXT: " + h.FunctionName() + ": producer stopped"
	return res, nil
}

func (h *l1ProducerStopCmd) Help() string {
	return h.FunctionName() + ": stop L1 rollup info producer"
}
