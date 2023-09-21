package main

import (
	"fmt"
	"github.com/pygzfei/issuer-gateway/cmd/modules"
	"github.com/spf13/cobra"
	"log"
	"strings"
	"time"
)

var rootCmd = &cobra.Command{
	Use:   "ig",
	Short: "A High Performance Gateway with Certificate Issuance and Flexible Proxy",
	Long: `
··········································································································································
·                                                                                                                                        ·
·  ######## #######  #######  ##     ##  #######  #######          #######     ###   ########  #######  ##          ##    ###   ##    ## ·
·     ##    ##       ##       ##     ##  ##       ##    ##        ##          ## ##     ##     ##       ##   ###    ##   ## ##   ##  ##  ·
·     ##    #######  #######  ##     ##  ######   ## ##          ##   ####   ##   ##    ##     ######   ##  ##  ##  ##  ##   ##   ####   ·
·     ##         ##       ##  ##     ##  ##       ##  ##          ##    ##  ## ### ##   ##     ##        ## ##  ## ##  ## ### ##   ##    ·
·  ######## #######  #######   #######   #######  ##    ##         #######  ##     ##   ##     #######    ###    ###   ##     ##   ##    ·
·                                                                                                                                        ·
··········································································································································
                                                                                                                                          `,
	Run: func(cmd *cobra.Command, args []string) {
		mode, _ := cmd.Flags().GetString("mode")
		conf, _ := cmd.Flags().GetString("c")

		log.Println(fmt.Sprintf("Start engine: { config: %s, Mode: %s, args: %v }", conf, mode, args))

		if strings.ToUpper(mode) == "SINGLE" {
			tick := time.Tick(time.Second * 1)

			go modules.RunGrpcService(conf)
			go modules.RunIssuerService(conf)

			select {
			case <-tick:
				modules.RunGatewayService(conf)
			}
		}
	},
}

func main() {

	if err := rootCmd.Execute(); err != nil {
		log.Panic(err)
	}
}

func init() {
	rootCmd.Flags().StringP("mode", "m", "", "single or empty")
	rootCmd.Flags().StringP("c", "f", "", "Path to the configuration file")

	modules.GrpcCommand.Flags().StringP("c", "f", "", "Path to the configuration file")
	rootCmd.AddCommand(modules.GrpcCommand)

	modules.IssuerCommand.Flags().StringP("c", "f", "", "Path to the configuration file")
	rootCmd.AddCommand(modules.IssuerCommand)

	modules.GatewayCommand.Flags().StringP("c", "f", "", "Path to the configuration file")
	rootCmd.AddCommand(modules.GatewayCommand)
}
