ARG K8S_PLUGIN_NAME=vault-k8s-secret-engine
ARG GITHUB_REPO_HTTPS_URL=https://github.com/dzirg44/vault-k8s-secret-engine
ARG VAULT_IMAGE_VERSION
FROM vault:"${VAULT_IMAGE_VERSION}"
LABEL org.opencontainers.image.source=${GITHUB_REPO_HTTPS_URL}
RUN mkdir -p /home/vault/plugins
COPY bin/vault-k8s-secret-engine /home/vault/plugins/${K8S_PLUGIN_NAME}
RUN echo $(sha256sum /home/vault/plugins/${K8S_PLUGIN_NAME} | cut -d ' ' -f1) > /home/vault/plugins/SHA256SUMS
RUN chown -R vault: /home/vault/plugins && chmod 755 /home/vault/plugins/${K8S_PLUGIN_NAME}
RUN setcap cap_ipc_lock=+ep /home/vault/plugins/${K8S_PLUGIN_NAME}
