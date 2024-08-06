FROM eu.gcr.io/fifth-cab-359408/grpc-base:latest
COPY bin/notification-manager /notification-manager
ENTRYPOINT [ "/notification-manager" ]