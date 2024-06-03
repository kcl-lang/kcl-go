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
    gender: "unknown" | "male" | "female"
    gender_opt?: "unknown" | "male" | "female"

schema Company:
    name: str
    employees: [employee]
    persons?: Person

type gender = "unknown" | "male" | "female"
schema PersonPayloadOneOf0:
    int_payload: int

schema PersonPayloadOneOf1:
    string_payload: str

schema PersonPayloadOneOf2:
    bool_payload: bool

