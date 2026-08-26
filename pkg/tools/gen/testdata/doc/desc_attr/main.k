schema Server:
    """
    Server is the common user interface for long-running services.

    Attributes
    ----------
    containers : [Container], required
        The list of containers backing the service.
        Each container runs in its own isolated environment.
    workloadType : str, default is "Deployment", required
        The kind of long-running workload to schedule.
    """
    workloadType : str = "Deployment"
    containers: [Container]
