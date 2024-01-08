import regex

schema Probe:
    """ Probe describes a health check to be performed against a container to determine whether it is
    alive or ready to receive traffic. There are three probe types: readiness, liveness, and startup.

    Attributes
    ----------
    probeHandler: Exec | Http | Tcp, default is Undefined, required.
        The action taken to determine the alive or health of a container
    initialDelaySeconds: int, default is Undefined, optional.
        The number of seconds before health checking is activated.
        More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
    timeoutSeconds: int, default is Undefined, optional.
        The number of seconds after which the probe times out.
        More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
    periodSeconds: int, default is Undefined, optional.
        How often (in seconds) to perform the probe.
    successThreshold: int, default is Undefined, optional.
        Minimum consecutive successes for the probe to be considered successful after having failed.
    failureThreshold: int, default is Undefined, optional.
        Minimum consecutive failures for the probe to be considered failed after having succeeded.
    terminationGracePeriod: int, default is Undefined, optional.
        Duration in seconds before terminate gracefully upon probe failure.

    Examples
    --------
    import models.schema.v1.workload.container.probe as p

    probe = p.Probe {
        probeHandler: p.Http {
            path: "/healthz"
        }
        initialDelaySeconds: 10
    }
    """

    # The action taken to determine the health of a container
    probeHandler:               Exec | Http | Tcp

    # Number of seconds after the container has started before liveness probes are initiated.
    initialDelaySeconds?:       int

    # Number of seconds after which the probe times out.
    timeoutSeconds?:            int

    # How often (in seconds) to perform the probe.
    periodSeconds?:             int

    # Minimum consecutive successes for the probe to be considered successful after having failed.
    successThreshold?:          int

    # Minimum consecutive failures for the probe to be considered failed after having succeeded.
    failureThreshold?:          int

    # Duration in seconds before terminate gracefully upon probe failure.
    terminationGracePeriod?:    int

    check:
        initialDelaySeconds >= 0 if initialDelaySeconds, "initialDelaySeconds must be greater than or equal to 0"
        timeoutSeconds >= 0 if timeoutSeconds, "timeoutSeconds must be greater than or equal to 0"
        periodSeconds >= 0 if periodSeconds, "periodSeconds must be greater than or equal to 0"
        successThreshold >= 0 if successThreshold, "successThreshold must be greater than or equal to 0"
        failureThreshold >= 0 if failureThreshold, "failureThreshold must be greater than or equal to 0"
        terminationGracePeriod >= 0 if terminationGracePeriod, "terminationGracePeriod must be greater than or equal to 0"

schema Exec:
    """ Exec describes a "run in container" action.

    Attributes
    ----------
    command: str, default is Undefined, required.
        The command line to execute inside the container.

    Examples
    --------
    import models.schema.v1.workload.container.probe as p

    execProbe = p.Exec {
        command: ["probe.sh"]
    }
    """

    # The command line to execute inside the container.
    # Exit status of 0 is treated as live/healthy and non-zero is unhealthy.
    command:                    [str]

    check:
        len(command) > 0, "command must be specified"

schema Http:
    """ Http describes an action based on HTTP Get requests.

    Attributes
    ----------
    url: str, default is Undefined, required.
        The full qualified url to send HTTP requests.
    headers: {str:str}, default is Undefined, optional.
        Collection of custom headers to set in the request

    Examples
    --------
    import models.schema.v1.workload.container.probe as p

    httpProbe = p.Http {
        url: "http://localhost:80"
        headers: {
            "X-HEADER": "VALUE"
        }
    }
    """

    # The full qualified url to send HTTP requests.
    url:                        str

    # Custom headers to set in the request.
    headers?:                   {str:str}

    check:
        all header in headers {
            regex.match(header, r"^[-A-Za-z0-9]+$")
        } if headers, "a valid HTTP header must consist of alphanumeric characters or '-' e.g X-Header-Name"

schema Tcp:
    """ Tcp describes an action based on opening a socket.

    Attributes
    ----------
    url: str, default is Undefined, required.
        The full qualified url to open a socket.

    Examples
    --------
    import models.schema.v1.workload.container.probe as p

    tcpProbe = p.Tcp {
        url: "tcp://localhost:1234"
    }
    """

    # The full qualified url to open a socket.
    url:                        str
