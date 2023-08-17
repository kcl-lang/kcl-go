import k8s.deployment
import units

schema Server:
    """
    Server is the common user interface for long-running
    services adopting the best practice of Kubernetes.

    Attributes
    ----------
    workloadType : str, default is "Deployment", required
        Use this attribute to specify which kind of long-running service you want.
        Valid values: Deployment, CafeDeployment.
        See also: kusion_models/core/v1/workload_metadata.k.
    name : str, required
        A Server-level attribute.
        The name of the long-running service.
        See also: kusion_models/core/v1/metadata.k.
    labels : {str:str}, optional
        A Server-level attribute.
        The labels of the long-running service.
        See also: kusion_models/core/v1/metadata.k.

    Examples
    ----------------------
    myCustomApp = AppConfiguration {
        name = "componentName"
    }

    """
    workloadType : str = "Deployment"
    name: str
    labels?: {str: str}
    containers: [Container]
    age: int
    height: float
    port: int | str
    mainContainer: Container
    antiSelf: bool
    others: any
    litStr: "abc"
    litBool: True
    litInt: 123
    litFloat: 1.11
    union: "abc" | 123 | True | 1.11 | Container | units.NumberMultiplier | 1M
    union2: "abc" | "def" = "abc"
    dictAny: {str:}
    listAny: []
    backendWorkload: deployment.Deployment
    numMultiplier: units.NumberMultiplier = 1M
