package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/chenglch/golang-xcat3client/utils"
	"github.com/spf13/cobra"
)

type ShowNicOptions struct {
	fields string
	mac    string
}

var (
	showNicOpts  *ShowNicOptions
	VALID_FIELDS = []string{"uuid", "mac", "name", "ip", "netmask", "extra", "node"}
)

func ListNics(cmd *cobra.Command, args []string) {
	client, err := NewNicClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ret, err := client.Get("", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	nicSlice := utils.InterfaceToSlice(ret.(map[string]interface{})["nics"])
	if len(nicSlice) == 0 {
		fmt.Println("Could not find any record")
		os.Exit(1)
	}
	for _, value := range nicSlice {
		nic := value.(map[string]interface{})
		fmt.Printf("%s (uuid) %s (mac)\n", nic["uuid"], nic["mac"])
	}
}

func ListNicCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List nic(s) in xCAT3 service",
		Long:  `List nic(s) in xCAT3 service. Format: list`,
		Run:   ListNics,
	}
	return cmd
}

func ShowNic(cmd *cobra.Command, args []string) {
	var fields []string
	var result interface{}
	if showNicOpts.fields != "" {
		fields = strings.Split(showNicOpts.fields, ",")
		if exist, _ := utils.Contains(fields, "uuid"); !exist {
			fields = append(fields, "uuid")
		}
	}

	client, err := NewNicClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(args) == 1 {
		result, err = client.Show(args[0], fields)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		if showNicOpts.mac != "" {
			result, err = client.GetByMac(showNicOpts.mac, fields)
		}
	}
	utils.PrintJson(result)
}

func ShowNicCommand() *cobra.Command {
	showNicOpts = new(ShowNicOptions)
	cmd := &cobra.Command{
		Use:   "show <nic uuid>",
		Short: "Show detailed infomation about nic.",
		Long:  `Show detailed infomation about nic. Format show <nic uuid>, show --mac <mac>`,
		Run:   ShowNic,
	}
	cmd.Flags().StringVarP(&showNicOpts.fields, "fields", "i", "",
		`Fields seperated by comma. Only these fields will be fetched from the server.`)
	cmd.Flags().StringVarP(&showNicOpts.mac, "mac", "", "",
		`Search nic by mac address.`)
	return cmd
}

func CreateNics(cmd *cobra.Command, args []string) {
	var result interface{}
	attr_map, err := utils.KeyValueArrayToMap(args, "=")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	hasMac := false
	hasNode := false
	for k, _ := range attr_map {
		if exist, _ := utils.Contains(VALID_FIELDS, k); !exist {
			fmt.Printf("Only allow attributes '%s'. '%s' is given.\n", strings.Join(VALID_FIELDS, " "), k)
			os.Exit(1)
		}
		if k == "mac" {
			hasMac = true
		} else if k == "node" {
			hasNode = true
		}
	}
	if !hasMac || !hasNode {
		fmt.Println("Please specify the 'node' and 'mac' attributes")
		os.Exit(1)
	}
	client, err := NewNicClient()
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

func CreateNicCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create <key=val> [<key=val>]",
		Short: "Register nic into xCAT3 service.",
		Long: `Register nic into xCAT3 service. Format: create <key=val> [<key=val>]. Can be specified multiple times.
		Current valid fields uuid,mac,name,ip,netmask,extra,node`,
		Run: CreateNics,
	}
	return cmd
}

func DeleteNic(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("Please specify the uuid of nic to delete")
		os.Exit(1)
	}

	client, err := NewNicClient()
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

func DeleteNicCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <nic uuid>",
		Short: "Unregister nic from xCAT3 service.",
		Long:  `Unregister nic from xCAT3 service. Format: delete <nic uuid>`,
		Run:   DeleteNic,
	}
	return cmd
}

func UpdateNic(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		fmt.Println("Please specify the uuid of nic and attribute in key=val format to update")
		os.Exit(1)
	}
	patches := arg_array_to_patch(args[1:])
	client, err := NewNicClient()
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

func UpdateNicCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <nic uuid> <path=value> [<path=value>]",
		Short: "Update information about registered nic.",
		Long: `Update information about registered nic.
		Format: update <nic uuid> <path=value> [<path=value>].
		Can be specified multiple times. Current valid fields uuid,mac,name,ip,netmask,extra,node`,
		Run: UpdateNic,
	}
	return cmd
}

func NicCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nic",
		Short: "This is nic child command for xcat3",
		Long: `xcat3 nic --help and xcat3 nic help COMMAND to see the usage for specfied
	command.`,
	}
	return cmd
}

func init() {
	NicCmd := NicCommand()
	NicCmd.AddCommand(ListNicCommand())
	NicCmd.AddCommand(ShowNicCommand())
	NicCmd.AddCommand(CreateNicCommand())
	NicCmd.AddCommand(DeleteNicCommand())
	NicCmd.AddCommand(UpdateNicCommand())
	RootCmd.AddCommand(NicCmd)
}
