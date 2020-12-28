package main

import (
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	protogen.Options{
	}.Run(func(gen *protogen.Plugin) error {

		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

		// The type information for all extensions is in the source files,
		// so we need to extract them into a dynamically created protoregistry.Types.
		extTypes := new(protoregistry.Types)
		for _, file := range gen.Files {
			if err := registerAllExtensions(extTypes, file.Desc); err != nil {
				panic(err)
			}
		}

		// run through the files again, extracting and printing the Message options
		for _, sourceFile := range gen.Files {
			if !sourceFile.Generate {
				continue
			}

			// setup output file
			outputfile := gen.NewGeneratedFile("./out.txt", sourceFile.GoImportPath)

			for _, message := range sourceFile.Messages {
				outputfile.P(fmt.Sprintf("\nMessage %s:", message.Desc.Name()))

				// The MessageOptions as provided by protoc does not know about
				// dynamically created extensions, so they are left as unknown fields.
				// We round-trip marshal and unmarshal the options with
				// a dynamically created resolver that does know about extensions at runtime.
				options := message.Desc.Options().(*descriptorpb.MessageOptions)
				b, err := proto.Marshal(options)
				if err != nil {
					panic(err)
				}
				options.Reset()
				err = proto.UnmarshalOptions{Resolver: extTypes}.Unmarshal(b, options)
				if err != nil {
					panic(err)
				}

				// Use protobuf reflection to iterate over all the extension fields,
				// looking for the ones that we are interested in.
				options.ProtoReflect().Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
					if !fd.IsExtension() {
						return true
					}

					outputfile.P(fmt.Sprintf("Value of option %s is %s",fd.Name(), v.String()))

					// Make use of fd and v based on their reflective properties.

					return true
				})
			}
		}
		
		return nil
	})
}

// Recursively register all extensions into the provided protoregistry.Types,
// starting with the protoreflect.FileDescriptor and recursing into its MessageDescriptors,
// their nested MessageDescriptors, and so on.
//
// This leverages the fact that both protoreflect.FileDescriptor and protoreflect.MessageDescriptor 
// have identical Messages() and Extensions() functions in order to recurse through a single function
func registerAllExtensions(extTypes *protoregistry.Types, descs interface {
	Messages() protoreflect.MessageDescriptors
	Extensions() protoreflect.ExtensionDescriptors
}) error {
	mds := descs.Messages()
	for i := 0; i < mds.Len(); i++ {
		registerAllExtensions(extTypes, mds.Get(i))
	}
	xds := descs.Extensions()
	for i := 0; i < xds.Len(); i++ {
		if err := extTypes.RegisterExtension(dynamicpb.NewExtensionType(xds.Get(i))); err != nil {
			return err
		}
	}
	return nil
}