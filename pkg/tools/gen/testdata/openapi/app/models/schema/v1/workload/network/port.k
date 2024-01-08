schema Port:
    """ Port defines the exposed port of Service, which can be used to describe how the Service
    get accessed.

    Attributes
    ----------
    port: int, default is 80, required.
        The exposed port of the Service.
    targetPort: int, default is Undefined, optional.
        The backend container port. If empty, set it the same as the port.
    protocol: "TCP" | "UDP", default is "TCP", required.
        The protocol to access the port.
    public: bool, default is False, required.
        Public defines whether the port can be accessed through Internet.

    Examples
    --------
    import models.schema.v1.workload.network as n

    port = n.Port {
        port: 80
        targetPort: 8080
        protocol: "TCP"
        public: True
    }
    """

    # The exposed port of the Service.
    port:                       int = 80

    # The backend container port.
    targetPort?:                int

    # The protocol of port.
    protocol:                   "TCP" | "UDP" = "TCP"

    # Public defines whether to expose the port through Internet.
    public:                     bool = False

    check:
        1 <= port <= 65535, "port must be between 1 and 65535, inclusive"
        1 <= targetPort <= 65535 if targetPort, "targetPort must be between 1 and 65535, inclusive"