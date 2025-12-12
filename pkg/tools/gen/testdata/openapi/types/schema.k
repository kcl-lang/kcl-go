schema Test:
	name: str
	surname?: str

schema TestMap:
	[name: str]: Test = {name = name}

schema AppConfig:
    maps?: TestMap
    my_func?: (int) -> int
    replicas: int

app: AppConfig {
    replicas: 2
}
