# Copyright 2020 The KCL Authors. All rights reserved.

import kcl_plugin.hello as hello

schema TestHello_failed:
    name = hello.tolower('KCL')
    assert name == "kcl122", "got=${name}"

