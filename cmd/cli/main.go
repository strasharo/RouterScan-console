package main

import (
	"log"
	"net"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "routerscan",
		Commands: []*cli.Command{
			scanCommand(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func inet_aton(ip string) (ip_int uint32) {
	ip_byte := net.ParseIP(ip).To4()
	for i := 0; i < len(ip_byte); i++ {
		ip_int |= uint32(ip_byte[i])
		if i < 3 {
			ip_int <<= 8
		}
	}
	return
}
