# Copyright 2020 The KCL Authors. All rights reserved.

import some.pkg as pkgName


schema PersonProtocol:
    name: str
    age: int

mixin PersonMixin for PersonProtocol:
    name = 'kcl'
    age = 1

schema Person[id:str="id-default"]:
    mixin [PersonMixin]
    name: str
    age: int

go = Person("id-001") {
    name: 'golang'
}

if go.name == 'golang':
    print("hello KCL")
