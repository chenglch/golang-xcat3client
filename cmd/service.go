package cmd

import (
	"fmt"
	"os"

	"github.com/chenglch/golang-xcat3client/utils"
	"github.com/spf13/cobra"
)

func ListService(cmd *cobra.Command, args []string) {
	client, err := NewServiceClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ret, err := client.Get("", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	serviceSlice := utils.InterfaceToSlice(ret.(map[string]interface{})["services"])
	if len(serviceSlice) == 0 {
		fmt.Println("Could not find any record")
		os.Exit(1)
	}
	for _, value := range serviceSlice {
		service := value.(map[string]interface{})
		var online string
		var workers int
		hostname := service["hostname"].(string)
		if service["online"] == true {
			online = "online"
			workers = int(service["workers"].(float64))
			fmt.Printf("%s (hostname) %s %d (workers)\n", hostname, online, workers)

		} else {
			online = "offline"
			fmt.Printf("%s (hostname) %s\n", hostname, online)
		}
	}
}

func ListServiceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List service(s) in xCAT3 service",
		Long:  `List service(s) in xCAT3 service. Format: list`,
		Run:   ListService,
	}
	return cmd
}

func ShowService(cmd *cobra.Command, args []string) {
	var result interface{}

	if len(args) != 1 {
		fmt.Println("Please specify the service name.")
		os.Exit(1)
	}
	client, err := NewServiceClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	result, err = client.Show(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	utils.PrintJson(result)
}

func ShowServiceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show <service hostname>",
		Short: "Show detailed infomation about service.",
		Long:  `Show detailed infomation about service. Format: show <service hostname>`,
		Run:   ShowService,
	}
	return cmd
}

func ServiceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service",
		Short: "This is service child command for xcat3",
		Long: `xcat3 service --help and xcat3 service help COMMAND to see the usage for specfied
	command.`,
	}
	return cmd
}

func init() {
	ServiceCmd := ServiceCommand()
	ServiceCmd.AddCommand(ListServiceCommand())
	ServiceCmd.AddCommand(ShowServiceCommand())
	RootCmd.AddCommand(ServiceCmd)
}
