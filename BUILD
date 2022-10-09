load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:prefix github.com/todo-project
gazelle(name = "gazelle")

# disable proto things, so that it uses .pb.go files instead
# https://github.com/bazelbuild/rules_go/blob/master/proto/core.rst#option-2-use-pre-generated-pb-go-files
# gazelle:proto disable_global
