import models.schema.v1.workload.container.probe as p

schema Lifecycle:
    """ Lifecycle describes actions that the management system should take in response
    to container lifecycle events.

    Attributes
    ----------
    preStop: p.Exec | p.Http, default is Undefined, optional.
        The action to be taken before a container is terminated due to an API request or
        management event such as liveness/startup probe failure, preemption, resource contention, etc.
        More info: https://kubernetes.io/docs/concepts/containers/container-lifecycle-hooks/#container-hooks
    postStart: p.Exec | p.Http, default is Undefined, optional.
        The action to be taken after a container is created.
        More info: https://kubernetes.io/docs/concepts/containers/container-lifecycle-hooks/#container-hooks

    Examples
    --------
    import models.schema.v1.workload.container.probe as p
    import models.schema.v1.workload.container.lifecycle as lc

    lifecycleHook = lc.Lifecycle {
        preStop: p.Exec {
            command: ["preStop.sh"]
        }
        postStart: p.Http {
            url: "http://localhost:80"
        }
    }
    """

    # The action to be taken before a container is terminated.
    preStop?:                   p.Exec | p.Http

    # The action to be taken after a container is created.
    postStart?:                 p.Exec | p.Http
