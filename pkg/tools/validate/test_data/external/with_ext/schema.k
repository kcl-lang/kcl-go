import ext

schema Holder:
    name: str
    payload: ext.Value

    check:
        name != ""
