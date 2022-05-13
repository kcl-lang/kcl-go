# Copyright 2020 The KCL Authors. All rights reserved.

INFO = {
    'name': 'hello',
    'describe': 'hello doc',
    'long_describe': 'long describe',
    'version': '0.0.1',
}


def KMANGLED_add(a: int, b: int) -> int:
    """add two numbers, and return result"""
    return a + b


def KMANGLED_tolower(s: str) -> str:
    return s.lower()
