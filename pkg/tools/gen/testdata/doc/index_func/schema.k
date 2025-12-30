# Test schema with both index signature and function types

schema Resource:
    """
    Resource represents a generic resource with indexed metadata.
    """
    name: str
    metadata?: {str: str}

schema ResourceMap:
    """
    ResourceMap is a map of resources indexed by name.
    """
    [name: str]: Resource

schema Config:
    """
    Config contains a resource map and transformation functions.
    """
    resources?: ResourceMap
    transform?: (Resource) -> Resource
    filter?: (str, int) -> bool
    process?: (ResourceMap) -> [str]
    callback?: () -> str
