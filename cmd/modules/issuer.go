package modules

import (
	"errors"
	"github.com/pygzfei/issuer-gateway/issuer"
	"github.com/spf13/cobra"
)

var IssuerCommand = &cobra.Command{
	Use:   "issuer",
	Short: "run issuer service",
	RunE: func(cmd *cobra.Command, args []string) error {
		conPath, _ := cmd.Flags().GetString("c")
		if conPath == "" {
			return errors.New("config file is required! run: ig issuer -h")
		}
		issuer.Run(conPath)
		return nil
	},
}

func RunIssuerService(conPath string) {
	issuer.Run(conPath)
}
