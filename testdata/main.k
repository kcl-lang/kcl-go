import kcl_plugin.hello
import kcl_plugin.cmdb

name = "kcl-go"
age = 2

four = hello.add(1, 3)

schema Person:
    name: str = "kcl"
    age: int = 1

x0 = Person {}
x1 = Person {
    age = 101
}
