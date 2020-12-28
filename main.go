package main

import (
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"	
)

func main() {
	protogen.Options{
	}.Run(func(gen *protogen.Plugin) error {

		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

		for _, sourceFile := range gen.Files {
			if !sourceFile.Generate {
				continue
			}

			// setup output file
			outputfile := gen.NewGeneratedFile("./out.txt", sourceFile.GoImportPath)
			
			// pull out the particular message we know is in the file for this example
			options := sourceFile.Messages[0].Desc.Options()

			// how to get the custom options values out of the 'options' protoreflect.ProtoMessage?
			// just print the value to prove they are in there
			outputfile.P(fmt.Sprintf("Options value is:%v",options))

		}
		return nil
	})
}
