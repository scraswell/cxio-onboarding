FROM menci/archlinuxarm:base-devel

ENV NVM_DIR /usr/local/nvm
ENV MY_NODE_VERSION 16.4.2
ENV REMOTE_USER_NAME vscode

RUN echo -e "export NVM_DIR=\"/usr/local/nvm\"\\n[ -s \"/usr/local/nvm/nvm.sh\" ] && \\. \"/usr/local/nvm/nvm.sh\"" >> /etc/skel/.bashrc
RUN echo -e "eval \"\$(starship init bash)\"" >> /etc/skel/.bashrc

# TODO: Fix delve installation; it installs only in the context of the user running the go command.
#       Check go install for flags... or maybe just switch user contexts and run the command.
RUN mkdir -pv ${NVM_DIR} && \
        pacman -Syu --noconfirm && \
        pacman -S --noconfirm curl cloud-guest-utils git sudo openssh go go-tools gopls revive staticcheck yaegi && \
        go install github.com/go-delve/delve/cmd/dlv@latest && \
        curl -L -O https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.1/install.sh && \
        bash install.sh && \
        rm -v install.sh && \
        curl -sS -L -O https://starship.rs/install.sh && \
        sh install.sh --yes && \
        rm -v install.sh && \
        . $NVM_DIR/nvm.sh && \
        nvm install ${MY_NODE_VERSION} && \
        nvm use default && \
        npm install -g npm@8.6.0 && \
        npm install -g pnpm@next-7

RUN echo -e "%wheel    ALL=NOPASSWD: ALL" >> /etc/sudoers.d/00_wheel-no-passwd
RUN chmod 0440 /etc/sudoers.d/00_wheel-no-passwd
RUN cp -v /etc/skel/.bashrc /root/.bashrc

RUN useradd -m -s /bin/bash -d /home/${REMOTE_USER_NAME} ${REMOTE_USER_NAME} && \
        usermod -aG wheel ${REMOTE_USER_NAME}
