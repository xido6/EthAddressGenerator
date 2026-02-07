## 以太坊地址生成器

根据参数生成特殊以太坊地址

### 资源

下载即可用

- 对于Windows平台

  资源位于 [/resource/windows_amd64/generator.exe](https://github.com/xido6/EthAddressGenerator/blob/main/resource/windows_amd64/generator.exe)

- 对于mac M系列芯片

  资源位于 [/resource/macOS/arm64/generator](https://github.com/xido6/EthAddressGenerator/blob/main/resource/macOS/arm64/generator)

- 对于mac intel系列芯片

  资源位于 [/resource/macOS/amd64/generator](https://github.com/xido6/EthAddressGenerator/blob/main/resource/macOS/amd64/generator)

### 接受参数

|     flag     |         explain          |
| :----------: | :----------------------: |
|  -lead-char  |         前缀字符         |
| -lead-count  |        前缀字符数        |
| -trail-char  |         后缀字符         |
| -trail-count |        后缀字符数        |
|   -workers   | 并发数(程序并发协程数量) |

**note：前缀字符数参数和后缀字符数参数不可同时为0，并且当字符数参数不为0时，其对应字符参数必须为合法的单个16进制字符，workers为0或过大时，默认使用当前进程最大可见cpu核心数，eg：**

```
./generator.exe -lead-char 8 -lead-count 2 -trail-char 6 -trail-count 2 -workers 20
```

**将以20个协程，生成一个拥有2个前缀"8"，2个后缀"6"的eth地址，例如0x88ff2ad3c8d1e2826a3000f6b4bca17aba6f4d66，程序运行结果将保存至当前目录下的search_result.json文件中**

### 简单测试

在一台13代i7机器上测试结果

(随手测试，取3次运行结果的中位数，非标准benchmark测试)

| 前缀数+后缀数 | 并发数 | 生成时间  |
| :-----------: | :----: | :-------: |
|       1       |   1    |   ~0ms    |
|       2       |   1    |  1.57ms   |
|       3       |   1    |  30.79ms  |
|       3       |   5    |  10.83ms  |
|       3       |   10   |  10.18ms  |
|       3       |   20   |  12.33ms  |
|       4       |   1    | 725.61ms  |
|       4       |   5    | 284.46ms  |
|       4       |   10   | 105.52ms  |
|       4       |   20   |  27.1ms   |
|       5       |   1    | 8085.19ms |
|       5       |   5    | 7565.51ms |
|       5       |   10   | 5310.74ms |
|       5       |   20   | 3387.96ms |
|      ...      |  ...   |    ...    |
|       8       |   20   |   1.5h    |
|      ...      |  ...   |    ...    |

misc：

- 前缀+后缀总数＜4时难度过低，测试结果受运气影响过大，生成此类地址时易发生过度并发问题，建议不要使用过多worker，因为CPU上下文切换耗时可能比真正生成地址耗时还大，反而拖慢整体运行时间；

- 前缀+后缀数每增加1，难度指数级增长，要求8位特定字符时，程序耗时以小时为单位，10位及以上特定字符时，概率为1/16^10，按照每秒生成100w个地址计算，需要耗时约12.7天，个人PC大概率无法完成任务，建议视自己机器性能而定，不要追求过高难度。