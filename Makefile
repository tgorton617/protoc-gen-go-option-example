default: clean example

install:
	go install .

example: install
	mkdir -p ./generatedresult/; \
	protoc -I=./proto \
		--option-example_out=./generatedresult/ \
		proto/options.proto ;

clean:
	rm -Rf generatedresult

.PHONY: install example clean