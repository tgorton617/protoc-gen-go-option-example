# protoc-gen-go-option-example

Demonstration of how to use custom options in a golang compiler/protogen protoc plugin using the solution provided by @dsnet at [golang/protobuf issue # 1260](https://github.com/golang/protobuf/issues/1260). Thanks @dsnet!

## How to run

First you need protoc installed (and, of course, go.)

Then run:

`make`

The output (in ./generatedresult/out.txt) should look like:

```text

Message Foo:
Value of option example_annotation_int32 is 1234
Value of option example_annotation_string is hello

Message Bar:
Value of option example_annotation_string is world
```
