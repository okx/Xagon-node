package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/log"
	xl_pool "github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/pool/pgpoolstorage"
	"github.com/0xPolygonHermez/zkevm-node/test/pendingtx/pool"
	"github.com/0xPolygonHermez/zkevm-node/test/pendingtx/sequencer"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

var (
	levelCount   = 50
	addrPerLevel = 400
	txPerAddr    = 5
)

func main() {
	c, err := loadConfig()
	if err != nil {
		log.Error(err)
		return
	}
	log.Init(c.Log)
	//var cancelFuncs []context.CancelFunc
	ctx := context.Background()

	poolInstance := createPool(c.Pool)
	start := time.Now()
	poolInstance.PrepareTx(ctx, levelCount, addrPerLevel, txPerAddr)
	cost := time.Since(start)

	time.Sleep(3 * time.Second)
	count, err := poolInstance.CountPendingTransactions(ctx)
	if err != nil {
		log.Error(err)
		return
	}
	fmt.Println("===================")
	fmt.Printf("prepared tx count:%d. time cost:%s\n", count, cost)
	fmt.Println("===================")

	finishedCh := make(chan int, 1)
	seq := createSequencer(*c, poolInstance)
	go seq.Start(ctx, levelCount, addrPerLevel*txPerAddr, finishedCh)

	//time.Sleep(5 * time.Second)
	go poolInstance.Speed(ctx)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	for {
		select {
		case <-sigCh:
			log.Info("terminating application gracefully...")
			os.Exit(0)
		case level := <-finishedCh:
			fmt.Println("===================")
			fmt.Printf("Speed level [%d] finished.\n", level)
			fmt.Println("===================")
			if level == 0 {
				os.Exit(0)
			}
		}
	}
}

func createSequencer(cfg config.Config, pool *pool.Pool) *sequencer.Sequencer {
	seq, err := sequencer.New(cfg.Sequencer, cfg.State.Batch, pool)
	if err != nil {
		log.Fatal(err)
	}
	return seq
}

func createPool(cfgPool xl_pool.Config) *pool.Pool {
	runPoolMigrations(cfgPool.DB)
	poolStorage, err := pgpoolstorage.NewPostgresPoolStorage(cfgPool.DB)
	if err != nil {
		log.Fatal(err)
	}
	poolInstance := pool.NewPool(cfgPool, poolStorage)
	return poolInstance
}

func runPoolMigrations(c db.Config) {
	runMigrations(c, db.PoolMigrationName)
}

func runMigrations(c db.Config, name string) {
	log.Infof("running migrations for %v", name)
	err := db.RunMigrationsUp(c, name)
	if err != nil {
		log.Fatal(err)
	}
}

func waitSignal(cancelFuncs []context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	for sig := range signals {
		switch sig {
		case os.Interrupt, os.Kill:
			log.Info("terminating application gracefully...")

			exitStatus := 0
			for _, cancel := range cancelFuncs {
				cancel()
			}
			os.Exit(exitStatus)
		}
	}
}

// Load loads the configuration
func loadConfig() (*config.Config, error) {
	cfg, err := config.Default()
	if err != nil {
		return nil, err
	}
	configFilePath := "./config.toml"
	dirName, fileName := filepath.Split(configFilePath)

	fileExtension := strings.TrimPrefix(filepath.Ext(fileName), ".")
	fileNameWithoutExtension := strings.TrimSuffix(fileName, "."+fileExtension)

	viper.AddConfigPath(dirName)
	viper.SetConfigName(fileNameWithoutExtension)
	viper.SetConfigType(fileExtension)

	err = viper.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if ok {
			log.Infof("config file not found")
		} else {
			log.Infof("error reading config file: ", err)
			return nil, err
		}
	}

	decodeHooks := []viper.DecoderConfigOption{
		// this allows arrays to be decoded from env var separated by ",", example: MY_VAR="value1,value2,value3"
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(mapstructure.TextUnmarshallerHookFunc(), mapstructure.StringToSliceHookFunc(","))),
	}

	err = viper.Unmarshal(&cfg, decodeHooks...)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
