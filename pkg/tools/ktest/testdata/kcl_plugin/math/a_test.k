# Copyright 2020 The KCL Authors. All rights reserved.

import kcl_plugin.math as math

schema TestMath:
    v = math.sum(10, 32)
    assert v == 42, "got=${v}"


schema TestBool:
    a: bool = True
