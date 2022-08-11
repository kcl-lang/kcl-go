# KCL Go Plugin 架构

目标：支持 Go 静态插件、支持 Go 动态插件、支持 Py 动态插件。

## 背景

前期的 kclvm-go 工作原理是底层通过多进程方式启动 KCLVM 命令，然后通过类似 RPC 的形式调用 KCLVM 进程的服务。
此方案不仅仅带来了因多语言混合编程的复杂性，也带来了性能和并发的瓶颈，同时无法通过 Go 直接开发 KCL 的 Plugin。
Python 开发的 KCL Plugin 有以下几个问题：首先是需要依赖外部的 Python 环境，增加了复杂性；其次插件的并发能力
受到 Python 的全局锁限制无法扩展；最后是动态语言开发灵活但是也带来更多的不确定性。

鉴于 Kusion 总体发展是希望去掉对 Python 的依赖，底层完全采用 Rust 实现。因此 kclvm-go 新的方法将采用 CGO 方式调用 KCLVM
导出的 C 函数，因此可以通过 Go 或 Rust 来开发 KCL Plugin。而作为 KCLVM 的 Go SDK 提供 Go 开发 KCL Plugin 也
就成了自然的需求。

## KCLVM 透出的 C 接口

通过 `/kclvm_capi` 目录下的 Rust 工程引入了 KCLVM 实现，然后构建出的 动态库包含 CGO 需要的 C 函数接口：

```c
typedef struct kclvm_service kclvm_service;

kclvm_service * kclvm_service_new(uint64_t plugin_agent);
void kclvm_service_delete(kclvm_service *);

const char* kclvm_service_call(kclvm_service* c,const char * method,const char * args);
const char* kclvm_service_get_error_buffer(kclvm_service* c);
void kclvm_service_clear_error_buffer(kclvm_service* c);

void kclvm_service_free_string(const char * res);
```

其中 `kclvm_service_new` 的 `plugin_agent` 表示 KCL Plugin 函数的代理者，对应以下的函数类型：

```c
extern char* plugin_agent(
	char* method,
	char* args_json,
	char* kwargs_json
);
```

此函数对应的 C 函数地址可以通过 `kusionstack.io/kclvm-go/pkg/kcl_plugin"` 包中的 `kcl_plugin.GetInvokeJsonProxyPtr()` 函数获取，
然后当插件方法被调用时将被映射到对应的 Go 函数。

比如，KCL 中调用 `kcl_plugin.hello.add(1, 2)` 将对应以下的 C 函数调用：

```c
const char* s = plugin_agent("kcl_plugin.hello.add", "1", "2");
// "3"
```

传入参数和返回的结果目前均采用 JSON 格式编码。

## 错误处理

目前返回的结果如果包含 `__kcl_PanicInfo__` 属性，则表示返回了错误：

```go
type PanicInfo struct {
	X__kcl_PanicInfo__ bool `json:"__kcl_PanicInfo__"`

	RustFile string `json:"rust_file,omitempty"`
	RustLine int    `json:"rust_line,omitempty"`
	RustCol  int    `json:"rust_col,omitempty"`

	KclPkgPath string `json:"kcl_pkgpath,omitempty"`
	KclFile    string `json:"kcl_file,omitempty"`
	KclLine    int    `json:"kcl_line,omitempty"`
	KclCol     int    `json:"kcl_col,omitempty"`
	KclArgMsg  string `json:"kcl_arg_msg,omitempty"`

	// only for schema check
	KclConfigMetaFile   string `json:"kcl_config_meta_file,omitempty"`
	KclConfigMetaLine   int    `json:"kcl_config_meta_line,omitempty"`
	KclConfigMetaCol    int    `json:"kcl_config_meta_col,omitempty"`
	KclConfigMetaArgMsg string `json:"kcl_config_meta_arg_msg,omitempty"`

	Message     string `json:"message"`
	ErrTypeCode string `json:"err_type_code,omitempty"`
	IsWarning   string `json:"is_warning,omitempty"`
}
```

比如下面是一种可能的错误形式：

```json
{"__kcl_PanicInfo__":true,"message":"[\"KusionStack\",\"KCL\",123]"}
```

在框架层面会将正常的返回值和错误的返回值分离开。

## 对 Go 语言的支持

CGO 通过对 `plugin_agent` 代理函数的封装，比如通过以下方式实现 `kcl_plugin.hello.add` 插件和方法：

```go
import "kusionstack.io/kclvm-go/pkg/kcl_plugin"

func main() {
	kcl_plugin.RegisterPlugin(kcl_plugin.Plugin{
		Name:      "hello",
		ResetFunc: func() {},
		MethodMap: map[string]PluginMethod{
			"add": {
				Type: &kcl_plugin.MethodType{},
				Body: func(args *kcl_plugin.MethodArgs) (*kcl_plugin.MethodResult, error) {
					v := args.IntArg(0) + args.IntArg(1)
					return &kcl_plugin.MethodResult{V: float64(v)}, nil
				},
			},
		},
	})

	output := kcl_plugin.Invoke("kcl_plugin.hello.add", []interface{}{111, 22}, nil)
	// output: 133
}
```

在定义插件方法时，除了方法的实现函数外，还可以通过 `kcl_plugin.MethodType` 定义方法对应的类型信息（可以用于后续结合 KCL 定义进行类型检查）。

以上的插件可以被编译到最终的二进制中。

## 对 Go 开发的动态插件

目前尚不支持，未来可能包含基于动态库的方案：
Go 开发的插件在 macOS 和 Linux 平台也可以编译成 Go plugin 格式动态库加载，在 Windows 平台可以构建出 dll 动态库加载。


## Python 插件的支持

Python 支持是一个可选的外部依赖。因此 CGO 将通过 dlopen 的方式加在对应的 Python 动态库，然后接入插件的能力。

目前 Python 插件的支持还在开发中。

