FROM alpine:3.2
ADD organization-srv /organization-srv
ENTRYPOINT [ "/organization-srv" ]
