package flow

import (
	"context"
	"flag"
	"github.com/adverax/configs"
	"github.com/adverax/core"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

func Main(
	config interface{},
	action func(ctx context.Context),
) {
	configFiles := flag.String("config", "./configs/local.yml", "Comma-separated list of configuration files")
	flag.Parse()

	err := configs.LoadYamlConfigs(config, makeConfigFiles(*configFiles)...)
	if err != nil {
		panic(err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-sigs
		cancel()
	}()

	action(ctx)
}

func makeConfigFiles(conf string) []string {
	files := strings.Split(conf, ",")
	for i, file := range files {
		files[i] = core.Must(filepath.Abs(file))
	}
	return files
}
