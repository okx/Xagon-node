DOCKERCOMPOSE := docker compose -f docker-compose.yml 
DOCKERCOMPOSEAPPSEQ := zkevm-sequencer
DOCKERCOMPOSEAPPSEQV1TOV2 := zkevm-sequencer-v1tov2
DOCKERCOMPOSEAPPSEQSENDER := zkevm-sequence-sender
DOCKERCOMPOSEAPPSEQSENDERV1TOV2 := zkevm-sequence-sender-v1tov2
DOCKERCOMPOSEAPPL2GASP := zkevm-l2gaspricer
DOCKERCOMPOSEAPPL2GASPV1TOV2 := zkevm-l2gaspricer-v1tov2
DOCKERCOMPOSEAPPAGG := zkevm-aggregator
DOCKERCOMPOSEAPPAGGV1TOV2 := zkevm-aggregator-v1tov2
DOCKERCOMPOSEAPPRPC := zkevm-json-rpc
DOCKERCOMPOSEAPPRPCV1TOV2 := zkevm-json-rpc-v1tov2
DOCKERCOMPOSEAPPSYNC := zkevm-sync
DOCKERCOMPOSEAPPSYNCV1TOV2 := zkevm-sync-v1tov2
DOCKERCOMPOSEAPPETHTXMANAGER := zkevm-eth-tx-manager
DOCKERCOMPOSEAPPETHTXMANAGERV1TOV2 := zkevm-eth-tx-manager-v1tov2
DOCKERCOMPOSESTATEDB := zkevm-state-db
DOCKERCOMPOSEPOOLDB := zkevm-pool-db
DOCKERCOMPOSEEVENTDB := zkevm-event-db
DOCKERCOMPOSENETWORK := zkevm-mock-l1-network
DOCKERCOMPOSEV1TOV2NETWORK := zkevm-v1tov2-l1-network
DOCKERCOMPOSEEXPLORERL1 := zkevm-explorer-l1
DOCKERCOMPOSEEXPLORERL1DB := zkevm-explorer-l1-db
DOCKERCOMPOSEEXPLORERL2 := zkevm-explorer-l2
DOCKERCOMPOSEEXPLORERL2DB := zkevm-explorer-l2-db
DOCKERCOMPOSEEXPLORERRPC := zkevm-explorer-json-rpc
DOCKERCOMPOSEZKPROVER := zkevm-prover
DOCKERCOMPOSEPERMISSIONLESSDB := zkevm-permissionless-db
DOCKERCOMPOSEPERMISSIONLESSNODE := zkevm-permissionless-node
DOCKERCOMPOSEPERMISSIONLESSZKPROVER := zkevm-permissionless-prover
DOCKERCOMPOSENODEAPPROVE := zkevm-approve
DOCKERCOMPOSENODEAPPROVEV1TOV2 := zkevm-approve-v1tov2
DOCKERCOMPOSEMETRICS := zkevm-metrics
DOCKERCOMPOSEGRAFANA := grafana

RUNSTATEDB := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSESTATEDB)
RUNPOOLDB := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEPOOLDB)
RUNEVENTDB := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEEVENTDB)
RUNSEQUENCER := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPPSEQ)
RUNV1TOV2SEQUENCER := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPPSEQV1TOV2)
RUNSEQUENCESENDER := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPPSEQSENDER)
RUNV1TOV2SEQUENCESENDER := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPPSEQSENDERV1TOV2)
RUNL2GASPRICER := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPPL2GASP)
RUNV1TOV2L2GASPRICER := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPPL2GASPV1TOV2)
RUNAGGREGATOR := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPPAGG)
RUNV1TOV2AGGREGATOR := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPPAGGV1TOV2)
RUNJSONRPC := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPPRPC)
RUNV1TOV2JSONRPC := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPPRPCV1TOV2)
RUNSYNC := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPPSYNC)
RUNV1TOV2SYNC := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPPSYNCV1TOV2)
RUNETHTXMANAGER := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPPETHTXMANAGER)
RUNV1TOV2ETHTXMANAGER := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPPETHTXMANAGERV1TOV2)
RUNGRAFANA := DOCKERGID=`stat -c '%g' /var/run/docker.sock` $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEGRAFANA)

RUNL1NETWORK := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSENETWORK)
RUNV1TOV2L1NETWORK := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEV1TOV2NETWORK)
RUNEXPLORERL1 := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEEXPLORERL1)
RUNEXPLORERL1DB := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEEXPLORERL1DB)
RUNEXPLORERL2 := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEEXPLORERL2)
RUNEXPLORERL2DB := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEEXPLORERL2DB)
RUNEXPLORERJSONRPC := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEEXPLORERRPC)
RUNZKPROVER := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEZKPROVER)

RUNPERMISSIONLESSDB := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEPERMISSIONLESSDB)
RUNPERMISSIONLESSNODE := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEPERMISSIONLESSNODE)
RUNPERMISSIONLESSZKPROVER := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEPERMISSIONLESSZKPROVER)

RUNAPPROVE := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSENODEAPPROVE)
RUNV1TOV2APPROVE := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSENODEAPPROVEV1TOV2)

RUNMETRICS := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEMETRICS)

RUN := $(DOCKERCOMPOSE) up -d

STOPSTATEDB := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSESTATEDB) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSESTATEDB)
STOPPOOLDB := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEPOOLDB) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEPOOLDB)
STOPEVENTDB := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEEVENTDB) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEEVENTDB)
STOPSEQUENCER := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPPSEQ) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPPSEQ)
STOPV1TOV2SEQUENCER := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPPSEQV1TOV2) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPPSEQV1TOV2)
STOPSEQUENCESENDER := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPPSEQSENDER) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPPSEQSENDER)
STOPV1TOV2SEQUENCESENDER := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPPSEQSENDERV1TOV2) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPPSEQSENDERV1TOV2)
STOPL2GASPRICER := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPPL2GASP) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPPL2GASP)
STOPV1TOV2L2GASPRICER := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPPL2GASPV1TOV2) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPPL2GASPV1TOV2)
STOPAGGREGATOR := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPPAGG) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPPAGG)
STOPV1TOV2AGGREGATOR := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPPAGGV1TOV2) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPPAGGV1TOV2)
STOPJSONRPC := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPPRPC) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPPRPC)
STOPV1TOV2JSONRPC := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPPRPCV1TOV2) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPPRPCV1TOV2)
STOPSYNC := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPPSYNC) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPPSYNC)
STOPV1TOV2SYNC := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPPSYNCV1TOV2) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPPSYNCV1TOV2)
STOPETHTXMANAGER := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPPETHTXMANAGER) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPPETHTXMANAGER)
STOPV1TOV2ETHTXMANAGER := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPPETHTXMANAGERV1TOV2) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPPETHTXMANAGERV1TOV2)
STOPGRAFANA := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEGRAFANA) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEGRAFANA)

STOPNETWORK := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSENETWORK) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSENETWORK)
STOPV1TOV2NETWORK := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEV1TOV2NETWORK) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEV1TOV2NETWORK)
STOPEXPLORERL1 := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEEXPLORERL1) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEEXPLORERL1)
STOPEXPLORERL1DB := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEEXPLORERL1DB) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEEXPLORERL1DB)
STOPEXPLORERL2 := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEEXPLORERL2) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEEXPLORERL2)
STOPEXPLORERL2DB := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEEXPLORERL2DB) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEEXPLORERL2DB)
STOPEXPLORERJSONRPC := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEEXPLORERRPC) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEEXPLORERRPC)
STOPZKPROVER := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEZKPROVER) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEZKPROVER)

STOPPERMISSIONLESSDB := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEPERMISSIONLESSDB) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEPERMISSIONLESSDB)
STOPPERMISSIONLESSNODE := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEPERMISSIONLESSNODE) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEPERMISSIONLESSNODE)
STOPPERMISSIONLESSZKPROVER := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEPERMISSIONLESSZKPROVER) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEPERMISSIONLESSZKPROVER)

STOPAPPROVE := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSENODEAPPROVE) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSENODEAPPROVE)
STOPV1TOV2APPROVE := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSENODEAPPROVEV1TOV2) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSENODEAPPROVEV1TOV2)

STOPMETRICS := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEMETRICS) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEMETRICS)

STOP := $(DOCKERCOMPOSE) down --remove-orphans

.PHONY: test-full-non-e2e
test-full-non-e2e: stop ## Runs non-e2e tests checking race conditions
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNEVENTDB)
	sleep 2
	$(RUNZKPROVER)
	sleep 7
	$(RUNL1NETWORK)
	sleep 15
	docker ps -a
	docker logs $(DOCKERCOMPOSEZKPROVER)
	trap '$(STOP)' EXIT; MallocNanoZone=0 go test -count=1 -short -race -p 1 -covermode=atomic -coverprofile=../coverage.out  -coverpkg ./... -timeout 200s ../...

.PHONY: test-full-non-e2e-sonar
test-full-non-e2e-sonar: stop ## Runs non-e2e tests checking race conditions
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNEVENTDB)
	sleep 2
	$(RUNZKPROVER)
	sleep 7
	$(RUNL1NETWORK)
	sleep 15
	docker ps -a
	docker logs $(DOCKERCOMPOSEZKPROVER)
	trap '$(STOP)' EXIT; MallocNanoZone=0 go test -count=1 -short -race -p 1 -covermode=atomic -coverprofile=../coverage.out  -coverpkg ./... -timeout 200s ../... -json > ../report.json


.PHONY: test-e2e-group-1
test-e2e-group-1: stop ## Runs group 1 e2e tests checking race conditions
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNEVENTDB)
	sleep 5
	$(RUNZKPROVER)
	docker ps -a
	docker logs $(DOCKERCOMPOSEZKPROVER)
	trap '$(STOP)' EXIT; MallocNanoZone=0 go test -count=1 -failfast -race -v -p 1 -timeout 600s ../ci/e2e-group1/...

.PHONY: test-e2e-group-2
test-e2e-group-2: stop ## Runs group 2 e2e tests checking race conditions
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNEVENTDB)
	sleep 5
	$(RUNZKPROVER)
	docker ps -a
	docker logs $(DOCKERCOMPOSEZKPROVER)
	trap '$(STOP)' EXIT; MallocNanoZone=0 go test -count=1 -failfast -race -v -p 1 -timeout 600s ../ci/e2e-group2/...

.PHONY: test-e2e-group-3
test-e2e-group-3: stop ## Runs group 3 e2e tests checking race conditions
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNEVENTDB)
	sleep 5
	$(RUNZKPROVER)
	docker ps -a
	docker logs $(DOCKERCOMPOSEZKPROVER)
	trap '$(STOP)' EXIT; MallocNanoZone=0 go test -count=1 -failfast -race -v -p 1 -timeout 600s ../ci/e2e-group3/...

.PHONY: test-e2e-group-4
test-e2e-group-4: stop ## Runs group 4 e2e tests checking race conditions
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNEVENTDB)
	sleep 5
	$(RUNZKPROVER)
	docker ps -a
	docker logs $(DOCKERCOMPOSEZKPROVER)
	trap '$(STOP)' EXIT; MallocNanoZone=0 go test -count=1 -failfast -race -v -p 1 -timeout 600s ../ci/e2e-group4/...

.PHONY: test-e2e-group-5
test-e2e-group-5: stop ## Runs group 5 e2e tests checking race conditions
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNEVENTDB)
	sleep 5
	$(RUNZKPROVER)
	docker ps -a
	docker logs $(DOCKERCOMPOSEZKPROVER)
	trap '$(STOP)' EXIT; MallocNanoZone=0 go test -count=1 -failfast -race -v -p 1 -timeout 1200s ../ci/e2e-group5/...

.PHONY: test-e2e-group-6
test-e2e-group-6: stop ## Runs group 6 e2e tests checking race conditions
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNEVENTDB)
	sleep 5
	$(RUNZKPROVER)
	docker ps -a
	docker logs $(DOCKERCOMPOSEZKPROVER)
	trap '$(STOP)' EXIT; MallocNanoZone=0 go test -count=1 -failfast -race -v -p 1 -timeout 600s ../ci/e2e-group6/...

.PHONY: test-e2e-group-7
test-e2e-group-7: stop ## Runs group 7 e2e tests checking race conditions
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNEVENTDB)
	sleep 5
	$(RUNZKPROVER)
	docker ps -a
	docker logs $(DOCKERCOMPOSEZKPROVER)
	trap '$(STOP)' EXIT; MallocNanoZone=0 go test -count=1 -failfast -race -v -p 1 -timeout 720s ../ci/e2e-group7/...

.PHONY: test-e2e-group-8
test-e2e-group-8: stop ## Runs group 8 e2e tests checking race conditions
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNEVENTDB)
	sleep 5
	$(RUNZKPROVER)
	docker ps -a
	docker logs $(DOCKERCOMPOSEZKPROVER)
	trap '$(STOP)' EXIT; MallocNanoZone=0 go test -count=1 -failfast -race -v -p 1 -timeout 1200s ../ci/e2e-group8/...


.PHONY: test-e2e-group-9
test-e2e-group-9: stop ## Runs group 9 e2e tests checking race conditions
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNEVENTDB)
	sleep 5
	$(RUNZKPROVER)
	docker ps -a
	docker logs $(DOCKERCOMPOSEZKPROVER)
	trap '$(STOP)' EXIT; MallocNanoZone=0 go test -count=1 -failfast -race -v -p 1 -timeout 2000s ../ci/e2e-group9/...


.PHONY: test-e2e-group-10
test-e2e-group-10: stop ## Runs group 10 e2e tests checking race conditions
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNEVENTDB)
	sleep 5
	$(RUNZKPROVER)
	docker ps -a
	docker logs $(DOCKERCOMPOSEZKPROVER)
	trap '$(STOP)' EXIT; MallocNanoZone=0 go test -count=1 -failfast -race -v -p 1 -timeout 2000s ../ci/e2e-group10/...


.PHONY: test-e2e-group-11
test-e2e-group-11: stop ## Runs group 11 e2e tests checking race conditions
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNEVENTDB)
	sleep 5
	$(RUNZKPROVER)
	docker ps -a
	docker logs $(DOCKERCOMPOSEZKPROVER)
	trap '$(STOP)' EXIT; MallocNanoZone=0 go test -count=1 -failfast -race -v -p 1 -timeout 2000s ../ci/e2e-group11/...

.PHONY: benchmark-sequencer-eth-transfers
benchmark-sequencer-eth-transfers: stop
	$(RUNL1NETWORK)
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNEVENTDB)
	sleep 5
	$(RUNZKPROVER)
	$(RUNSYNC)
	sleep 2
	$(RUNL2GASPRICER)
	$(RUNJSONRPC)
	docker ps -a
	docker logs $(DOCKERCOMPOSEZKPROVER)
	@ cd benchmarks/sequencer/e2e/eth-transfers ; \
 	mkdir -p results ; \
 	touch ./results/out.dat ; \
	go test -bench=. -timeout=600m | tee ./results/out.dat ;

.PHONY: benchmark-sequencer-erc20-transfers
benchmark-sequencer-erc20-transfers: stop
	$(RUNL1NETWORK)
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNEVENTDB)
	sleep 5
	$(RUNZKPROVER)
	$(RUNSYNC)
	sleep 2
	$(RUNL2GASPRICER)
	$(RUNJSONRPC)
	docker ps -a
	docker logs $(DOCKERCOMPOSEZKPROVER)
	@ cd benchmarks/sequencer/e2e/erc20-transfers ; \
 	mkdir -p results ; \
 	touch ./results/out.dat ; \
	go test -bench=. -timeout=600m | tee ./results/out.dat ;


.PHONY: benchmark-sequencer-uniswap-transfers
benchmark-sequencer-uniswap-transfers: stop
	$(RUNL1NETWORK)
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNEVENTDB)
	sleep 5
	$(RUNZKPROVER)
	$(RUNSYNC)
	sleep 2
	$(RUNL2GASPRICER)
	$(RUNJSONRPC)
	docker ps -a
	docker logs $(DOCKERCOMPOSEZKPROVER)
	@ cd benchmarks/sequencer/e2e/uniswap-transfers ; \
 	mkdir -p results ; \
 	touch ./results/out.dat ; \
	go test -bench=. -timeout=600m | tee ./results/out.dat ;

.PHONY: run-db
run-db: ## Runs the node database
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNEVENTDB)

.PHONY: stop-db
stop-db: ## Stops the node database
	$(STOPEVENTDB)
	$(STOPPOOLDB)
	$(STOPSTATEDB)

.PHONY: run-node
run-node: ## Runs the node
	$(RUNSYNC)
	sleep 4
	$(RUNETHTXMANAGER)
	$(RUNSEQUENCER)
	$(RUNSEQUENCESENDER)
	$(RUNL2GASPRICER)
	$(RUNAGGREGATOR)
	$(RUNJSONRPC)

.PHONY: stop-node
stop-node: ## Stops the node
	$(STOPSEQUENCER)
	$(STOPSEQUENCESENDER)
	$(STOPJSONRPC)
	$(STOPL2GASPRICER)
	$(STOPAGGREGATOR)
	$(STOPSYNC)
	$(STOPETHTXMANAGER)

.PHONY: run-network
run-network: ## Runs the l1 network
	$(RUNL1NETWORK)

.PHONY: run-network-v1tov2
run-network-v1tov2: ## Runs the v1tov2 l1 network
	$(RUNV1TOV2L1NETWORK)

.PHONY: stop-network
stop-network: ## Stops the l1 network
	$(STOPNETWORK)

.PHONY: stop-network-v1tov2
stop-network-v1tov2: ## Stops the l1 network
	$(STOPV1TOV2NETWORK)

.PHONY: run-zkprover
run-zkprover: ## Runs zkprover
	$(RUNZKPROVER)

.PHONY: stop-zkprover
stop-zkprover: ## Stops zkprover
	$(STOPZKPROVER)

.PHONY: run-l1-explorer
run-l1-explorer: ## Runs L1 blockscan explorer 
	$(RUNEXPLORERL1DB)
	$(RUNEXPLORERL1)

.PHONY: run-l2-explorer
run-l2-explorer: ## Runs L2 blockscan explorer
	$(RUNEXPLORERL2DB)
	$(RUNEXPLORERJSONRPC)
	$(RUNEXPLORERL2)

.PHONY: run-l2-explorer-json-rpc
run-l2-explorer-json-rpc: ## Runs L2 explorer json rpc
	$(RUNEXPLORERJSONRPC)

.PHONY: stop-l2-explorer-json-rpc
stop-l2-explorer-json-rpc: ## Stops L2 explorer json rpc
	$(STOPEXPLORERJSONRPC)

.PHONY: run-explorer
run-explorer: run-l1-explorer run-l2-explorer ## Runs both L1 and L2 explorers

.PHONY: stop-explorer
stop-explorer: ## Stops the explorer
	$(STOPEXPLORERL2)
	$(STOPEXPLORERL1)
	$(STOPEXPLORERJSONRPC)
	$(STOPEXPLORERL2DB)
	$(STOPEXPLORERL1DB)

.PHONY: run-explorer-db
run-explorer-db: ## Runs the explorer database
	$(RUNEXPLORERL1DB)
	$(RUNEXPLORERL2DB)

.PHONY: stop-explorer-db
stop-explorer-db: ## Stops the explorer database
	$(STOPEXPLORERL2DB)
	$(STOPEXPLORERL1DB)

.PHONY: run-seq
run-seq: ## runs the sequencer
	$(RUNSEQUENCER)

.PHONY: run-seq-v1tov2
run-seq-v1tov2: ## runs the sequencer
	$(RUNV1TOV2SEQUENCER)

.PHONY: stop-seq
stop-seq: ## stops the sequencer
	$(STOPSEQUENCER)

.PHONY: stop-seq-v1tov2
stop-seq-v1tov2: ## stops the sequencer
	$(STOPV1TOV2SEQUENCER)

.PHONY: run-seqsender
run-seqsender: ## runs the sequencer sender
	$(RUNSEQUENCESENDER)

.PHONY: run-seqsender-v1tov2
run-seqsender-v1tov2: ## runs the sequencer sender
	$(RUNV1TOV2SEQUENCESENDER)

.PHONY: stop-seqsender
stop-seqsender: ## stops the sequencer sender
	$(STOPSEQUENCESENDER)

.PHONY: stop-seqsender-v1tov2
stop-seqsender-v1tov2: ## stops the sequencer sender
	$(STOPV1TOV2SEQUENCESENDER)
	
.PHONY: run-sync
run-sync: ## runs the synchronizer
	$(RUNSYNC)

.PHONY: run-sync-v1tov2
run-sync-v1tov2: ## runs the synchronizer
	$(RUNV1TOV2SYNC)

.PHONY: stop-sync
stop-sync: ## stops the synchronizer
	$(STOPSYNC)

.PHONY: stop-sync-v1tov2
stop-sync-v1tov2: ## stops the synchronizer
	$(STOPV1TOV2SYNC)

.PHONY: run-json-rpc
run-json-rpc: ## runs the JSON-RPC
	$(RUNJSONRPC)

.PHONY: run-json-rpc-v1tov2
run-json-rpc-v1tov2: ## runs the JSON-RPC
	$(RUNV1TOV2JSONRPC)

.PHONY: stop-json-rpc
stop-json-rpc: ## stops the JSON-RPC
	$(STOPJSONRPC)

.PHONY: stop-json-rpc-v1tov2
stop-json-rpc-v1tov2: ## stops the JSON-RPC
	$(STOPV1TOV2JSONRPC)

.PHONY: run-l2gaspricer
run-l2gaspricer: ## runs the L2 Gas Price component
	$(RUNL2GASPRICER)

.PHONY: run-l2gaspricer-v1tov2
run-l2gaspricer-v1tov2: ## runs the L2 Gas Price component
	$(RUNV1TOV2L2GASPRICER)

.PHONY: stop-l2gaspricer
stop-l2gaspricer: ## stops the L2 Gas Price component
	$(STOPL2GASPRICER)

.PHONY: stop-l2gaspricer-v1tov2
stop-l2gaspricer-v1tov2: ## stops the L2 Gas Price component
	$(STOPV1TOV2L2GASPRICER)

.PHONY: run-eth-tx-manager
run-eth-tx-manager: ## Runs the eth tx manager service
	$(RUNETHTXMANAGER)

.PHONY: run-eth-tx-manager-v1tov2
run-eth-tx-manager-v1tov2: ## Runs the eth tx manager service
	$(RUNV1TOV2ETHTXMANAGER)

.PHONY: stop-eth-tx-manager
stop-eth-tx-manager: ## Stops the eth tx manager service
	$(STOPETHTXMANAGER)

.PHONY: stop-eth-tx-manager-v1tov2
stop-eth-tx-manager-v1tov2: ## Stops the eth tx manager service
	$(STOPV1TOV2ETHTXMANAGER)
	
.PHONY: run-agg
run-agg: ## Runs the aggregator service
	$(RUNAGGREGATOR)

.PHONY: run-agg-v1tov2
run-agg-v1tov2: ## Runs the aggregator service
	$(RUNV1TOV2AGGREGATOR)

.PHONY: stop-agg
stop-agg: ## Stops the aggregator service
	$(STOPAGGREGATOR)

.PHONY: stop-agg-v1tov2
stop-agg-v1tov2: ## Stops the aggregator service
	$(STOPV1TOV2AGGREGATOR)

.PHONY: run-grafana
run-grafana: ## Runs the grafana service
	$(RUNGRAFANA)

.PHONY: stop-grafana
stop-grafana: ## Stops the grafana service
	$(STOPGRAFANA)

.PHONY: run-permissionless
run-permissionless: run-node run-permissionless-dependencies ## Runs the trusted and permissionless node
	$(RUNPERMISSIONLESSDB)
	sleep 3
	$(RUNPERMISSIONLESSZKPROVER)
	$(RUNPERMISSIONLESSNODE)

.PHONY: stop-permissionless
stop-permissionless: stop-node stop-permissionless-dependencies ## Stops the permissionless node
	$(STOPPERMISSIONLESSNODE)


PHONY: run-permissionless-dependencies
run-permissionless-dependencies:  ## Runs the permissionless dependencies (db + prover) without the node
	$(RUNPERMISSIONLESSDB)
	sleep 3
	$(RUNPERMISSIONLESSZKPROVER)
	

PHONY: stop-permissionless-dependencies
stop-permissionless-dependencies:  ## Stop the permissionless dependencies (db + prover) without the node
	$(STOPPERMISSIONLESSNODE)
	$(STOPPERMISSIONLESSZKPROVER)
	$(STOPPERMISSIONLESSDB)

.PHONY: run-approve-pol
run-approve-pol: ## Runs approve in node container
	$(RUNAPPROVE)

.PHONY: run-approve-pol-v1tov2
run-approve-pol-v1tov2: ## Runs approve in node container
	$(RUNV1TOV2APPROVE)

.PHONY: stop-approve-pol
stop-approve-pol: ## Stops approve in node container
	$(STOPAPPROVE)

.PHONY: stop-approve-pol-v1tov2
stop-approve-pol-v1tov2: ## Stops approve in node container
	$(STOPV1TOV2APPROVE)

.PHONY: run
run: ## Runs a full node
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNEVENTDB)
	$(RUNL1NETWORK)
	sleep 1
	$(RUNZKPROVER)
	$(RUNAPPROVE)
	sleep 3
	$(RUNSYNC)
	sleep 4
	$(RUNETHTXMANAGER)
	$(RUNSEQUENCER)
	$(RUNSEQUENCESENDER)
	$(RUNL2GASPRICER)
	$(RUNAGGREGATOR)
	$(RUNJSONRPC)

.PHONY: run-v1tov2
run-v1tov2: ## Runs a full node using v1tov2 network
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNEVENTDB)
	$(RUNV1TOV2L1NETWORK)
	sleep 1
	$(RUNZKPROVER)
	$(RUNV1TOV2APPROVE)
	sleep 3
	$(RUNV1TOV2SYNC)
	sleep 4
	$(RUNV1TOV2ETHTXMANAGER)
	$(RUNV1TOV2SEQUENCER)
	$(RUNV1TOV2SEQUENCESENDER)
	$(RUNV1TOV2L2GASPRICER)
	$(RUNV1TOV2AGGREGATOR)
	$(RUNV1TOV2JSONRPC)

.PHONY: stop
stop: ## Stops all services
	$(STOP)

.PHONY: ship
ship: ## Builds docker images and run them
	cd .. && make build-docker && cd ./test && make run

.PHONY: reship
reship: stop ship ## Executes `make stop` and `make sail` commands

.PHONY: restart
restart: stop run ## Executes `make stop` and `make run` commands

.PHONY: run-metrics
run-metrics: ## Runs the metrics container
	$(RUNMETRICS)

.PHONY: stop-metrics
stop-metrics: ## Stops the metrics container
	$(STOPMETRICS)

.PHONY: init-network
init-network: ## Initializes the network
	go run ./scripts/init_network/main.go .

.PHONY: show-logs
show-logs:  ## Show logs for running docker 
	$(DOCKERCOMPOSE) logs

.PHONY: deploy-sc
deploy-sc: ## deploys some examples of transactions and smart contracts
	go run ./scripts/deploy_sc/main.go .

.PHONY: send-transfers
send-transfers: ## sends some ETH transfers txs to test the network
	go run ./scripts/send_transfers/main.go .

.PHONY: deploy-uniswap
deploy-uniswap: ## deploy the uniswap environment to the network
	go run ./scripts/uniswap/main.go .

.PHONY: run-db-scripts
run-db-scripts: ## Executes scripts on the db after it has been initialized, potentially using info from the environment
	./scripts/postgres/run.sh

.PHONY: install-mockery
install-mockery: ## Installs mockery with the correct version to generate the mocks
	go install github.com/vektra/mockery/v2@v2.39.0

.PHONY: generate-mocks
generate-mocks: generate-mocks-jsonrpc generate-mocks-sequencer generate-mocks-sequencesender generate-mocks-synchronizer generate-mocks-etherman generate-mocks-aggregator generate-mocks-state ## Generates mocks for the tests, using mockery tool

.PHONY: generate-mocks-jsonrpc
generate-mocks-jsonrpc: ## Generates mocks for jsonrpc , using mockery tool
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=storageInterface --dir=../jsonrpc --output=../jsonrpc --outpkg=jsonrpc --inpackage --structname=storageMock --filename=mock_storage.go
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=PoolInterface --dir=../jsonrpc/types --output=../jsonrpc/mocks --outpkg=mocks --structname=PoolMock --filename=mock_pool.go
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=StateInterface --dir=../jsonrpc/types --output=../jsonrpc/mocks --outpkg=mocks --structname=StateMock --filename=mock_state.go
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=EthermanInterface --dir=../jsonrpc/types --output=../jsonrpc/mocks --outpkg=mocks --structname=EthermanMock --filename=mock_etherman.go

.PHONY: generate-mocks-sequencer
generate-mocks-sequencer: ## Generates mocks for sequencer , using mockery tool
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=workerInterface --dir=../sequencer --output=../sequencer --outpkg=sequencer --inpackage --structname=WorkerMock --filename=mock_worker.go
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=stateInterface --dir=../sequencer --output=../sequencer --outpkg=sequencer --inpackage --structname=StateMock --filename=mock_state.go
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=txPool --dir=../sequencer --output=../sequencer --outpkg=sequencer --inpackage  --structname=PoolMock --filename=mock_pool.go
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=Tx --srcpkg=github.com/jackc/pgx/v4 --output=../sequencer --outpkg=sequencer --structname=DbTxMock --filename=mock_dbtx.go
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=ethermanInterface --dir=../sequencer --output=../sequencer --outpkg=sequencer --inpackage --structname=EthermanMock --filename=mock_etherman.go

.PHONY: generate-mocks-sequencesender
generate-mocks-sequencesender: ## Generates mocks for sequencesender , using mockery tool
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=stateInterface --dir=../sequencesender --output=../sequencesender --outpkg=sequencesender --inpackage --structname=StateMock --filename=mock_state.go
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=etherman --dir=../sequencesender --output=../sequencesender --outpkg=sequencesender --inpackage --structname=EthermanMock --filename=mock_etherman.go
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=ethTxManager --dir=../sequencesender --output=../sequencesender --outpkg=sequencesender --inpackage --structname=EthTxManagerMock --filename=mock_ethtxmanager.go


SYNC_L1_PARALLEL_FOLDER="../synchronizer/l1_parallel_sync"
SYNC_L1_PARALLEL_MOCKS_FOLDER="../synchronizer/l1_parallel_sync/mocks"
SYNC_L1_PARALLEL_PARAMS=--inpackage --outpkg=l1_parallel_sync
COMMON_MOCKERY_PARAMS=--disable-version-string --with-expecter
.PHONY: generate-mocks-synchronizer
generate-mocks-synchronizer: ## Generates mocks for synchronizer , using mockery tool
	## mocks for synchronizer
	#export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=EthermanInterface 		--dir=../synchronizer        			--output=../synchronizer --outpkg=synchronizer 	--structname=ethermanMock 						--filename=mock_etherman.go      ${COMMON_MOCKERY_PARAMS}
	#export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=stateInterface    		--dir=../synchronizer        			--output=../synchronizer --outpkg=synchronizer 	--structname=StateMock 							--filename=mock_state.go         ${COMMON_MOCKERY_PARAMS}
	#export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=ethTxManager       		--dir=../synchronizer       			--output=../synchronizer --outpkg=synchronizer 	--structname=ethTxManagerMock 					--filename=mock_ethtxmanager.go  ${COMMON_MOCKERY_PARAMS}
	#export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=poolInterface 			--dir=../synchronizer 					--output=../synchronizer --outpkg=synchronizer 	--structname=poolMock 							--filename=mock_pool.go          ${COMMON_MOCKERY_PARAMS}
	for i in l1RollupProducerInterface l1RollupConsumerInterface worker synchronizerProcessBlockRangeInterface workersInterface L1ParallelEthermanInterface; do \
		camelcase=$$(echo $$i | sed 's/\([a-z0-9]\)\([A-Z]\)/\1_\2/g' | tr '[:upper:]' '[:lower:]') ; \
		echo $$camelcase ; \
		export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=$$i --dir=../synchronizer/l1_parallel_sync --output=../synchronizer/l1_parallel_sync/mocks --outpkg=l1_parallel_sync --structname=$$i"Mock" --filename=mock_$$camelcase.go  --inpackage ${COMMON_MOCKERY_PARAMS}  ; \
	done

	rm -Rf ../synchronizer/l2_sync/l2_sync_etrog/mocks
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --all --case snake --dir ../synchronizer/l2_sync/l2_sync_etrog --output ../synchronizer/l2_sync/l2_sync_etrog/mocks --outpkg mock_l2_sync_etrog ${COMMON_MOCKERY_PARAMS}
	
	rm -Rf ../synchronizer/l2_sync/l2_shared/mocks
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --all --case snake --dir ../synchronizer/l2_sync/l2_shared --output ../synchronizer/l2_sync/l2_shared/mocks --outpkg mock_l2_shared ${COMMON_MOCKERY_PARAMS}
	
	rm -Rf ../synchronizer/common/syncinterfaces/mocks
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --all --case snake --dir ../synchronizer/common/syncinterfaces --output ../synchronizer/common/syncinterfaces/mocks --outpkg mock_syncinterfaces ${COMMON_MOCKERY_PARAMS}
	
	rm -Rf ../synchronizer/actions/elderberry/mocks
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --all --case snake --dir ../synchronizer/actions/elderberry --output ../synchronizer/actions/elderberry/mocks --outpkg mock_elderberry ${COMMON_MOCKERY_PARAMS}
	
	rm -Rf ../synchronizer/l1_check_block/mocks
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --all --case snake --dir ../synchronizer/l1_check_block --output ../synchronizer/l1_check_block/mocks --outpkg mock_l1_check_block ${COMMON_MOCKERY_PARAMS}
	
	

	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=Tx --srcpkg=github.com/jackc/pgx/v4 --output=../synchronizer/mocks --structname=DbTxMock --filename=mock_dbtx.go ${COMMON_MOCKERY_PARAMS}
	
.PHONY: generate-mocks-etherman	
generate-mocks-etherman: ## Generates mocks for etherman , using mockery tool
	## mocks for etherman
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=GasPricer --srcpkg=github.com/ethereum/go-ethereum --output=../etherman --outpkg=etherman --structname=etherscanMock --filename=mock_etherscan.go ${COMMON_MOCKERY_PARAMS}
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=GasPricer --srcpkg=github.com/ethereum/go-ethereum --output=../etherman --outpkg=etherman --structname=ethGasStationMock --filename=mock_ethgasstation.go ${COMMON_MOCKERY_PARAMS}

	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=ChainReader --srcpkg=github.com/ethereum/go-ethereum --output=../etherman --outpkg=etherman --structname=ChainReaderMock --filename=mock_chainreader.go ${COMMON_MOCKERY_PARAMS}


	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=ethermanInterface --dir=../ethtxmanager --output=../ethtxmanager --outpkg=ethtxmanager --structname=ethermanMock --filename=mock_etherman_test.go ${COMMON_MOCKERY_PARAMS}
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=stateInterface --dir=../ethtxmanager --output=../ethtxmanager --outpkg=ethtxmanager --structname=stateMock --filename=mock_state_test.go ${COMMON_MOCKERY_PARAMS}

	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=poolInterface --dir=../gasprice --output=../gasprice --outpkg=gasprice --structname=poolMock --filename=mock_pool.go ${COMMON_MOCKERY_PARAMS}
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=ethermanInterface --dir=../gasprice --output=../gasprice --outpkg=gasprice --structname=ethermanMock --filename=mock_etherman.go ${COMMON_MOCKERY_PARAMS}

	rm -Rf ../etherman/mockseth
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --all --case snake --dir ../etherman/ --output ../etherman/mockseth --outpkg mockseth ${COMMON_MOCKERY_PARAMS}
	

.PHONY: generate-mocks-aggregator	
generate-mocks-aggregator: ## Generates mocks for aggregator , using mockery tool
	## mocks for the aggregator tests
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=stateInterface --dir=../aggregator --output=../aggregator/mocks --outpkg=mocks --structname=StateMock --filename=mock_state.go
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=proverInterface --dir=../aggregator --output=../aggregator/mocks --outpkg=mocks --structname=ProverMock --filename=mock_prover.go
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=etherman --dir=../aggregator --output=../aggregator/mocks --outpkg=mocks --structname=Etherman --filename=mock_etherman.go
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=ethTxManager --dir=../aggregator --output=../aggregator/mocks --outpkg=mocks --structname=EthTxManager --filename=mock_ethtxmanager.go
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=aggregatorTxProfitabilityChecker --dir=../aggregator --output=../aggregator/mocks --outpkg=mocks --structname=ProfitabilityCheckerMock --filename=mock_profitabilitychecker.go
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=Tx --srcpkg=github.com/jackc/pgx/v4 --output=../aggregator/mocks --outpkg=mocks --structname=DbTxMock --filename=mock_dbtx.go

.PHONY: generate-mocks-state	
generate-mocks-state: ## Generates mocks for state , using mockery tool
	## mocks for the aggregator tests
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=storage --dir=../state --output=../state/mocks --outpkg=mocks --structname=StorageMock --filename=mock_storage.go --disable-version-string --with-expecter
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=ExecutorServiceClient --dir=../state/runtime/executor/ --output=../state/mocks --outpkg=mocks --structname=ExecutorServiceClientMock --filename=mock_executor_service_client.go --disable-version-string --with-expecter
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=Tx --srcpkg=github.com/jackc/pgx/v4 --output=../state/mocks --outpkg=mocks --structname=DbTxMock --filename=mock_dbtx.go --disable-version-string --with-expecter


.PHONY: run-benchmarks
run-benchmarks: run-db ## Runs benchmars
	go test -bench=. ./state/tree

.PHONY: compile-scs
compile-scs: ## Compiles smart contracts, configuration in test/contracts/index.yaml
	go run ./scripts/cmd... compilesc --input ./contracts

## Help display.
## Pulls comments from beside commands and prints a nicely formatted
## display with the commands and their usage information.
.DEFAULT_GOAL := help

.PHONY: help
help: ## Prints this help
		@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
