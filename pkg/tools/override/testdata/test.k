import pkg

schema Config:
    image: str
    replicas: int

if True:
    configOther = Config {image = "image/other:v1"}


config = Config {
    image = "image/image:v1"
    replicas: 1
}

