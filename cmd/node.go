package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/chenglch/golang-xcat3client/utils"
	"github.com/spf13/cobra"
)

type CreateNodeOptions struct {
	nics    []string
	control string
}

type ShowNodeOptions struct {
	fields string
}

type ExportNodeOptions struct {
	filepath string
}

type DeployNodeOptions struct {
	state   string
	osimage string
	delete  bool
}

var (
	createOpts      *CreateNodeOptions
	SUCCESS_RESULTS = map[string]bool{"ok": true,
		"updated":   true,
		"deleted":   true,
		"on":        true,
		"off":       true,
		"net":       true,
		"cdrom":     true,
		"disk":      true,
		"provision": true}
	FIELD_MAP = map[string]string{"control": "control_info",
		"nics": "nics_info"}
	showOpts   *ShowNodeOptions
	exportOpts *ExportNodeOptions
	deployOpts *DeployNodeOptions

	exportFields     = []string{"name", "mgt", "netboot", "type", "arch", "nics_info", "control_info"}
	allowBootDev     = []string{"disk", "net", "cdrom", "status"}
	allowPowerStatus = []string{"on", "off", "boot", "status"}
)

func arg_array_to_patch(args []string) []map[string]string {
	patches := make([]map[string]string, 0)
	for _, arg := range args {
		if !strings.HasPrefix(arg, "/") {
			arg = "/" + arg
		}
		items := strings.Split(arg, "=")
		path := items[0]
		value := items[1]
		patch := make(map[string]string)
		if value != "" {
			patch["op"] = "add"
			patch["path"] = path
			patch["value"] = value
		} else {
			patch["op"] = "remove"
			patch["path"] = path
		}
		patches = append(patches, patch)
	}
	return patches
}

func _print_node_result(ret map[string]interface{}) {
	var success uint64
	var failed uint64
	for k, v := range utils.InterfaceToMap(ret["nodes"]) {
		if _, ok := SUCCESS_RESULTS[v.(string)]; ok {
			success += 1
		} else {
			failed += 1
		}
		fmt.Printf("%s: %s\n", k, v)
	}
	fmt.Printf("\nSuccess: %d Failed: %d\n", success, failed)
}

func _post(data map[string]interface{}, result map[string]interface{}, wg *sync.WaitGroup) {
	client, err := NewNodeClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ret, err := client.Post("", data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	utils.MergeMap(result["nodes"].(map[string]interface{}), ret["nodes"].(map[string]interface{}))
	defer wg.Done()
}

func _patch(data map[string]interface{}, result map[string]interface{}, wg *sync.WaitGroup) {
	client, err := NewNodeClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ret, err := client.Patch("", data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	utils.MergeMap(result["nodes"].(map[string]interface{}), ret["nodes"].(map[string]interface{}))
	defer wg.Done()
}

func _parallel_create(data map[string]interface{}) map[string]interface{} {
	// split into 4 parts
	part_num := 4
	i := 0
	nodes := make([]interface{}, part_num)
	size := len(data["nodes"].([]interface{}))
	per_size := size / part_num
	elements := utils.InterfaceToSlice(data["nodes"])
	for i = 0; i < part_num-1; i++ {
		nodes[i] = elements[i*per_size : (i+1)*per_size]
	}
	nodes[i] = elements[i*per_size:]

	var wg sync.WaitGroup
	runtime.GOMAXPROCS(part_num)
	results := make([]map[string]interface{}, part_num)
	for i := 0; i < part_num; i++ {
		wg.Add(1)
		dataMap := make(map[string]interface{})
		results[i] = make(map[string]interface{})
		results[i]["nodes"] = make(map[string]interface{})
		dataMap["nodes"] = nodes[i]
		go _post(dataMap, results[i], &wg)
	}
	wg.Wait()
	ret := make(map[string]interface{})
	ret["nodes"] = make(map[string]interface{})
	for i := 0; i < part_num; i++ {
		utils.MergeMap(ret["nodes"].(map[string]interface{}), results[i]["nodes"].(map[string]interface{}))
	}
	return ret
}

func _parallel_update(data map[string]interface{}) map[string]interface{} {
	// split into 4 parts
	part_num := 4
	i := 0
	nodes := make([]interface{}, part_num)
	size := len(data["nodes"].([]interface{}))
	per_size := size / part_num
	elements := utils.InterfaceToSlice(data["nodes"])
	for i = 0; i < part_num-1; i++ {
		nodes[i] = elements[i*per_size : (i+1)*per_size]
	}
	nodes[i] = elements[i*per_size:]
	patches := data["patches"]

	var wg sync.WaitGroup
	runtime.GOMAXPROCS(part_num)
	results := make([]map[string]interface{}, part_num)
	for i := 0; i < part_num; i++ {
		wg.Add(1)
		dataMap := make(map[string]interface{})
		results[i] = make(map[string]interface{})
		results[i]["nodes"] = make(map[string]interface{})
		dataMap["nodes"] = nodes[i]
		dataMap["patches"] = patches
		go _patch(dataMap, results[i], &wg)
	}
	wg.Wait()
	ret := make(map[string]interface{})
	ret["nodes"] = make(map[string]interface{})
	for i := 0; i < part_num; i++ {
		utils.MergeMap(ret["nodes"].(map[string]interface{}), results[i]["nodes"].(map[string]interface{}))
	}
	return ret
}

func CreateNodes(cmd *cobra.Command, args []string) {
	var nics []interface{}
	var control map[string]interface{}
	var err error
	if len(createOpts.nics) != 0 {
		nics, err = utils.KeyValueArrayToMapArray(createOpts.nics)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	if createOpts.control != "" {
		control, err = utils.KeyValueToMap(createOpts.control, ",")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	if len(args) == 0 {
		fmt.Println("Could not find node argument")
		os.Exit(1)
	}
	names, err := utils.ToNodeArray(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	attr_map, err := utils.KeyValueArrayToMap(args[1:], "=")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var nodes []interface{}
	for _, name := range names {
		node := make(map[string]interface{})
		node["name"] = name
		if control != nil {
			node["control_info"] = control
		}
		if nics != nil {
			nics_dict := map[string]interface{}{"nics": nics}
			node["nics_info"] = nics_dict
		}
		utils.MergeMap(node, attr_map)
		nodes = append(nodes, node)
	}
	data := make(map[string]interface{})
	data["nodes"] = nodes
	var result map[string]interface{}
	if len(data["nodes"].([]interface{})) < 3000 {
		client, err := NewNodeClient()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		result, err = client.Post("", data)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		result = _parallel_create(data)
	}
	_print_node_result(result)
}

func CreateCommand() *cobra.Command {
	createOpts = new(CreateNodeOptions)
	cmd := &cobra.Command{
		Use:   "create <node range> --nic <key=val,key=val> --nic <key=val,key=val> --control <key=val,key=val>  [<key=val> <key=val>]",
		Short: "Enroll node(s) into xCAT3 service",
		Long: `Enroll node(s) into xCAT3 service.
		Format create <node range> --nic <key=val,key=val> --nic <key=val,key=val> --control <key=val,key=val>  [<key=val> <key=val>]`,
		Run: CreateNodes,
	}
	cmd.Flags().StringArrayVarP(&createOpts.nics, "nic", "i", []string{},
		`Key/value pairs split by comma to indicate network information, like:
		-i mac=42:87:0a:05:00:00,primary=True,name=eth0
		-i mac=42:87:0a:05:00:00,name=eth1`)
	cmd.Flags().StringVarP(&createOpts.control, "control", "c", "",
		`Key/value pairs split by comma used by the control plugin, such as
		bmc_address=11.0.0.0,bmc_password=password,bmc_username=admin`)
	return cmd
}

func ListNodes(cmd *cobra.Command, args []string) {
	client, err := NewNodeClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ret, err := client.Get("", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	nodeSlice := utils.InterfaceToSlice(ret.(map[string]interface{})["nodes"])
	if len(nodeSlice) == 0 {
		fmt.Println("Could not find any record")
		os.Exit(1)
	}

	if len(args) == 0 {
		for _, value := range nodeSlice {
			fmt.Printf("%s (node)\n", value)
		}
	} else if len(args) == 1 {
		names, err := utils.ToNodeArray(args[0])
		if err != nil {
			fmt.Println(err)
		}
		for _, value := range nodeSlice {
			// TODO(chenglch): Use map instread of slice to test if exist in the qeury list
			if ok, _ := utils.Contains(value, names); ok {
				fmt.Printf("%s (node)\n", value)
			}
		}
	}
}

func ListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [<node range>]",
		Short: "List node(s) in xCAT3 service",
		Long:  `List node(s) in xCAT3 service. Format list [<node range>]`,
		Run:   ListNodes,
	}
	return cmd
}

func ShowNodes(cmd *cobra.Command, args []string) {
	var fields []string
	if showOpts.fields != "" {
		fields = strings.Split(showOpts.fields, ",")
	}
	if len(args) == 1 {

	}
	client, err := NewNodeClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(args) != 1 {
		fmt.Println("show command should accept node(s) as the argument.")
		os.Exit(1)
	}
	names, err := utils.ToNodeArray(args[0])
	if err != nil {
		fmt.Println(err)
	}
	result, err := client.Show(names, fields)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	utils.PrintJson(result)
}

func ShowCommand() *cobra.Command {
	showOpts = new(ShowNodeOptions)
	cmd := &cobra.Command{
		Use:   "show <node range>",
		Short: "Show detailed information about node(s).",
		Long: `Show detailed information about node(s).
		Format: show <node range>`,
		Run: ShowNodes,
	}
	cmd.Flags().StringVarP(&showOpts.fields, "fields", "i", "",
		`Fields seperated by comma. Only these fields will be fetched from the server.`)
	return cmd
}

func DeleteNodes(cmd *cobra.Command, args []string) {
	client, err := NewNodeClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(args) != 1 {
		fmt.Println("Delete command should accept node(s) as the argument.")
		os.Exit(1)
	}
	names, err := utils.ToNodeArray(args[0])
	if err != nil {
		fmt.Println(err)
	}
	result, err := client.Delete(names)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_print_node_result(result)
}

func DeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <node range>",
		Short: "Unregister node(s) from the xCAT3 service.",
		Long: `SUnregister node(s) from the xCAT3 service.
		Format: delete <node range>`,
		Run: DeleteNodes,
	}
	return cmd
}

func ImportNodes(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("Import command should accept a json file as the argument.")
		os.Exit(1)
	}
	data, err := utils.ReadJsonFile(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var result map[string]interface{}
	if len(data["nodes"].([]interface{})) < 3000 {
		client, err := NewNodeClient()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		result, err = client.Post("", data)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	} else {
		result = _parallel_create(data)

	}
	_print_node_result(result)
}

func ImportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import <json file>",
		Short: "Import node(s) information from json data file.",
		Long: `Import node(s) information from json data file.
		Format: import <json file>`,
		Run: ImportNodes,
	}
	return cmd
}

func ExportNodes(cmd *cobra.Command, args []string) {
	if exportOpts.filepath == "" {
		fmt.Println("Please specified the output filepath")
		os.Exit(1)
	}
	client, err := NewNodeClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	names, err := utils.ToNodeArray(args[0])
	if err != nil {
		fmt.Println(err)
	}
	result, err := client.Show(names, exportFields)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	utils.WriteJsonFile(exportOpts.filepath, result.([]byte))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func ExportCommand() *cobra.Command {
	exportOpts = new(ExportNodeOptions)
	cmd := &cobra.Command{
		Use:   "export <node range> -o <filepath>",
		Short: "Export node(s) information as a specific json data file.",
		Long: `Export node(s) information as a specific json data file
		Format export <node range> -o <filepath>`,
		Run: ExportNodes,
	}
	cmd.Flags().StringVarP(&exportOpts.filepath, "output", "o", "",
		`The output file stores nodes data in json format.`)
	return cmd
}

func UpdateNodes(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		fmt.Println("show command should accept node(s) and attributes format like key=value as the arguments.")
		os.Exit(1)
	}
	names, err := utils.ToNodeArray(args[0])
	if err != nil {
		fmt.Println(err)
	}
	patches := arg_array_to_patch(args[1:])
	for _, p := range patches {
		path := p["path"]
		key := strings.Split(path, "/")[1]
		if _, ok := FIELD_MAP[key]; ok {
			p["path"] = strings.Replace(p["path"], key, FIELD_MAP[key], 1)
		}
	}
	data := make(map[string]interface{})
	data["nodes"] = make([]interface{}, 0)
	for _, name := range names {
		var node = map[string]string{"name": name}
		data["nodes"] = append(data["nodes"].([]interface{}), node)
	}
	data["patches"] = patches
	var result map[string]interface{}
	if len(data["nodes"].([]interface{})) < 3000 {
		client, err := NewNodeClient()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		result, err = client.Patch("", data)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		result = _parallel_update(data)
	}
	_print_node_result(result)
}

func UpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <node range> <key=val> [<key=val>]",
		Short: "Update information about registered node(s).",
		Long: `Update information about registered node(s).
		update <node range> <key=val> [<key=val>]`,
		Run: UpdateNodes,
	}
	return cmd
}

func BootDev(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		fmt.Println("bootdev command should accept node(s) and status/disk/net/cdrom as the arguments.")
		os.Exit(1)
	}
	names, err := utils.ToNodeArray(args[0])
	if err != nil {
		fmt.Println(err)
	}
	client, err := NewNodeClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var result interface{}
	data := client.ToNodesMap(names)
	if args[1] == "status" {
		result, err = client.Get("boot_device", data)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		if exist, _ := utils.Contains(allowBootDev, args[1]); !exist {
			if err != nil {
				fmt.Printf("Only allow %s\n", strings.Join(allowBootDev, " "))
				os.Exit(1)
			}
		}
		result, err = client.Put("boot_device", args[1], data)
	}
	_print_node_result(result.(map[string]interface{}))

}

func BootDevCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bootdev <node range> net/disk/cdrom/status",
		Short: "Set/Get next boot device (net or disk or cdrom).",
		Long: `Set/Get next boot device (net or disk or cdrom).
		Format: bootdev <node range> net/disk/cdrom/status`,
		Run: BootDev,
	}
	return cmd
}

func PowerNodes(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		fmt.Println("power command should accept node(s) and status/on/off/boot as the arguments.")
		os.Exit(1)
	}
	names, err := utils.ToNodeArray(args[0])
	if err != nil {
		fmt.Println(err)
	}
	client, err := NewNodeClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var result interface{}
	data := client.ToNodesMap(names)
	if args[1] == "status" {
		result, err = client.Get("power", data)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		if exist, _ := utils.Contains(allowPowerStatus, args[1]); !exist {
			if err != nil {
				fmt.Printf("Only allow %s\n", strings.Join(allowPowerStatus, " "))
				os.Exit(1)
			}
		}
		result, err = client.Put("power", args[1], data)
	}
	_print_node_result(result.(map[string]interface{}))
}

func PowerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "power <node range> status/on/off/boot",
		Short: "Power operation on/off/reset/status for nodes.",
		Long:  `Power operation on/off/reset/status for nodes. Format: power <node range> status/on/off/boot`,
		Run:   PowerNodes,
	}
	return cmd
}

func DeployNodes(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Please specified nodes")
		os.Exit(1)
	}
	names, err := utils.ToNodeArray(args[0])
	if err != nil {
		fmt.Println(err)
	}
	client, err := NewNodeClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var result interface{}
	data := client.ToNodesMap(names)
	result, err = client.Deploy(deployOpts.osimage, deployOpts.state, deployOpts.delete, data)
	_print_node_result(result.(map[string]interface{}))
}

func DeployCommand() *cobra.Command {
	deployOpts = new(DeployNodeOptions)
	cmd := &cobra.Command{
		Use:   "deploy <node range> --osimage <osimage> [--state dhcp/nodeset] [-d]",
		Short: "Deploy node(s) into specified state.",
		Long:  `Deploy node(s) into specified state. Format: deploy <node range> --osimage <osimage> [--state dhcp/nodeset] [-d]`,
		Run:   DeployNodes,
	}
	cmd.Flags().StringVarP(&deployOpts.state, "state", "", "nodeset",
		`nodeset' or 'dhcp.`)
	cmd.Flags().StringVarP(&deployOpts.osimage, "osimage", "", "",
		`osimage name`)
	cmd.Flags().BoolVarP(&deployOpts.delete, "delete", "d", false,
		`Recover from deploy state`)
	return cmd
}

func init() {
	RootCmd.AddCommand(CreateCommand())
	RootCmd.AddCommand(ListCommand())
	RootCmd.AddCommand(ShowCommand())
	RootCmd.AddCommand(DeleteCommand())
	RootCmd.AddCommand(ImportCommand())
	RootCmd.AddCommand(UpdateCommand())
	RootCmd.AddCommand(ExportCommand())
	RootCmd.AddCommand(BootDevCommand())
	RootCmd.AddCommand(PowerCommand())
	RootCmd.AddCommand(DeployCommand())
}
