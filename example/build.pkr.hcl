source "manifestory-builder" "foo" {
  mock = local.foo
}

source "manifestory-builder" "bar-example" {
  mock = local.bar
}

build {
  sources = [
    "source.manifestory-builder.foo",
  ]

  provisioner "manifestory-provisioner" {
    only = ["manifestory-builder.foo-example"]
    mock = "foo: ${local.foo}"
  }

  post-processor "manifestory-post-processor" {}
}
