package cmd

import (
	"github.com/hsukvn/go-mt4-tracker/server"
	"github.com/spf13/cobra"
)

var (
	debug bool
	port  int
)

var RootCmd = &cobra.Command{
	Use:   "mt4-tracker",
	Short: "API server to collect mt4 orders",
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := server.NewServer(&server.Config{
			Debug: debug,
		})
		if err != nil {
			return err
		}
		s.Run(port)
		return nil
	},
}

func init() {
	RootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "debug mode")
	RootCmd.Flags().IntVarP(&port, "port", "p", 9527, "port number")
}
