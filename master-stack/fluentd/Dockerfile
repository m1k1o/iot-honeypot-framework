ARG FLUENTD_VERSION=latest
FROM fluent/fluentd:${FLUENTD_VERSION}

USER root

# https://docs.fluentd.org/output/elasticsearch
RUN gem install fluent-plugin-elasticsearch --no-rdoc --no-ri; \
    gem install fluent-plugin-multi-format-parser

ENTRYPOINT ["fluentd", "-c", "/fluentd/etc/fluent.conf"]
