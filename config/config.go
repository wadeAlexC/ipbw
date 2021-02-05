package config

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/wadeAlexC/ipbw/crawler"
	"github.com/wadeAlexC/ipbw/crawler/modules"
	"go.uber.org/fx"
)

type Config struct {
	Modules []fx.Option
	Invokes []fx.Option

	hello []string
}

const (
	// crawler
	FlagNumWorkers    = "num-workers"
	FlagCrawlDuration = "crawl-duration"

	// modules/status
	FlagEnableStatus   = "enable-status"
	FlagStatusInterval = "status-interval"

	// modules/identifier
	FlagEnableIdentifier = "enable-identifier"
)

// Default returns a Config that sets up the crawler
// - Provide the NewCrawler constructor to the app
// - Invoke SetNumWorkers on app startup
func Default(cctx *cli.Context) *Config {
	config := &Config{
		Modules: []fx.Option{fx.Provide(crawler.NewCrawler)},
		Invokes: []fx.Option{fx.Invoke(func(c *crawler.Crawler) error {
			return c.SetNumWorkers(cctx.Uint(FlagNumWorkers))
		})},
		hello: []string{
			"IPBW config:\n",
			"====================\n",
			fmt.Sprintf("number of workers: %d\n", cctx.Uint(FlagNumWorkers)),
			fmt.Sprintf("crawl duration: %d min\n", cctx.Uint(FlagCrawlDuration)),
			"====================\n",
		},
	}

	return config
}

// ConfigStatus adds setup for modules/status, if enabled
func (c *Config) ConfigStatus(cctx *cli.Context) {
	if cctx.Bool(FlagEnableStatus) {
		c.Modules = append(c.Modules, fx.Provide(modules.NewStatus))
		c.Invokes = append(c.Invokes, fx.Invoke(func(s *modules.Status) error {
			return s.SetInterval(cctx.Uint(FlagStatusInterval))
		}))
		c.hello = append(c.hello,
			"modules/status: enabled\n",
			fmt.Sprintf("- status interval: %d min\n", cctx.Uint(FlagStatusInterval)),
		)
	}
}

// ConfigIdentifier adds setup for modules/identifier, if enabled
func (c *Config) ConfigIdentifier(cctx *cli.Context) {
	if cctx.Bool(FlagEnableIdentifier) {
		c.Modules = append(c.Modules, fx.Provide(modules.NewIdentifier))
		c.Invokes = append(c.Invokes, fx.Invoke(func(i *modules.Identifier) error {
			return i.Setup()
		}))
		c.hello = append(c.hello,
			"modules/identifier: enabled\n",
		)
	}
}

// Hello prints a summary of config to the console
func (c *Config) Hello() {
	hello := strings.Join(c.hello, "")
	fmt.Printf(hello)
	fmt.Printf("====================\n")
}
