import models.schema.v1.workload as wl
import models.schema.v1.trait as t
import models.schema.v1.accessories.database as db
import models.schema.v1.monitoring as m

schema AppConfiguration:
    """ AppConfiguration is a developer-centric definition that describes how to run an Application.
    This application model builds upon a decade of experience at AntGroup running super large scale
    internal developer platform, combined with best-of-breed ideas and practices from the community.

    Attributes
    ----------
    workload: wl.Service | wl.Job, default is Undefined, required.
        Workload defines how to run your application code. Currently supported workload profile
        includes Service and Job.
    monitoring: m.Prometheus, default is Undefined, optional.
        Monitoring specifies how to scrape Prometheus metrics for the workload.
    database: db.Database, default is Undefined, optional.
        Database describes a locally deployed or a cloud provider managed database instance for the workload.
    opsRule: t.OpsRule, default is Undefined, optional.
        OpsRule specifies collection of rules that will be checked for Day-2 operation.
    labels: {str:str}, default is Undefined, optional.
        Labels can be used to attach arbitrary metadata as key-value pairs to resources.
    annotations: {str:str}, default is Undefined, optional.
        Annotations are key/value pairs that attach arbitrary non-identifying metadata to resources.

    Examples
    --------
    # Instantiate an App with a long-running service and its image is "nginx:v1"

    import models.schema.v1 as ac
    import models.schema.v1.workload as wl
    import models.schema.v1.workload.container as c

    appConfiguration = ac.AppConfiguration {
        workload: wl.Service {
            containers: {
                "nginx": c.Container {
                    image: "nginx:v1"
                }
            }
        }
    }
    """

    # Workload defines how to run your application code.
    workload:                   wl.Service | wl.Job

    # Monitoring specifies how to scrape metrics for the workload.
    monitoring?:                m.Prometheus
    
    # OpsRule specifies collection of rules that will be checked for Day-2 operation.
    opsRule?:                   t.OpsRule

    # Database describes a locally deployed or a cloud provider managed database instance for the workload.
    database?:                  db.Database

    ###### Other metadata info
    # Labels and annotations can be used to attach arbitrary metadata as key-value pairs to resources.
    @info(hidden=True)
    labels?:                    {str:str}
    @info(hidden=True)
    annotations?:               {str:str}
