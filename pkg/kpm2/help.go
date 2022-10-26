package kpm2

const (
	CliHelp = `Usage:

        kpm  <command> [arguments]

The commands are:
        add     add dependencies pkg
        del    del  dependencies pkg
        download    download dependencies pkg to local cache and link to workspace
        graph       print dependencies pkg requirement graph
        init        initialize new module in current directory
        store       Reads and performs actions on kpm store that is on the current filesystem
        publish     publish dependencies pkg
        search      search  dependencies pkg
`
	CliNotFound  = `unknown command`
	CliStoreHelp = `Usage: kpm store <command>

Reads and performs actions on kpm store that is on the current filesystem.

Commands:
      add     <pkg>...         Adds new packages to the store. Example: kpm store add konfig@1.0.0
      addfile <pkg>...         Adds path to the store. Example: kpm store add /root/code`
	CliStoreAddHelp     = `Usage: kpm store add <pkg>...`
	CliStoreAddFileHelp = `Usage: kpm store addfile <path>...`
	CliAddHelp          = `Usage: kpm  add <pkg>...`
	CliDelHelp          = `Usage: kpm del <pkg>...`
	CliInitHelp         = `Usage: kpm init <pkg>`
	CliSearchHelp       = `Usage: kpm search <pkg>`
	CliPublishHelp      = `Usage: kpm publish <pkg>`

	//CliDownloadHelp=`Usage: kpm store add <pkg>...`
	//CliTidyHelp=""
	//CliGraphHelp=""
)
