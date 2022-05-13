# Copyright 2020 The KCL Authors. All rights reserved.

import kcl_plugin.hello as hello

schema TestHello:
    name = hello.tolower('KCL')
    assert name == "kcl", "${name}"

    three = hello.add(1, 2)
    assert three == 3
