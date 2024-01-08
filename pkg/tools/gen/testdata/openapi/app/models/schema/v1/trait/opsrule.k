schema OpsRule:
    """ OpsRule describes operation rules for various Day-2 Operations. Once declared, these
    operation rules will be checked before any Day-2 operations.

    Attributes
    ----------
    maxUnavailable: str or int, default is Undefined, optional.
        The maximum percentage of the total pod instances in the component that can be
        simultaneously unhealthy.

    Examples
    --------
    import models.schema.v1.trait as t

    opsRule = t.OpsRule {
        maxUnavailable: "30%"
    }
    """

    # The maximum percentage of the total pod instances in the component that can be
    # simultaneously unhealthy.
    maxUnavailable?:            int | str = "25%"