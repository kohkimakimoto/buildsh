FROM centos:centos7

# deps
RUN yum -y install \
  epel-release \
  && yum clean all && \
  yum -y install \
  sudo

# sudo config
RUN sed -i -e 's/Defaults    env_reset/#Defaults    env_reset/' /etc/sudoers \
    && sed -i -e 's/Defaults    secure_path/#Defaults    secure_path/' /etc/sudoers \
    && echo 'Defaults    env_keep += "PATH"' >> /etc/sudoers

WORKDIR /build
CMD ["/bin/bash"]
