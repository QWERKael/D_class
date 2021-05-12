# atk_D_class
一个 基于命令行-插件式-低依赖 的服务器运维工具
## 设计目标
D class 是一个以 `低依赖` 和 `高适用性` 为设计目标的运维工具
## 设计理念
D class 有三个设计理念: `插件化` `低依赖` `命令行友好`
> 插件化: 插件化作为 D class 最核心的功能, 它代表的是我对 D class `轻量` `可扩展` `高度可定制` 的期望.

> 低依赖: 在运维中碰到的最恼人的事情莫过于各种命令在不同的环境下没办法正常运行了, 所以 D class 一开始就以低依赖作为其最重要的设计理念之一. 同时, 这也就是我为什么使用Go作为 D class 的核心开发语言的原因了.

> 命令行友好: 命令行友好的理念其实和低依赖是相辅相成的, 我们有很多的工具都有很好的web界面, 这很方便我们的日常使用, 但实际工作中, 我们需要处理故障的机器并没有提供我们完善的web管理环境, 也来不及部署这样的环境, 这个时候, 我希望 D class 能够作为一个有效的工具来为我们的工作提供高效的帮助.

## 架构与功能
D class 是最典型的C/S架构
### Server端
D class 的server端本体只提供最简单的三个功能, 1.重启 2.上传(是的, 连下载都不提供) 3.插件执行.
### Client端
D class 的client端提供了可配置的命令行提示功能, 可以将常用的和复杂的命令记录到配置文件, 在命令行输入几个字母就可以自动补全
### 插件
插件作为最核心的功能, 可以由用户自定义, 目前提供一些常用的插件供使用
#### show
#### cmd
#### mysql
#### watcher/sentry
#### async

### 默认端口管理
8881: server端默认端口
8880: watcher插件的默认端口
8883: transfer插件的默认端口