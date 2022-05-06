# 开发文档

目录结构说明:

- `/`: 根目录对应 kclvm-go 导出的公共函数
- `/cmds/kcl-go`: 基于 kclvm-go 导出的公共函数构建的命令行工具
- `/cmds/kcl-go/command`: kcl-go 命令行对应的包
- `/docs`: 开发文档和图片等其它文件资源
- `/examples`: 独立可执行的例子
- `/pkg`: 内部包(API可能发生变化)
- `/pkg/3rdparty`: 第三方代码或者基于第三方改造等代码
- `/pkg/ast`: KCL 语言语法树
- `/pkg/compiler/parser`: 从 KCL 代码解析 语法树
- `/pkg/kcl`: 导出的公共函数的内部实现
- `/pkg/kclvm_runtime`: 底层 KCLVM 命令包装为进程服务, 为上层的函数提供支持
- `/pkg/langserver`: LSP 服务支持
- `/pkg/logger`: 日志处理
- `/pkg/play`: Playground 实现
- `/pkg/service`: Rest 和 GRPC 等服务的实现
- `/pkg/settings`: `kcl.yaml` 文件解析
- `/pkg/spec`: Protobuf 生成的文件, 请参考 https://github.com/KusionStack/kclvm-go/tree/master/pkg/spec
- `/pkg/tools`: 命令行子命令实现
- `/pkg/utils`: 内部辅助函数
