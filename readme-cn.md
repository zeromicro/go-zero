<img align="right" width="150px" src="https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/go-zero.png">

# go-zero

***缩短从需求到上线的距离***

[English](readme.md) | 简体中文 | [한국어](readme-ko.md)

[![Go Report Card](https://goreportcard.com/badge/github.com/zeromicro/go-zero)](https://goreportcard.com/report/github.com/zeromicro/go-zero)
[![goproxy](https://goproxy.cn/stats/github.com/zeromicro/go-zero/badges/download-count.svg)](https://goproxy.cn/stats/github.com/zeromicro/go-zero/badges/download-count.svg)
[![codecov](https://codecov.io/gh/zeromicro/go-zero/branch/master/graph/badge.svg)](https://codecov.io/gh/zeromicro/go-zero)
[![Release](https://img.shields.io/github/v/release/zeromicro/go-zero.svg?style=flat-square)](https://github.com/zeromicro/go-zero)
[![Go Reference](https://pkg.go.dev/badge/github.com/zeromicro/go-zero.svg)](https://pkg.go.dev/github.com/zeromicro/go-zero)
[![Awesome Go](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/avelino/awesome-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

<a href="https://trendshift.io/repositories/3263" target="_blank"><img src="https://trendshift.io/api/badge/repositories/3263" alt="zeromicro%2Fgo-zero | Trendshift" style="width: 250px; height: 55px;" width="250" height="55"/></a>
<a href="https://www.producthunt.com/posts/go-zero?utm_source=badge-featured&utm_medium=badge&utm_souce=badge-go&#0045;zero" target="_blank"><img src="https://api.producthunt.com/widgets/embed-image/v1/featured.svg?post_id=334030&theme=light" alt="go&#0045;zero - A&#0032;web&#0032;&#0038;&#0032;rpc&#0032;framework&#0032;written&#0032;in&#0032;Go&#0046; | Product Hunt" style="width: 250px; height: 54px;" width="250" height="54" /></a>

## go-zero 介绍

go-zero（收录于 CNCF 云原生技术全景图：[https://landscape.cncf.io/?selected=go-zero](https://landscape.cncf.io/?selected=go-zero)）是一个集成了各种工程实践的 web 和 rpc 框架。通过弹性设计保障了大并发服务端的稳定性，经受了充分的实战检验。

go-zero 包含极简的 API 定义和生成工具 goctl，可以根据定义的 api 文件一键生成 Go, iOS, Android, Kotlin, Dart, TypeScript, JavaScript 代码，并可直接运行。

使用 go-zero 的好处：

* 经过千万日活服务验证的稳定性
* 内建弹性保护:级联超时、限流、熔断、降载(无需配置)
* 极简 API 语法生成多端代码
* 自动参数校验和丰富的微服务工具包

![架构图](https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/architecture.png)

## go-zero 框架背景

18 年初，我们决定从 `Java+MongoDB` 的单体架构迁移到微服务架构，选择:

* **基于 Go 语言** - 高效性能、简洁语法、极致部署体验、极低资源成本
* **自研微服务框架** - 更快速的问题定位、更便捷的新特性增加

## go-zero 框架设计思考

go-zero 遵循以下核心设计准则:

* **保持简单** - 简单是第一原则
* **高可用** - 高并发、易扩展
* **弹性设计** - 面向故障编程
* **工具驱动** - 工具大于约定和文档
* **业务友好** - 封装复杂度、一事一法

## go-zero 项目实现和特点

go-zero 集成各种工程实践，主要特点:

* **强大工具支持** - 尽可能少的代码编写
* **极简接口** - 完全兼容 net/http
* **高性能** - 优化的速度和效率
* **弹性设计** - 内建限流、熔断、降载，自动触发、自动恢复
* **服务治理** - 内建服务发现、负载均衡、链路跟踪
* **开发工具** - API 参数自动校验、超时级联控制、自动缓存控制

![弹性设计](https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/resilience.jpg)

## 我们使用 go-zero 的基本架构图

<img width="1067" alt="image" src="https://user-images.githubusercontent.com/1918356/171880582-11a86658-41c3-466c-95e7-7b1220eecc52.png">

觉得不错的话，别忘 **star** 👏

## Installation

在项目目录下通过如下命令安装：

```shell
GO111MODULE=on GOPROXY=https://goproxy.cn/,direct go get -u github.com/zeromicro/go-zero
```

## AI 原生开发

go-zero 团队构建了完整的 AI 工具生态，让 Claude、GitHub Copilot、Cursor 生成符合 go-zero 规范的代码。

### 三大核心项目

**[ai-context](https://github.com/zeromicro/ai-context)** - AI 的工作流程指南

**[zero-skills](https://github.com/zeromicro/zero-skills)** - 模式库和示例

**[mcp-zero](https://github.com/zeromicro/mcp-zero)** - 基于 MCP 的代码生成工具

### 快速配置

#### GitHub Copilot
```bash
git submodule add https://github.com/zeromicro/ai-context.git .github/ai-context
ln -s ai-context/00-instructions.md .github/copilot-instructions.md  # macOS/Linux
# Windows: mklink .github\copilot-instructions.md .github\ai-context\00-instructions.md
git submodule update --remote .github/ai-context  # 更新
```

#### Cursor
```bash
git submodule add https://github.com/zeromicro/ai-context.git .cursorrules
git submodule update --remote .cursorrules  # 更新
```

#### Windsurf
```bash
git submodule add https://github.com/zeromicro/ai-context.git .windsurfrules
git submodule update --remote .windsurfrules  # 更新
```

#### Claude Desktop
```bash
git clone https://github.com/zeromicro/mcp-zero.git && cd mcp-zero && go build
# 配置: ~/Library/Application Support/Claude/claude_desktop_config.json
# 或: claude mcp add --transport stdio mcp-zero --env GOCTL_PATH=/path/to/goctl -- /path/to/mcp-zero
```

### 协同工作原理

AI 助手通过三个工具协同配合:
1. **ai-context** - 工作流程指导
2. **zero-skills** - 实现模式
3. **mcp-zero** - 实时代码生成

**示例**:创建新的 REST API → AI 读取 **ai-context** 了解工作流 → 调用 **mcp-zero** 生成代码 → 参考 **zero-skills** 实现模式 → 生成符合规范的代码 ✅

## Quick Start

0. 完整示例请查看

    [快速构建高并发微服务](https://github.com/zeromicro/zero-doc/blob/main/doc/shorturl.md)

    [快速构建高并发微服务 - 多 RPC 版](https://github.com/zeromicro/zero-doc/blob/main/docs/zero/bookstore.md)

1. 安装 goctl 工具


    ```shell
    # Go
    GOPROXY=https://goproxy.cn/,direct go install github.com/zeromicro/go-zero/tools/goctl@latest

    # For Mac
    brew install goctl

    # docker for all platforms
    docker pull kevinwan/goctl
    # run goctl
    docker run --rm -it -v `pwd`:/app kevinwan/goctl --help
    ```

    确保 goctl 可执行并在 $PATH 环境变量里。

2. 快速生成 api 服务

    ```shell
    goctl api new greet
    cd greet
    go mod tidy
    go run greet.go -f etc/greet-api.yaml
    ```

    默认侦听在 `8888` 端口（可以在配置文件里修改），可以通过 `curl` 请求：

    ```shell
    curl -i http://localhost:8888/from/you
    ```

    返回如下：

    ```http
    HTTP/1.1 200 OK
    Content-Type: application/json
    Date: Thu, 22 Oct 2020 14:03:18 GMT
    Content-Length: 14

    {"message":""}
    ```

    编写业务代码：

      * api 文件定义了服务对外 HTTP 接口，可参考 [api 规范](https://github.com/zeromicro/zero-doc/blob/main/docs/zero/goctl-api.md)
      * 可以在 `servicecontext.go` 里面传递依赖给 logic，比如 mysql, redis 等
      * 在 api 定义的 `get/post/put/delete` 等请求对应的 logic 里增加业务处理逻辑

3. 生成多语言客户端代码

    ```shell
    goctl api java -api greet.api -dir greet
    goctl api dart -api greet.api -dir greet
    ...
    ```

## Benchmark

![benchmark](https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/benchmark.png)

[测试代码见这里](https://github.com/smallnest/go-web-framework-benchmark)

## 文档

* API 文档

  [https://go-zero.dev](https://go-zero.dev)

* awesome 系列（更多文章见『微服务实践』公众号）

  * [快速构建高并发微服务](https://github.com/zeromicro/zero-doc/blob/main/doc/shorturl.md)
  * [快速构建高并发微服务 - 多 RPC 版](https://github.com/zeromicro/zero-doc/blob/main/docs/zero/bookstore.md)
  * [goctl 使用帮助](https://github.com/zeromicro/zero-doc/blob/main/doc/goctl.md)
  * [Examples](https://github.com/zeromicro/zero-examples)

* 精选 `goctl` 插件

  | 插件    | 用途  |
  | ------------- |:-------------|
  | [goctl-swagger](https://github.com/zeromicro/goctl-swagger) | 一键生成 `api` 的 `swagger` 文档 |
  | [goctl-android](https://github.com/zeromicro/goctl-android) | 生成 `java (android)` 端 `http client` 请求代码 |
  | [goctl-go-compact](https://github.com/zeromicro/goctl-go-compact) | 合并 `api` 里同一个 `group` 里的 `handler` 到一个 `go` 文件 |

## go-zero 用户

go-zero 已被众多公司用于生产部署，场景涵盖在线教育、电商、游戏、区块链等。目前使用 go-zero 的公司包括但不限于:

>1. 好未来
>2. 上海晓信信息科技有限公司（晓黑板）
>3. 上海玉数科技有限公司
>4. 常州千帆网络科技有限公司
>5. 上班族科技
>6. 英雄体育（VSPN）
>7. githubmemory
>8. 释空(上海)品牌策划有限公司(senkoo)
>9. 鞍山三合众鑫科技有限公司
>10. 广州星梦工场网络科技有限公司
>11. 杭州复杂美科技有限公司
>12. 赛凌科技
>13. 捞月狗
>14. 浙江三合通信科技有限公司
>15. 爱克萨
>16. 郑州众合互联信息技术有限公司
>17. 三七游戏
>18. 成都创道夫科技有限公司
>19. 联想Lenovo
>20. 云犀
>21. 高盈国际
>22. 北京中科生活服务有限公司
>23. Indochat 印尼艾希英
>24. 数赞
>25. 量冠科技
>26. 杭州又拍云科技有限公司
>27. 深圳市点购电子商务控股股份有限公司
>28. 深圳市宁克沃德科技有限公司
>29. 桂林优利特医疗电子有限公司
>30. 成都智橙互动科技有限公司
>31. 深圳市班班科技有限公司
>32. 飞视（苏州）数字技术有限公司
>33. 上海鲸思智能科技有限公司
>34. 南宁宸升计算机科技有限公司
>35. 秦皇岛2084team
>36. 天翼云股份有限公司
>37. 南京速优云信息科技有限公司
>38. 北京小鸦科技有限公司
>39. 深圳无边界技术有限公司
>40. 马鞍山百助网络科技有限公司
>41. 上海阿莫尔科技有限公司
>42. 发明者量化
>43. 济南超级盟网络科技有限公司
>44. 苏州互盟信息存储技术有限公司
>45. 成都艾途教育科技集团有限公司
>46. 上海游族网络
>47. 深信服
>48. 中免日上科技互联有限公司
>49. ECLOUDVALLEY TECHNOLOGY (HK) LIMITED
>50. 馨科智（深圳）科技有限公司
>51. 成都松珀科技有限公司
>52. 亿景智联
>53. 上海扩博智能技术有限公司
>54. 一犀科技成都有限公司
>55. 北京术杰科技有限公司
>56. 时代脉搏网络科技（云浮市）有限公司
>57. 店有帮
>58. 七牛云
>59. 费芮网络
>60. 51CTO
>61. 聿旌科技
>62. 山东胜软科技股份有限公司
>63. 上海芯果科技有限公司(好特卖)
>64. 成都高鹿科技有限公司
>65. 飞视（苏州）数字技术有限公司
>66. 上海幻析信息科技有限公司
>67. 统信软件技术有限公司
>68. 得物
>69. 鼎翰文化股份有限公司
>70. 茶码纹化（云南）科技发展有限公司
>71. 湖南度思信息技术有限公司
>72. 深圳圆度
>73. 武汉沃柒科技有限公司(茄椒)
>74. 驭势科技
>75. 叮当跳动
>76. Keep
>77. simba innovation
>78. ZeroCMF
>79. 安徽寻梦投资发展集团
>80. 广州腾思信息科技有限公司
>81. 广州机智云物联网科技有限公司
>82. 厦门亿联网络技术股份有限公司
>83. 北京麦芽田网络科技有限公司
>84. 佛山市振联科技有限公司
>85. 苏州智言信息科技有限公司
>86. 中国移动上海产业研究院
>87. 天枢数链（浙江）科技有限公司
>88. 北京娱人共享智能科技有限公司
>89. 北京数智方科技有限公司
>90. 元匠科技
>91. 宁波甬风信息科技有限公司
>92. 深圳市万佳安物联科技股份有限公司
>93. 武侯区编程之美软件开发工作室
>94. 西安交通大学智慧能源与碳中和研究中心
>95. 成都创软科技有限责任公司
>96. Sonderbase Technologies
>97. 上海荣时信息科技有限公司
>98. 上海同犀智能科技有限公司
>99. 新华三技术有限公司
>100. 上海邑脉科技有限公司
>101. 上海巨瓴科技有限公司
>102. 深圳市兴海物联科技有限公司
>103. 爱芯元智半导体股份有限公司
>104. 杭州升恒科技有限公司
>105. 昆仑万维科技股份有限公司
>106. 无锡盛算信息技术有限公司
>107. 深圳市聚货通信息科技有限公司
>108. 浙江银盾云科技有限公司
>109. 南京造世网络科技有限公司
>110. 温州飞儿云信息技术有限公司
>111. 统信软件
>112. 深圳坐标软件集团有限公司

如果贵公司也已使用 go-zero，欢迎在 [登记地址](https://github.com/zeromicro/go-zero/issues/602) 登记，仅仅为了推广，不做其它用途。

## CNCF 云原生技术全景图

<p float="left">
<img src="https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/cncf-logo.svg" width="200"/>&nbsp;&nbsp;&nbsp;
<img src="https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/cncf-landscape-logo.svg" width="150"/>
</p>

go-zero 收录在 [CNCF Cloud Native 云原生技术全景图](https://landscape.cncf.io/?selected=go-zero)。

## 微信公众号

`go-zero` 相关文章和视频都会在 `微服务实践` 公众号整理呈现，欢迎扫码关注 👏

<img src="https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/zeromicro.jpg" alt="wechat" width="600" />

## 微信交流群

如果文档中未能覆盖的任何疑问，欢迎您在群里提出，我们会尽快答复。

您可以在群内提出使用中需要改进的地方，我们会考虑合理性并尽快修改。

如果您发现 ***bug*** 请及时提 ***issue***，我们会尽快确认并修改。

加群之前有劳点一下 ***star***，一个小小的 ***star*** 是作者们回答海量问题的动力！🤝

<img src="https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/wechat.jpg" alt="wechat" width="300" />