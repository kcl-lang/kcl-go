# Copyright 2021 The KCL Authors. All rights reserved.

schema Person:
    name: str = "kcl"
    age: int = 1

hello = Person {
    name: "hello kcl"
    age: 102
}
