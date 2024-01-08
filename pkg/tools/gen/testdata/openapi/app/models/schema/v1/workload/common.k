import models.schema.v1.workload.container as c
import models.schema.v1.workload.secret as sec

schema WorkloadBase:
    """ WorkloadBase defines set of attributes shared by different workload profile, e.g Service
    and Job. You can inherit this Schema to reuse these common attributes.

    Attributes
    ----------
    containers: {str:c.Container}, default is Undefined, required.
        Containers defines the templates of containers to be ran.
        More info: https://kubernetes.io/docs/concepts/containers
    secrets: {str:sec.BasicAuthSecret|sec.TokenSecret|sec.OpaqueSecret|sec.TLSSecret|sec.ExternalSecret},
        default is Undefined, optional.
        Secrets can be used to store small amount of sensitive data e.g. password, token
    replicas: int, default is 2, required.
        Number of container replicas based on this configuration that should be ran.
    labels: {str:str}, default is Undefined, optional.
        Labels are key/value pairs that are attached to the workload.
    annotations: {str:str}, default is Undefined, optional.
        Annotations are key/value pairs that attach arbitrary non-identifying metadata to the workload.
    """

    # The templates of containers to be ran.
    containers:                 {str:c.Container}
    # Secrets store small amount of sensitive data such as a password, a token, or a key.
    secrets?:                   {str:sec.Secret}

    # The number of containers that should be ran.
    # Default is 2 to meet high availability requirements.
    replicas:                   int = 2

    ###### Other metadata info
    # Labels and annotations can be used to attach arbitrary metadata as key-value pairs to resources.
    labels?:                    {str:str}
    annotations?:               {str:str}
