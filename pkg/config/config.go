package config

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/ardanlabs/conf"
	"github.com/eddie023/wex-tag/internal/build"
)

type ApiConfig struct {
	conf.Version
	Web struct {
		ReadTimeout     time.Duration `conf:"default:5s"`
		WriteTimeout    time.Duration `conf:"default:10s"`
		IdleTimeout     time.Duration `conf:"default:120s"`
		ShutdownTimeout time.Duration `conf:"default:20s"`
		APIHost         string        `conf:"default:0.0.0.0:8000,env:API_HOST"`
		DebugHost       string        `conf:"default:0.0.0.0:4000"`
	}

	Db struct {
		Host     string `conf:"default:localhost,env:DB_HOST"`
		Port     string `conf:"default:5432,env:DB_PORT"`
		User     string `conf:"default:user,env:DB_USER"`
		Dbname   string `conf:"default:wex_tag,env:DB_NAME"`
		Password string `conf:"default:pass,env:DB_PASSWORD"`
	}
}

func GetParsedConfig() (*ApiConfig, error) {
	cfg := ApiConfig{
		Version: conf.Version{
			SVN:  build.Build,
			Desc: "",
		},
	}

	help, err := conf.ParseOSArgs("API", &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil, err
		}
		slog.Error("unable to parse config", err)
		return nil, err
	}

	return &cfg, nil
}
