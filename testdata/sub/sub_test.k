# Copyright 2021 The KCL Authors. All rights reserved.

schema TestFoo:
    assert version == "dev"

schema TestSub_failed:
    assert version == "dev-??"

schema TestSub_unit_failed:
    assert 1Ki == 1000, "${1Ki}"

