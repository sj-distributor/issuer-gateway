package modules

import (
	"github.com/spf13/cobra"
	"issuer-gateway/gateway"
)

var GatewayCommand = &cobra.Command{
	Use:   "gateway",
	Short: "run gateway service",
	Run: func(cmd *cobra.Command, args []string) {
		conPath, _ := cmd.Flags().GetString("c")
		gateway.Run(conPath)
	},
}

func RunGatewayService(conPath string) {
	gateway.Run(conPath)
}
