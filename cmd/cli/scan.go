package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/dilap54/RouterScan-console/pkg/routerscan"
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
		},
		Action: func(c *cli.Context) error {
			if err := routerscan.Initialize(); err != nil {
				return err
			}
			count, err := routerscan.GetModuleCount()
			if err != nil {
				return err
			}
			for i := 0; i < count; i++ {
				info, err := routerscan.GetModuleInfo(i)
				if err != nil {
					return err
				}
				if strings.Contains(info.Name, "RouterScanRouter") && info.Enabled != c.Bool("module-scanrouter") {
					if err := routerscan.SwitchModule(i, c.Bool("module-scanrouter")); err != nil {
						log.Fatalf("cannot switch module RouterScanRouter to %t", c.Bool("module-scanrouter"))
					} else {
						log.Printf("module RouterScanRouter enabled = %t", c.Bool("module-scanrouter"))
					}
				}
				if strings.Contains(info.Name, "ProxyCheckDetect") && info.Enabled != c.Bool("module-proxycheck") {
					if err := routerscan.SwitchModule(i, c.Bool("module-proxycheck")); err != nil {
						log.Fatalf("cannot switch module ProxyCheckDetect to %t", c.Bool("module-proxycheck"))
					} else {
						log.Printf("module ProxyCheckDetect enabled = %t", c.Bool("module-proxycheck"))
					}
				}
				if strings.Contains(info.Name, "HNAPUse") && info.Enabled != c.Bool("module-hnap") {
					if err := routerscan.SwitchModule(i, c.Bool("module-hnap")); err != nil {
						log.Fatalf("cannot switch module HNAPUse to %t", c.Bool("module-hnap"))
					} else {
						log.Printf("module HNAPUse enabled = %t", c.Bool("module-hnap"))
					}
				}
				if strings.Contains(info.Name, "SQLiteSQLite") && info.Enabled != c.Bool("module-sqlite") {
					if err := routerscan.SwitchModule(i, c.Bool("module-sqlite")); err != nil {
						log.Fatalf("cannot switch module SQLiteSQLite to %t", c.Bool("module-sqlite"))
					} else {
						log.Printf("module SQLiteSQLite enabled = %t", c.Bool("module-sqlite"))
					}
				}
				if strings.Contains(info.Name, "HudsonHudson") && info.Enabled != c.Bool("module-hudson") {
					if err := routerscan.SwitchModule(i, c.Bool("module-hudson")); err != nil {
						log.Fatalf("cannot switch module HudsonHudson to %t", c.Bool("module-hudson"))
					} else {
						log.Printf("module HudsonHudson enabled = %t", c.Bool("module-hudson"))
					}
				}
				if strings.Contains(info.Name, "PMAphpMyAdmin") && info.Enabled != c.Bool("module-phpmyadmin") {
					if err := routerscan.SwitchModule(i, c.Bool("module-phpmyadmin")); err != nil {
						log.Fatalf("cannot switch module PMAphpMyAdmin to %t", c.Bool("module-phpmyadmin"))
					} else {
						log.Printf("module PMAphpMyAdmin enabled = %t", c.Bool("module-phpmyadmin"))
					}
				}
			}
			if c.IsSet("st-enable-debug") {
				log.Printf("setting stEnableDebug = %t", c.Bool("st-enable-debug"))
				if err := routerscan.SetParamBool(routerscan.StEnableDebug, c.Bool("st-enable-debug")); err != nil {
					panic(err)
				}
			}

			if err := routerscan.SetParamString(routerscan.StUserAgent, "Mozilla/5.0 (Windows NT 5.1; rv:9.0.1) Gecko/20100101 Firefox/9.0.1"); err != nil {
				panic(err)
			}
			if err := routerscan.SetParamString(routerscan.StPairsBasic, "admin\tadmin"); err != nil {
				panic(err)
			}
			if err := routerscan.SetParamString(routerscan.StPairsDigest, "admin\tadmin"); err != nil {
				panic(err)
			}
			if err := routerscan.SetParamString(routerscan.StPairsForm, "admin\tadmin"); err != nil {
				panic(err)
			}

			if err := routerscan.SetSetTableDataCallback(func(row uint, name string, value string) {
				fmt.Printf("%d %s %s\n", row, name, value)
			}); err != nil {
				panic(err)
			}
			if err := routerscan.SetWriteLogCallback(func(str string, verbosity int) {
				fmt.Printf("%s %d\n", str, verbosity)
			}); err != nil {
				panic(err)
			}
			host := strings.Split(c.String("target"), ":")
			port, _ := strconv.ParseUint(host[1], 10, 16)
			router, err := routerscan.PrepareRouter(1, inet_aton(host[0]), uint16(port))
			if err != nil {
				panic(err)
			}
			if err := router.Scan(); err != nil {
				panic(err)
			}
			if err := router.Free(); err != nil {
				panic(err)
			}
			return nil
		},
	}
}
