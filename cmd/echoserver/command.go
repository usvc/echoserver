package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	command   *cobra.Command
	Commit    string
	Version   string
	Timestamp string
)

func GetCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "echoserver",
		Version: fmt.Sprintf("%s-%s / %s", Version, Commit, Timestamp),
		Run: func(cmd *cobra.Command, args []string) {
			bindAddress := fmt.Sprintf("%s:%v", conf.GetString("server_addr"), conf.GetUint("server_port"))
			log.Debugf("server binding to address '%s'", bindAddress)
			server := getServer(bindAddress)
			if err := server.ListenAndServe(); err != nil {
				log.Errorf("failed to start the http server: '%s'", err)
				os.Exit(1)
			}
		},
	}
}
