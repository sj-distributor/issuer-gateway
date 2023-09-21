package modules

import (
	"github.com/spf13/cobra"
	"issuer-gateway/grpc"
)

var GrpcCommand = &cobra.Command{
	Use:   "grpc",
	Short: "run grpc service",
	Run: func(cmd *cobra.Command, args []string) {
		conPath, _ := cmd.Flags().GetString("c")
		grpcServer.Run(conPath)
	},
}

func RunGrpcService(confPath string) {
	grpcServer.Run(confPath)
}
