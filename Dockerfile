# 使用指定的基础镜像
FROM ghcr.nju.edu.cn/dendianugerah/reubah:latest

# 将本地的 components、js 和 pages 目录复制到容器内的 /app/templates 目录下
# 设置所有者uid=1000(appuser) gid=1000(appgroup) groups=1000(appgroup),1000(appgroup) 权限为 777
COPY --chown=appuser:appuser --chmod=777 components /app/templates/components
COPY --chown=appuser:appuser --chmod=777 js /app/templates/js
COPY --chown=appuser:appuser --chmod=777 pages /app/templates/pages

# docker build -t firfe/reubah:25.01.19 .