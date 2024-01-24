schema Container:
    """
    Container is the common user interface for long-running services.

    Attributes
    ----------
    name: str, required
        The name of the long-running container.
    """
    name: str
    image: str
