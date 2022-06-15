schema Exec:
    $type: "exec" = "exec"
    command: [str]
    check:
        len(command) > 0, "command must be specified"
