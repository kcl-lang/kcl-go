import kcl_plugin.hello as hello

schema TestPlugin_hello:
    assert hello.tolower('KCL') == 'kcl'
    assert hello.add(1, 2) == 3

schema TestPlugin_hello_bug:
    assert hello.add(1, 2) == 4
