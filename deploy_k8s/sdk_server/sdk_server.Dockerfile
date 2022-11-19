FROM ubuntu

# 设置固定的项目路径
ENV WORKDIR /FreeChat_Service
ENV CMDDIR $WORKDIR/cmd
ENV CONFIG_NAME $WORKDIR/cmd/main/config.yaml

# 将可执行文件复制到目标目录
ADD ./open_im_sdk_server $WORKDIR/main
ADD ../config/config.yaml $WORKDIR/cmd/main

# 创建用于挂载的几个目录，添加可执行权限
RUN mkdir $WORKDIR/logs $WORKDIR/config $WORKDIR/db && \
  chmod +x $WORKDIR/main

VOLUME ["/FreeChat_Service/logs","/FreeChat_Service/config","/FreeChat_Service/script","/FreeChat_Service/db/sdk"]

WORKDIR $CMDDIR
CMD ./main