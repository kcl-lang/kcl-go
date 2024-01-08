schema Database: 
    """ As an important supporting accessory, Database describes the attributes
    to locally deploy or create a cloud provider managed database instance for 
    the workload. 

    Attributes
    ----------
    type: str, default is Undefined, required. 
        Type defines the local deployment mode or the specific cloud vendor that 
        provides the relational database service (rds). 
    engine: str, default is Undefined, required. 
        Engine defines the database engine to use. 
    version: str, default is Undefined, required. 
        Version defines the database engine version to use. 
    instanceType: str, default is Undefined, optional. 
        InstanceType defines the type of the database which is required when 
        creating an rds instance provided by the cloud vendor. 
    size: int, default is 10, optional. 
        Size defines the allocated storage size of the rds instance provided by 
        the cloud vendor in GB.  
    category: str, default is "Basic", optional. 
        Category defines the edition of the rds instance provided by the cloud 
        vendor. 
    username: str, default is "root", optional. 
        Username defines the operation account for the database. 
    securityIPs: [str], default is ["0.0.0.0/0"], optional. 
        SecurityIPs defines the list of IP addresses allowed to access the rds 
        instance provided by the cloud vendor. 
    subnetID: str, default is Undefined, optional. 
        SubnetID defines the virtual subnet ID associated with the VPC that the rds 
        instance will be created in. The rds instance won't be created in user's own VPC 
        if this field is not provided. 
    privateRouting: bool, default is True, optional. 
        PrivateRouting defines whether the host address of the rds instance for the workload 
        to connect with is via public network or private network of the cloud vendor. 
    extraMap: {str:str}, default is Undefined, optional. 
        ExtraMap defines the diversified rds configuration items from different
        cloud vendors. 

    Examples
    --------
    Instantiate an aws rds with mysql 5.7. 

    import models.schema.v1.accessories.database as db

    database: db.Database {
        type: "aws"
        engine: "mysql"
        version: "5.7"
        instanceType: "db.t3.micro"
    }
    """

    # The local deployment mode or the specific cloud vendor that provides the 
    # relational database service (rds). 
    type:                       str

    # The database engine to use. 
    engine:                     str

    # The database engine version to use. 
    version:                    str

    # The type of the database which is required when creating an rds instance
    # provided by the cloud vendor. 
    instanceType?:              str

    # The allocated storage size of the rds instance provided by the cloud vendor 
    # in GB. 
    size?:                      int = 10

    # The edition of the rds instance provided by the cloud vendor. 
    category?:                  str = "Basic"

    # The operation account for the database. 
    username?:                  str = "root"

    # The list of IP addresses allowed to access the rds instance provided by the 
    # cloud vendor. 
    securityIPs?:               [str] = ["0.0.0.0/0"]

    # The virtual subnet ID associated with the VPC that the rds instance will be
    # created in. 
    subnetID?:                  str

    # Whether the host address of the rds instance for the workload to connect with 
    # is via public network or private network of the cloud vendor. 
    privateRouting?:            bool = True

    # The diversified rds configuration items from different cloud vendors. 
    extraMap?:                  {str:str}

    check: 
        instanceType if type != "local", "instanceType is required for cloud provider managed database"