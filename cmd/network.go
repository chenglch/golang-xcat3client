package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/chenglch/golang-xcat3client/utils"
	"github.com/spf13/cobra"
)

type ShowNetworkOptions struct {
	fields string
}

var (
	showNetworkOpts *ShowNetworkOptions
)

func ListNetwork(cmd *cobra.Command, args []string) {
	client, err := NewNetworkClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ret, err := client.Get("", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	networkSlice := utils.InterfaceToSlice(ret.(map[string]interface{})["networks"])
	if len(networkSlice) == 0 {
		fmt.Println("Could not find any record")
		os.Exit(1)
	}
	for _, value := range networkSlice {
		network := value.(map[string]interface{})
		fmt.Printf("%s (network)\n", network["name"])
	}
}

func ListNetworkCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List network(s) in xCAT3 service",
		Long:  `List network(s) in xCAT3 service. Format: list`,
		Run:   ListNetwork,
	}
	return cmd
}

func ShowNetwork(cmd *cobra.Command, args []string) {
	var fields []string
	var result interface{}
	if showNetworkOpts.fields != "" {
		fields = strings.Split(showNetworkOpts.fields, ",")
		if exist, _ := utils.Contains(fields, "name"); !exist {
			fields = append(fields, "name")
		}
	}

	client, err := NewNetworkClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(args) == 1 {
		result, err = client.Show(args[0], fields, nil, true)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	utils.PrintJson(result)
}

func ShowNetworkCommand() *cobra.Command {
	showNetworkOpts = new(ShowNetworkOptions)
	cmd := &cobra.Command{
		Use:   "show <network name>",
		Short: "Show detailed infomation about network.",
		Long:  `Show detailed infomation about network. Format: show <network name>`,
		Run:   ShowNetwork,
	}
	cmd.Flags().StringVarP(&showNetworkOpts.fields, "fields", "i", "",
		`Fields seperated by comma. Only these fields will be fetched from the server.`)
	return cmd
}

func CreateNetwork(cmd *cobra.Command, args []string) {
	var result interface{}
	attr_map, err := utils.KeyValueArrayToMap(args[1:], "=")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	attr_map["name"] = args[0]
	client, err := NewNetworkClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(attr_map) > 1 {
		result, err = client.Post("", attr_map, true)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Pleace specify the attribute key values in key1=val1 key2=val2 format")
		os.Exit(1)
	}
	utils.PrintJson(result)
}

func CreateNetworkCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create <network name> <key=val> [key=val]",
		Short: "Register network into xCAT3 service.",
		Long: `Register network into xCAT3 service. Format: create <network name> <key=val> [key=val]
		Current valid fields 'subnet', 'netmask', 'gateway', 'dhcpserver', 'dynamic_range',
                'nameservers', 'domain'`,
		Run: CreateNetwork,
	}
	return cmd
}

func DeleteNetwork(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("Please specify the uuid of network to delete")
		os.Exit(1)
	}

	client, err := NewNetworkClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = client.Delete(args[0], nil, true)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%s: deleted\n", args[0])
}

func DeleteNetworkCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <network name>",
		Short: "Unregister network from xCAT3 service.",
		Long:  `Unregister network from xCAT3 service. Format: delete <network name>`,
		Run:   DeleteNetwork,
	}
	return cmd
}

func UpdateNetwork(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		fmt.Println("Please specify the name of network and attribute in key=val format to update")
		os.Exit(1)
	}
	patches := arg_array_to_patch(args[1:])
	client, err := NewNetworkClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = client.Patch(args[0], patches, true)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%s: updated\n", args[0])
}

func UpdateNetworkCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <network name> <key=val> [<key=val>]",
		Short: "Update information about registered network.",
		Long: `Update information about registered network. Format: update <network name> <key=val> [<key=val>]
		Current valid fields 'subnet', 'netmask', 'gateway', 'dhcpserver', 'dynamic_range',
                'nameservers', 'domain'`,
		Run: UpdateNetwork,
	}
	return cmd
}

func NetworkCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network",
		Short: "This is network child command for xcat3",
		Long: `xcat3 network --help and xcat3 network help COMMAND to see the usage for specfied
	command.`,
	}
	return cmd
}

func init() {
	NetworkCmd := NetworkCommand()
	NetworkCmd.AddCommand(ListNetworkCommand())
	NetworkCmd.AddCommand(ShowNetworkCommand())
	NetworkCmd.AddCommand(CreateNetworkCommand())
	NetworkCmd.AddCommand(DeleteNetworkCommand())
	NetworkCmd.AddCommand(UpdateNetworkCommand())
	RootCmd.AddCommand(NetworkCmd)
}
