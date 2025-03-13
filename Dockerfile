# 使用指定的基础镜像
FROM ghcr.nju.edu.cn/dendianugerah/reubah:latest

# 将本地的 components、js 和 pages 目录复制到容器内的 /app/templates 目录下，并设置权限为 777
COPY --chown=appuser:appuser --chmod=777 components /app/templates/components
COPY --chown=appuser:appuser --chmod=777 js /app/templates/js
COPY --chown=appuser:appuser --chmod=777 pages /app/templates/pages

# docker build -t firfe/reubah:2025.01.19 .