package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	"github.com/ardanlabs/conf"
	"github.com/eddie023/wex-tag/ent"
	"github.com/eddie023/wex-tag/ent/migrate"
	"github.com/eddie023/wex-tag/ent/migrate/migrations"

	gomigrate "github.com/golang-migrate/migrate/v4"

	"github.com/eddie023/wex-tag/pkg/config"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"

	"github.com/urfave/cli/v2"

	atlas "ariga.io/atlas/sql/migrate"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"
)

var CreateCommand = cli.Command{
	Name: "create",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "name", Required: true},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		err := godotenv.Load()
		if err != nil {
			return err
		}

		cfg, err := config.GetParsedConfig()
		if err != nil {
			return err
		}

		err = conf.Parse(os.Args, "", cfg)
		if err != nil {
			return err
		}

		connectionURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.Db.User, cfg.Db.Password, cfg.Db.Host, cfg.Db.Port, cfg.Db.Dbname)

		path := "ent/migrate/migrations"
		// Create a local migration directory able to understand Atlas migration file format for replay.
		if err := os.MkdirAll(path, 0755); err != nil {
			log.Fatalf("creating migration directory: %v", err)
		}
		dir, err := atlas.NewLocalDir(path)
		if err != nil {
			log.Fatalf("failed creating atlas migration directory: %v", err)
		}
		// Migrate diff options.
		opts := []schema.MigrateOption{
			schema.WithDir(dir),                          // provide migration directory
			schema.WithMigrationMode(schema.ModeInspect), // provide migration mode
			schema.WithDialect(dialect.Postgres),
			schema.WithDropColumn(true),
			schema.WithDropIndex(true),
		}

		fmt.Println("the connection usr is", connectionURL)

		err = migrate.NamedDiff(ctx, connectionURL, c.String("name"), opts...)
		if err != nil {
			log.Fatalf("failed named diff with err: %s", err)
		}

		fmt.Println("successfully created new migration file")

		return nil
	},
}

// migrate
var UpCommand = cli.Command{
	Name: "up",
	Action: func(ctx *cli.Context) error {

		err := godotenv.Load()
		if err != nil {
			return err
		}

		cfg, err := config.GetParsedConfig()
		if err != nil {
			return err
		}

		err = conf.Parse(os.Args, "", cfg)
		if err != nil {
			return err
		}

		connectionURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.Db.User, cfg.Db.Password, cfg.Db.Host, cfg.Db.Port, cfg.Db.Dbname)

		fmt.Printf("using postgres connection url '%s'", connectionURL)

		d, err := migrations.MigrationFS()
		if err != nil {
			return errors.WithMessage(err, "iofs.New")
		}

		m, err := gomigrate.NewWithSourceInstance("iofs", d, connectionURL)
		if err != nil {
			return errors.WithMessage(err, "NewWithSourceInstance")
		}
		if err := m.Up(); err != nil {
			return errors.WithMessage(err, "m.up")
		}

		fmt.Println("successfully completed migration up")

		return nil
	},
}

var InitCommand = cli.Command{
	Name: "init",
	Action: func(ctx *cli.Context) error {
		return Migrate(ctx.Context)
	},
}

var Migration = cli.Command{
	Name:        "migrate",
	Subcommands: []*cli.Command{&CreateCommand, &InitCommand, &UpCommand},
}

func main() {
	app := &cli.App{
		Name: "dev",
		Commands: []*cli.Command{
			&Migration,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func Migrate(ctx context.Context) error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	cfg, err := config.GetParsedConfig()
	if err != nil {
		return err
	}

	err = conf.Parse(os.Args, "", cfg)
	if err != nil {
		return err
	}

	connectionURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", cfg.Db.Host, cfg.Db.Port, cfg.Db.User, cfg.Db.Dbname, cfg.Db.Password)

	fmt.Printf("using postgres connection url '%s'", connectionURL)

	client, err := ent.Open("postgres", connectionURL)
	if err != nil {
		return err
	}

	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	return nil
}
