package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/dilap54/RouterScan-console/internal/routerscan"
	"github.com/urfave/cli/v2"
)

func scanCommand() *cli.Command {
	return &cli.Command{
		Name:  "scan",
		Usage: "run scanning routers in network",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "target",
				Usage: "<ip>:<port>",
				Value: "192.168.1.1:80",
			},
			&cli.PathFlag{
				Name:  "auth-basic",
				Value: "auth_basic.txt",
				Usage: "./path/to/<ip>\\t<port>\\n file with credentials dictionary for basic auth",
			},
			&cli.PathFlag{
				Name:  "auth-digest",
				Value: "auth_digest.txt",
				Usage: "./path/to/<ip>\\t<port>\\nfile with credentials dictionary for digest auth",
			},
			&cli.PathFlag{
				Name:  "auth-form",
				Value: "auth_form.txt",
				Usage: "./path/to/<ip>\\t<port>\\n file with credentials dictionary for form auth",
			},
			&cli.BoolFlag{
				Name:  "module-scanrouter",
				Value: true,
			},
			&cli.BoolFlag{
				Name:  "module-proxycheck",
				Value: false,
			},
			&cli.BoolFlag{
				Name:  "module-hnap",
				Value: true,
			},
			&cli.BoolFlag{
				Name:  "module-sqlite",
				Value: false,
			},
			&cli.BoolFlag{
				Name:  "module-hudson",
				Value: false,
			},
			&cli.BoolFlag{
				Name:  "module-phpmyadmin",
				Value: false,
			},
			&cli.BoolFlag{
				Name:  "st-enable-debug",
				Value: false,
			},
			&cli.IntFlag{
				Name:  "st-debug-verbosity",
				Value: 0,
				Usage: "supported levels: 1, 2, 3",
			},
			&cli.StringFlag{
				Name:  "st-user-agent",
				Value: "Mozilla/5.0 (Windows NT 5.1; rv:9.0.1) Gecko/20100101 Firefox/9.0.1",
			},
			&cli.PathFlag{
				Name:    "input",
				Aliases: []string{"i"},
				Value:   "-",
				Usage:   "/path/to/<ip>:<port>\n file with list of targets, or '-' for read targets from stdin",
			},
			&cli.IntFlag{
				Name:    "threads",
				Aliases: []string{"t"},
				Value:   5,
			},
		},
		Action: func(c *cli.Context) error {
			authBasic, err := ioutil.ReadFile(c.Path("auth-basic"))
			if err != nil {
				return err
			}
			authDigest, err := ioutil.ReadFile(c.Path("auth-digest"))
			if err != nil {
				return err
			}
			authForm, err := ioutil.ReadFile(c.Path("auth-form"))
			if err != nil {
				return err
			}
			moduleOptions := routerscan.ModulesOptions{
				Hnap:       c.Bool("module-hnap"),
				Hudson:     c.Bool("module-hudson"),
				Phpmyadmin: c.Bool("module-phpmyadmin"),
				Proxycheck: c.Bool("module-proxycheck"),
				Scanrouter: c.Bool("module-scanrouter"),
				Sqlite:     c.Bool("module-sqlite"),
			}
			stOptions := routerscan.StOptions{
				EnableDebug:    c.Bool("st-enable-debug"),
				DebugVerbosity: c.Int("st-debug-verbosity"),
				UserAgent:      c.String("st-user-agent"),
				PairsBasic:     string(authBasic),
				PairsDigest:    string(authDigest),
				PairsForm:      string(authForm),
			}
			rsscan, err := routerscan.New(moduleOptions, stOptions)
			if err != nil {
				return err
			}

			ch := make(chan string)
			wg := sync.WaitGroup{}
			for thread := 0; thread < c.Int("threads"); thread++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for in := range ch {
						target, err := routerscan.ParseTarget(in)
						if err != nil {
							log.Println(err)
							continue
						}
						result, err := rsscan.Scan(target)
						if err != nil {
							log.Println(err)
							continue
						}
						jsonBytes, err := json.Marshal(result)
						if err != nil {
							log.Println(err)
							continue
						}
						fmt.Println(string(jsonBytes))
					}
				}()
			}

			if c.IsSet("input") {
				var reader io.Reader
				if c.Path("input") == "-" {
					reader = os.Stdin
				} else {
					file, err := os.Open(c.Path("input"))
					if err != nil {
						return err
					}
					reader = file
				}
				scanner := bufio.NewScanner(reader)
				for scanner.Scan() {
					ch <- scanner.Text()
				}
			} else {
				ch <- c.String("target")
			}
			close(ch)

			wg.Wait()

			return nil
		},
	}
}
