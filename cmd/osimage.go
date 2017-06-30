package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/chenglch/golang-xcat3client/utils"
	"github.com/spf13/cobra"
)

type ShowOsimageOptions struct {
	fields string
}

var (
	showOsimageOpts *ShowOsimageOptions
)

func ListOsimage(cmd *cobra.Command, args []string) {
	client, err := NewOsimageClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ret, err := client.Get("", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	osimageSlice := utils.InterfaceToSlice(ret.(map[string]interface{})["images"])
	if len(osimageSlice) == 0 {
		fmt.Println("Could not find any record")
		os.Exit(1)
	}
	for _, value := range osimageSlice {
		osimage := value.(map[string]interface{})
		fmt.Printf("%s (osimage)\n", osimage["name"])
	}
}

func ListOsimageCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List osimage(s) in xCAT3 service",
		Long:  `List osimage(s) in xCAT3 service. Format: list`,
		Run:   ListOsimage,
	}
	return cmd
}

func ShowOsimage(cmd *cobra.Command, args []string) {
	var fields []string
	var result interface{}
	if showOsimageOpts.fields != "" {
		fields = strings.Split(showOsimageOpts.fields, ",")
		if exist, _ := utils.Contains(fields, "name"); !exist {
			fields = append(fields, "name")
		}
	}

	client, err := NewOsimageClient()
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

func ShowOsimageCommand() *cobra.Command {
	showOsimageOpts = new(ShowOsimageOptions)
	cmd := &cobra.Command{
		Use:   "show <osimage name>",
		Short: "Show detailed infomation about osimage.",
		Long:  `Show detailed infomation about osimage. Format: show <osimage name>`,
		Run:   ShowOsimage,
	}
	cmd.Flags().StringVarP(&showOsimageOpts.fields, "fields", "i", "",
		`Fields seperated by comma. Only these fields will be fetched from the server.`)
	return cmd
}

func DeleteOsimage(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("Please specify the uuid of osimage to delete")
		os.Exit(1)
	}

	client, err := NewOsimageClient()
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

func DeleteOsimageCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <osimage name>",
		Short: "Unregister osimage from xCAT3 service.",
		Long:  `Unregister osimage from xCAT3 service. Format: delete <osimage name>`,
		Run:   DeleteOsimage,
	}
	return cmd
}

func UpdateOsimage(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		fmt.Println("Please specify the name of osimage and attribute in key=val format to update")
		os.Exit(1)
	}
	patches := arg_array_to_patch(args[1:])
	client, err := NewOsimageClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	result, err := client.Patch(args[0], patches, true)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	utils.PrintJson(result)
}

func UpdateOsimageCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <osimage name> <key=val> [<key=val>]",
		Short: "Update information about registered osimage.",
		Long:  `Update information about registered osimage. Format: update <osimage name> <key=val> [<key=val>]`,
		Run:   UpdateOsimage,
	}
	return cmd
}

func OsimageCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "osimage",
		Short: "This is osimage child command for xcat3",
		Long: `xcat3 osimage --help and xcat3 osimage help COMMAND to see the usage for specfied
	command.`,
	}
	return cmd
}

func init() {
	OsimageCmd := OsimageCommand()
	OsimageCmd.AddCommand(ListOsimageCommand())
	OsimageCmd.AddCommand(ShowOsimageCommand())
	OsimageCmd.AddCommand(DeleteOsimageCommand())
	OsimageCmd.AddCommand(UpdateOsimageCommand())
	RootCmd.AddCommand(OsimageCmd)
}
