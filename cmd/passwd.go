package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/chenglch/golang-xcat3client/utils"
	"github.com/spf13/cobra"
)

type ShowPasswdOptions struct {
	fields string
}

var (
	showPasswdOpts *ShowPasswdOptions
)

func ListPasswd(cmd *cobra.Command, args []string) {
	client, err := NewPasswdClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ret, err := client.Get("", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	passwdSlice := utils.InterfaceToSlice(ret.(map[string]interface{})["passwds"])
	if len(passwdSlice) == 0 {
		fmt.Println("Could not find any record")
		os.Exit(1)
	}
	for _, value := range passwdSlice {
		passwds := value.(map[string]interface{})
		fmt.Printf("%s (passwd)\n", passwds["key"])
	}
}

func ListPasswdCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List passwds(s) in xCAT3 service",
		Long:  `List passwds(s) in xCAT3 service. Format: list`,
		Run:   ListPasswd,
	}
	return cmd
}

func ShowPasswd(cmd *cobra.Command, args []string) {
	var fields []string
	var result interface{}
	if showPasswdOpts.fields != "" {
		fields = strings.Split(showPasswdOpts.fields, ",")
		if exist, _ := utils.Contains(fields, "key"); !exist {
			fields = append(fields, "key")
		}
	}

	client, err := NewPasswdClient()
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

func ShowPasswdCommand() *cobra.Command {
	showPasswdOpts = new(ShowPasswdOptions)
	cmd := &cobra.Command{
		Use:   "show <passwd key>",
		Short: "Show detailed infomation about passwds.",
		Long:  `Show detailed infomation about passwds. Format: show <passwd name>`,
		Run:   ShowPasswd,
	}
	cmd.Flags().StringVarP(&showPasswdOpts.fields, "fields", "i", "",
		`Fields seperated by comma. Only these fields will be fetched from the server.`)
	return cmd
}

func CreatePasswd(cmd *cobra.Command, args []string) {
	var result interface{}
	attr_map, err := utils.KeyValueArrayToMap(args[1:], "=")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	attr_map["key"] = args[0]
	client, err := NewPasswdClient()
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

func CreatePasswdCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create <passwd key> <key=val> [key=val]",
		Short: "Register passwds into xCAT3 service.",
		Long: `Register passwds into xCAT3 service. Format: create <passwd name> <key=val> [key=val]
		Current valied fields 'username', 'password', 'crypt_method'`,
		Run: CreatePasswd,
	}
	return cmd
}

func DeletePasswd(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("Please specify the uuid of passwds to delete")
		os.Exit(1)
	}

	client, err := NewPasswdClient()
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

func DeletePasswdCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <passwd key>",
		Short: "Unregister passwds from xCAT3 service.",
		Long:  `Unregister passwds from xCAT3 service. Format: delete <passwd key>`,
		Run:   DeletePasswd,
	}
	return cmd
}

func UpdatePasswd(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		fmt.Println("Please specify the name of passwds and attribute in key=val format to update")
		os.Exit(1)
	}
	patches := arg_array_to_patch(args[1:])
	client, err := NewPasswdClient()
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

func UpdatePasswdCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <passwd name> <key=val> [<key=val>]",
		Short: "Update information about registered passwds.",
		Long: `Update information about registered passwds. Format: update <passwd name> <key=val> [<key=val>]
		Current valied fields 'username', 'password', 'crypt_method'`,
		Run: UpdatePasswd,
	}
	return cmd
}

func PasswdCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "passwd",
		Short: "This is passwds child command for xcat3",
		Long: `xcat3 passwds --help and xcat3 passwds help COMMAND to see the usage for specfied
	command.`,
	}
	return cmd
}

func init() {
	PasswdCmd := PasswdCommand()
	PasswdCmd.AddCommand(ListPasswdCommand())
	PasswdCmd.AddCommand(ShowPasswdCommand())
	PasswdCmd.AddCommand(CreatePasswdCommand())
	PasswdCmd.AddCommand(DeletePasswdCommand())
	PasswdCmd.AddCommand(UpdatePasswdCommand())
	RootCmd.AddCommand(PasswdCmd)
}
