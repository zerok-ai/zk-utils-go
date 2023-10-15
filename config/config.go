package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	zklogger "github.com/zerok-ai/zk-utils-go/logs"
	"os"
)

var LogTag = "zk_config"

type Args struct {
	ConfigPath string
}

// ProcessArgs processes and handles CLI arguments
func ProcessArgs[T any](cfg *T) error {
	zklogger.DebugF(LogTag, "reading configs")
	var a Args

	flagSet := flag.NewFlagSet("server", 1)
	flagSet.StringVar(&a.ConfigPath, "c", "config.yaml", "Path to configuration file")

	fu := flagSet.Usage
	flagSet.Usage = func() {
		fu()
		envHelp, _ := cleanenv.GetDescription(cfg, nil)
		if _, err := fmt.Fprintln(flagSet.Output()); err != nil {
			return
		}

		_, err := fmt.Fprintln(flagSet.Output(), envHelp)
		if err != nil {
			return
		}
	}

	if err := flagSet.Parse(os.Args[1:]); err != nil {
		return err
	}

	return cleanenv.ReadConfig(a.ConfigPath, cfg)
}
