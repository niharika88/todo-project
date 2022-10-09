FROM gcr.io/cloud-builders/bazel

WORKDIR /build
COPY . /build

RUN bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64  //cmd:image --@io_bazel_rules_docker//transitions:enable=no 