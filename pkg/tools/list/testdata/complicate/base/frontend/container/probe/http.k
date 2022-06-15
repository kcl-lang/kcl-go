schema Http:
    $type: "http" = "http"
    path: str
    port: int
    scheme: "HTTP" | "HTTPS" = "HTTP"

    check:
        1 <= port <= 65535, "http port must be between 1 and 65535, inclusive"
