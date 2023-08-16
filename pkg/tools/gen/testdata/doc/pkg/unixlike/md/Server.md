## Schema Server

Server is the common user interface for long-running services adopting the best practice of Kubernetes.

### Attributes

**age** *required*

`int`

**antiSelf** *required*

`bool`

**backendWorkload** *required*

`Deployment`

**containers** *required*

`[Container]`

**dictAny** *required*

`{str:}`

**height** *required*

`float`

**labels**

`{str:str}`

A Server-level attribute.
The labels of the long-running service.
See also: kusion_models/core/v1/metadata.k.

**listAny** *required*

`[]`

**litBool** *required* *readOnly*

`True`

**litFloat** *required* *readOnly*

`1.11`

**litInt** *required* *readOnly*

`123`

**litStr** *required* *readOnly*

`"abc"`

**mainContainer** *required*

`Container`

**name** *required*

`str`

A Server-level attribute.
The name of the long-running service.
See also: kusion_models/core/v1/metadata.k.

**others** *required*

`any`

**port** *required*

`int | str`

**union** *required*

`"abc" | 123 | True | 1.11 | Container`

**union2** *required*

`"abc" | "def"`

**workloadType** *required*

`str`

Use this attribute to specify which kind of long-running service you want.
Valid values: Deployment, CafeDeployment.
See also: kusion_models/core/v1/workload_metadata.k.


## Source Files

- [Server](server.k)
