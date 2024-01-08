import regex

# mapping between secret type and valid data key
SECRET_TYPE_DATA_MAPPING: {str:[str]} = {
    "basic": ["username", "password"]
    "token": ["token"]
    "certificate": ["tls.crt", "tls.key"]
}

schema Secret:
    """ Secrets are used to provide data that is considered sensitive like passwords, API keys,
    TLS certificates, tokens or other credentials.

    Attributes
    ----------
    type: str, default is Undefined, required.
        Type of secret, used to facilitate programmatic handling of secret data.
    params: {str,str}, default is Undefined, optional.
        Collection of parameters used to facilitate programmatic handling of secret data.
    data: {str,str}, default is Undefined, optional.
        Data contains the non-binary secret data in string form.
    immutable: bool, default is Undefined, optional.
        Immutable, if set to true, ensures that data stored in the Secret cannot be updated.

    Examples
    --------
    import models.schema.v1.workload.secret as sec

    basicAuth = sec.Secret {
        type: "basic"
        data: {
            "username": ""
            "password": ""
        }
    }
    """

    # Types of secrets available to use.
    type:                       "basic" | "token" | "opaque" | "certificate" | "external"

    # Params defines extra parameters used to customize secret handling.
    params?:                    {str:str}

    # Data defines the keys and data that will be used by secret.
    data?:                      {str:str}

    # If immutable set to true, ensures that data stored in the Secret cannot be updated.
    immutable?:                 bool

    check:
        all k in data {
            regex.match(k, r"[A-Za-z0-9_.-]*")
        } if data, "a valid secret data key must consist of alphanumeric characters, '-', '_' or '.'"
        all k in data {
            k in SECRET_TYPE_DATA_MAPPING[type] if len(SECRET_TYPE_DATA_MAPPING[type]) > 0
        } if data, "a valid secret data key name must be one of ${SECRET_TYPE_DATA_MAPPING[type]} for ${type} type secret"

