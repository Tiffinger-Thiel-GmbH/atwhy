---
meta:
  title: Architecture 
server:
  index: true
---
# Architecture

The idea behind atwhy is to have several interfaces, each one for a small purpose. These interfaces are then
implemented by concrete implementations.

It is possible to replace or mock each part of the application at any time. You can use atwhy as lib and just provide
your own implementations.

The interfaces are:
{{ .Tag.atwhy_interfaces }}  
{{ .Tag.atwhy_struct_link }}  
{{ .Tag.atwhy_struct_code }}