data "manifestory-my-datasource" "test" {
  mock = "mock-config"
}

locals {
  foo = data.manifestory-my-datasource.test.foo
  bar = data.manifestory-my-datasource.test.bar
}

source "null" "basic-example" {
  communicator = "none"
}

build {
  sources = [
    "source.null.basic-example"
  ]

  provisioner "shell-local" {
    inline = [
      "echo foo: ${local.foo}",
      "echo bar: ${local.bar}",
    ]
  }
}
