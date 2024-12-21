> 翻译的不是源代码,是编译后的 `templates`文件

这是项目[reubah](https://github.com/dendianugerah/reubah)的汉化。  
原项目目前没有发布明确的版本，这个是2024年12月20日翻译的，在这个时间后如果原项目中的`templates`目录有修改，这个汉化可能不能使用，请注意甄别！！


## 翻译说明
首先感谢原作者的开源。  
本人不会英语，翻译此项目，只是自己有需求，但是发现这个项目没有提供中文，因此翻译了此项目；  
项目的所有英文翻译都是用 **通义千问** 通过问答的形式翻译的，如有不准确，请指出；  
本人提供这个项目在 NAS、服务器等的有偿远程部署服务，15元/次。有意者可联系。  
微信号 `E-0_0-` 二维码在最后  
闲鱼搜索 `明月人间`

具体翻译了哪些内容,请阅读[翻译说明](./翻译说明.md)

# 使用说明
下载本项目的`24.12.20`的发行版文件,将其中的`components`、`js`、`pages`目录映射到容器中的`/app/templates/`目录中的同名目录。

挂载示例
```yml
services:
  reubah:
    container_name: reubah
    # ghcr.nju.edu.cn 是国内的ghcr.io镜像地址
    image: ghcr.nju.edu.cn/dendianugerah/reubah:latest
    network_mode: bridge
    restart: always
    tty: true
    stdin_open: true
    cpus: 1
    mem_limit: 512m
    logging:
      driver: json-file
      options:
        max-size: 1m
        max-file: "3"
    environment:
      TZ: Asia/Shanghai
      TIME_ZONE: Asia/Shanghai
      PORT: 端口
    ports:
      - 端口:端口
    volumes:
      - 缓存:/app/tmp
      - 文档缓存:/tmp
      # 下面3个是汉化文件的目录
      - /templates/components:/app/templates/components
      - /templates/js:/app/templates/js
      - /templates/pages:/app/templates/pages
```

# 微信

![微信](./微信.png)