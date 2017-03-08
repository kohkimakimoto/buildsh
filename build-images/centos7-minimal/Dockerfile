FROM centos:centos7

# docker yum repo
# https://docs.docker.com/engine/installation/linux/centos/
RUN echo $'[dockerrepo]\n\
name=Docker Repository\n\
baseurl=https://yum.dockerproject.org/repo/main/centos/7/\n\
enabled=1\n\
gpgcheck=1\n\
gpgkey=https://yum.dockerproject.org/gpg' > /etc/yum.repos.d/docker.repo

# deps
RUN yum -y install \
  epel-release \
  && yum clean all && \
  yum -y install \
  sudo \
  docker-engine \
  && mkdir -p /build \
  && chmod -R 777 "/build"

# sudo config
RUN sed -i -e 's/Defaults    env_reset/#Defaults    env_reset/' /etc/sudoers \
    && sed -i -e 's/Defaults    secure_path/#Defaults    secure_path/' /etc/sudoers \
    && echo "%sudo	ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers \
    && echo 'Defaults    env_keep += "PATH"' >> /etc/sudoers

COPY entrypoint.sh /entrypoint.sh

WORKDIR /build
ENTRYPOINT ["/entrypoint.sh"]
