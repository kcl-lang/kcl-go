schema Person:
    name: str
    age: int
    friends: [str]
    movies: {str:Movie}
    employee?: employee
    company: Company
    star_int?: int
    star_map: {str:str}
    inter: any

schema Movie:
    desc: str
    size: int
    kind: str
    unknown1: any
    unknown2: any

schema employee:
    name: str
    age: int
    friends: [str]
    movies: {str:Movie}
    bank_card: int
    nationality: str
    salary: float
    age_double: float
    is_married: bool

schema Company:
    name: str
    employees: [employee]
    persons?: Person

