# Terraform Provider LocalNull

This is a terraform provider to run local commands using configuration
variables specified to the resource. This is useful when you want to
run a command only when certain variables change.

When a `configuration` variable or `command` changes, the command
will be run again and the output updated.

```hcl
resource "localnull_with_variables" "example-1" {
  configuration = {
    var1 = "hello"
    var2 = "world"
  }

  command = "echo {{.var1}} {{.var2}}"
}

output "example1-command" {
  value = "${localnull_with_variables.example-1.interpolated_command}"
}

output "example1-command-output" {
  value = "${localnull_with_variables.example-1.output}"
}
```
