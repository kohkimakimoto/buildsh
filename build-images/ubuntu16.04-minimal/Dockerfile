FROM ubuntu:16.04		
		
RUN apt-get -y update		
RUN apt-get -y install \		
  sudo \		
  && mkdir -p /build \		
  && chmod -R 777 "/build"		
	
# sudo config		
RUN sed -i -e 's/Defaults[ 	]env_reset/#Defaults    env_reset/' /etc/sudoers \		
    && sed -i -e 's/Defaults[ 	]secure_path/#Defaults    secure_path/' /etc/sudoers \		
    && echo 'Defaults    env_keep += "PATH"' >> /etc/sudoers		
		
WORKDIR /build
CMD ["/bin/bash"]