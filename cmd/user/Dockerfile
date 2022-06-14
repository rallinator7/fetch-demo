FROM ubuntu:latest as build
RUN useradd -u 10001 fetch
ADD . .
# copy artifacts to a clean image
FROM scratch
COPY --from=build ./main /main
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /bin/tar /bin/tar
USER fetch
ENTRYPOINT [ "/main" ]