schema User:
    name: str
    age: int
    message?: str
    data: Data
    labels: {str:}
    hc: [int]
        
    check:
        age > 10

schema Data:
    id: int
    value: str
