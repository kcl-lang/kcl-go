schema Tcp:
    $type: "tcp" = "tcp"
    tcpSocket: int

    check:
        tcpSocket >= 1 and tcpSocket <= 65535, "tcpSocket must be between 1 and 65535, inclusive"
