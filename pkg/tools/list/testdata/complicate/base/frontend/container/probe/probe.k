schema Probe:
    handler:                 Exec | Http | Tcp
    initialDelaySeconds?:    int
    timeoutSeconds?:         int
    periodSeconds?:          int = 10
    successThreshold?:       int
    failureThreshold?:       int
