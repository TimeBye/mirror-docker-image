# Mirror docker images

[![Build Status](https://travis-ci.com/TimeBye/mirror-docker-image.svg?branch=mirror)](https://travis-ci.com/TimeBye/mirror-docker-image)[![CircleCI](https://circleci.com/gh/TimeBye/mirror-docker-image.svg?style=svg)](https://circleci.com/gh/TimeBye/mirror-docker-image)


### 怎样进行`docker image`同步

1. 在CI文件中定义登录`Docker Registry`语句
2. 将需要同步的镜像列表写入`config.yaml`文件中

    ```yaml
    # 批量同步镜像至同一仓库
    batch:
        # 同时并发同步个数
        maxConcurrentDownloads: 5
        # 目标仓库
        targetRegistry: registry.saas.hand-china.com/tools
        # 镜像列表
        images:
        - "redis:5.0.3"
        - "mysql:5.7.24"
        - "nginx:1.15.7"
    # 单个镜像同步
    single:
      # 镜像地址
    - image: "tomcat:7.0.92-jre8"
      # 目标仓库
      targetRegistry: registry.saas.hand-china.com/tools
    ```