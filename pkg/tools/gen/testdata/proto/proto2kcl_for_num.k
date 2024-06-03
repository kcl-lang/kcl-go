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
    payload: PersonPayloadOneOf0 | PersonPayloadOneOf1 | PersonPayloadOneOf2

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
    gender: 0 | 1 | 2
    gender_opt?: 0 | 1 | 2

schema Company:
    name: str
    employees: [employee]
    persons?: Person

type gender = 0 | 1 | 2
schema PersonPayloadOneOf0:
    int_payload: int

schema PersonPayloadOneOf1:
    string_payload: str

schema PersonPayloadOneOf2:
    bool_payload: bool

