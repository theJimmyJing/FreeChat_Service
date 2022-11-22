FROM ubuntu

# 设置固定的项目路径
ENV WORKDIR /FreeChat_Service
ENV CMDDIR $WORKDIR/cmd
ENV CONFIG_NAME $CMDDIR/config/config.yaml


# 将可执行文件复制到目标目录
ADD ./open_im_friend $WORKDIR/cmd/main
ADD ./config.yaml $WORKDIR/cmd/config/config.yaml

# 创建用于挂载的几个目录，添加可执行权限
RUN mkdir $WORKDIR/logs $WORKDIR/config $WORKDIR/script && \
  chmod +x $WORKDIR/cmd/main

VOLUME ["/FreeChat_Service/logs","/FreeChat_Service/config", "/FreeChat_Service/script"]

WORKDIR $CMDDIR
CMD ./main