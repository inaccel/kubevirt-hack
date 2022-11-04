FROM scratch
COPY kubevirt-hack /bin/kubevirt-hack
ENTRYPOINT ["kubevirt-hack"]
