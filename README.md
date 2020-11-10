# 补充、修改RxGo

#### Getting started：Hello World

The **Hello World** program:

```go
package main

import (
	"fmt"
	RxGo "github.com/Winszheng/rxgo"
)

func main() {
	RxGo.Just("Hello", "World", "!").Subscribe(func(x string) {
		fmt.Println(x)
	})
}
```

output:

```
Hello
World
!
```

#### 文档

https://godoc.org/github.com/Winszheng/rxgo

#### 设计说明文档

[specification.md

[](specification.md)

#### 实现文件

[filter.go](filter.go)

#### 测试文件

[filter_test.go](filter_test.go)

