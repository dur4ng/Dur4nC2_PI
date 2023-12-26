package command

const (
	// *** Jobs ***
	JobsStr       = "jobs"
	HttpStr       = "http"
	StagedHttpStr = "staged-http"

	// *** Generate ***
	BeaconStr = "beacon"

	// *** Beacons ***
	BeaconsStr           = "beacons"
	BeaconsUseStr        = "use"
	BeaconsBackgroundStr = "background"

	// *** System Information ***
	WhoamiStr = "whoami"

	// *** External Modules ***
	ExternalExtensionsStr        = "extensions"
	RegisterExternalExtensionStr = "register"
	InstallExternalExtensionStr  = "install"
	ListExternalExtensionStr     = "list"
	CallExternalExtensionStr     = "call"

	// *** Exec ***
	ExecuteAssemblyStr  = "execute-assembly"
	ExecuteShellcodeStr = "execute-shellcode"

	// *** Tasks ***
	TaskStr     = "tasks"
	ShowTaskStr = "show"

	// *** Hosts ***
	HostStr = "hosts"

	// *** Utils ***
	DownloadStr = "download"
	UploadStr   = "upload"
)

// Groups
const (
	GenericHelpGroup       = "Generic:"
	ImplantHelpGroup       = "Implant:"
	DBQueriesGroup         = "DB Queries:"
	ListenerGeneratorGroup = "Listeners & Generators:"
	ExtensionHelpGroup     = "Implant - 3rd Party extensions:"
)
