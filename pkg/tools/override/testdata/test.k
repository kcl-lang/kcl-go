import pkg

schema Config:
    image: str

if True:
    configOther = Config {image = "image/other:v1"}


config = Config {image = "image/image:v1"}

