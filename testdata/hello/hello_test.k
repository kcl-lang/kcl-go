# Copyright 2021 The KCL Authors. All rights reserved.

import testing

schema TestPerson:
    a = Person{}
    assert a.name == 'kcl'

schema TestPerson_age:
    a = Person{}
    assert a.age == 1

schema TestPerson_ok:
    a = Person{}
    assert a.name == "kcl"
    assert a.age == 1

schema TestOptions:
    testing.arguments("name", "ktest")
    testing.arguments("age", "123")

    testing.arguments("int0", 10)
    testing.arguments("float0", 0.0)
    testing.arguments("bool-true", True)
    testing.arguments("bool-false", False)

    name = option("name")
    assert name == "ktest"

    age = option("age")
    assert age == 123

    assert option("int0") == 10
    assert option("float0") == 0.0
    #assert option("bool-true") == True
    #assert option("bool-false") == False

schema TestOptions_setting:
    testing.setting_file("./xappinfo/settings.yaml")
    testing.arguments("file", "settings.yaml")

    assert option("global-tenant") == "MAIN_SITE"
    assert option("app-name") == "xappinfo"
    assert option("file") == "settings.yaml"
