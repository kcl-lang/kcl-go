import models.schema.v1.workload.network as n

schema Service(WorkloadBase):
    """ Service is a kind of workload profile that describes how to run your application code. This
    is typically used for long-running web applications that should "never" go down, and handle
    short-lived latency-sensitive web requests, or events.

    Attributes
    ----------
    ports: [n.Port], default is Undefined, optional.
        The list of ports which the Service should get exposed.

    Examples
    --------
    # Instantiate a long-running service and its image is "nginx:v1"

    import models.schema.v1.workload as wl
    import models.schema.v1.workload.container as c

    svc = wl.Service {
        containers: {
            "nginx": c.Container {
                image: "nginx:v1"
            }
        }
        ports: [
            n.Port {
                port: 80
                public: True
            }
            n.Port {
                port: 9090
            }
        ]
    }
    """

    # The list of ports get exposed.
    ports?:                      [n.Port]
