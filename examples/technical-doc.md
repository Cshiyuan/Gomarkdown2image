# 技术文档示例

## Go 语言最佳实践

### 错误处理模式

在 Go 中,错误处理是显式的:

```go
package main

import (
    "fmt"
    "os"
)

func ReadFile(filename string) ([]byte, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to read %s: %w", filename, err)
    }
    return data, nil
}

func main() {
    data, err := ReadFile("config.json")
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Read %d bytes\n", len(data))
}
```

### 并发模式

使用 goroutine 和 channel:

```go
func Pipeline() {
    // 生成器
    gen := func(nums ...int) <-chan int {
        out := make(chan int)
        go func() {
            for _, n := range nums {
                out <- n
            }
            close(out)
        }()
        return out
    }

    // 平方计算
    sq := func(in <-chan int) <-chan int {
        out := make(chan int)
        go func() {
            for n := range in {
                out <- n * n
            }
            close(out)
        }()
        return out
    }

    // 使用管道
    for n := range sq(sq(gen(2, 3))) {
        fmt.Println(n) // 16, 81
    }
}
```

## Python 装饰器模式

```python
import functools
import time

def timer(func):
    """计时装饰器"""
    @functools.wraps(func)
    def wrapper(*args, **kwargs):
        start = time.time()
        result = func(*args, **kwargs)
        end = time.time()
        print(f"{func.__name__} took {end - start:.2f}s")
        return result
    return wrapper

@timer
def slow_function():
    time.sleep(1)
    return "Done"

slow_function()  # 输出: slow_function took 1.00s
```

## JavaScript Promise 链

```javascript
// Promise 链式调用
fetch('https://api.example.com/data')
  .then(response => {
    if (!response.ok) {
      throw new Error('Network response was not ok');
    }
    return response.json();
  })
  .then(data => {
    console.log('Data:', data);
    return processData(data);
  })
  .then(processed => {
    console.log('Processed:', processed);
  })
  .catch(error => {
    console.error('Error:', error);
  });

// async/await 版本
async function fetchData() {
  try {
    const response = await fetch('https://api.example.com/data');
    if (!response.ok) throw new Error('Network error');
    const data = await response.json();
    const processed = await processData(data);
    console.log('Processed:', processed);
  } catch (error) {
    console.error('Error:', error);
  }
}
```

## SQL 查询示例

```sql
-- 创建用户表
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT true
);

-- 复杂查询:查找活跃用户及其订单统计
SELECT
    u.username,
    u.email,
    COUNT(o.id) as order_count,
    SUM(o.total_amount) as total_spent
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
WHERE u.is_active = true
    AND o.created_at >= NOW() - INTERVAL '30 days'
GROUP BY u.id, u.username, u.email
HAVING COUNT(o.id) > 0
ORDER BY total_spent DESC
LIMIT 10;
```

## 性能对比表格

| 语言 | 执行时间 (ms) | 内存占用 (MB) | 代码行数 |
|------|--------------|--------------|---------|
| Go | 150 | 45 | 120 |
| Python | 420 | 78 | 95 |
| JavaScript | 280 | 62 | 110 |
| Rust | 80 | 32 | 150 |

## 引用说明

> **重要提示**:
>
> 在生产环境中,始终要:
> - 验证所有外部输入
> - 使用参数化查询防止 SQL 注入
> - 实施适当的错误处理和日志记录
> - 定期更新依赖以修复安全漏洞

## 任务清单

开发前检查清单:

- [x] 代码审查
- [x] 单元测试(覆盖率 > 80%)
- [x] 文档更新
- [ ] 性能基准测试
- [ ] 安全扫描
- [ ] 部署到预生产环境

## 数学公式

时间复杂度分析:

- 线性查找: O(n)
- 二分查找: O(log n)
- 快速排序: O(n log n) 平均情况
- 哈希表查找: O(1) 平均情况

---

**文档版本**: 1.0.0
**更新日期**: 2025-12-14
