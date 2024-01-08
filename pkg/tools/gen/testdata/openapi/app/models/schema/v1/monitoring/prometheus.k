import regex

schema Prometheus:
    """ Prometheus can be used to define monitoring requirements

    Attributes
    ----------
    interval: str, default to the Prometheus global scraping interval, optional
        The time interval which Prometheus scrapes metrics data. Only applicable when operator mode is set to true.
        When operator mode is set to false, the scraping interval can only be set in the scraping job configuration, which kusion does not have permission to manage directly.
    timeout: str, default to the Prometheus global scraping timeout, optional
        The timeout when Prometheus scrapes metrics data. Only applicable when operator mode is set to true.
        When operator mode is set to false, the scraping timeout can only be set in the scraping job configuration, which kusion does not have permission to manage directly.
    path: str, default to the Prometheus global scraping path, which should be /metrics if not explicitly set in the Prometheus scrape config, optional
        The path to scrape metrics from.
    port: str, default to container ports when scraping pod (monitorType is pod), default to service port when scraping service (monitorType is service), optional
        The port to scrape metrics from. When using Prometheus operator, this needs to be the port NAME. Otherwise, this can be a port name or a number.
    scheme: str, default to http, optional
        The scheme to scrape metrics from. Possible values are http and https.

    Examples
    --------
    import models.schema.v1.monitoring as m

    monitoring: m.Prometheus{
        interval:       "30s"
        timeout:        "15s"
        path:           "/metrics"
        port:           "web"
        scheme:         "http"
    }
    """

    # Interval defines the time interval Prometheus scrapes for metrics. If not provided, this will default to the global scraping configuration for that Prometheus instance.
    interval?:       str

    # Timeout defines the time interval before Prometheus considers the scraping timed out. If not provided, this will default to the global scraping configuration for that Prometheus instance.
    timeout?:        str

    # Path defines the path from which Prometheus scrapes the target.
    path?:           str

    # Path defines the port from which Prometheus scrapes the target.
    port?:           str

    # Scheme defines the protocol Prometheus scrapes for metrics. Possible values are http and https.
    scheme?:         "http" | "https" = "http"

    check:
        regex.match(interval, r"[A-Za-z0-9_.-]*"), "a valid interval must starts with alphanumeric and end with y/w/d/h/m/s/ms"
        regex.match(timeout, r"[A-Za-z0-9_.-]*"), "a valid timeout must starts with alphanumeric and end with y/w/d/h/m/s/ms"
        port and path if operatorMode == False, "port and path must be set only when operator mode is false"