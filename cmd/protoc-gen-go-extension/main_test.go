package main

import (
	"testing"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func TestGenerateFile(t *testing.T) {
	// create test cases
	tests := []struct {
		name      string
		request   *pluginpb.CodeGeneratorRequest
		isPlugin  bool
		wantError bool
	}{
		{
			name: "without services",
			request: &pluginpb.CodeGeneratorRequest{
				ProtoFile: []*descriptorpb.FileDescriptorProto{
					{
						Name:    proto.String("testdata/foo/foo.proto"),
						Syntax:  proto.String(protoreflect.Proto3.String()),
						Package: proto.String("testdata.foo"),
						Options: &descriptorpb.FileOptions{
							GoPackage: proto.String("testdata.foo"),
						},
					},
				},
			},
		}, {
			name: "as service",
			request: &pluginpb.CodeGeneratorRequest{
				ProtoFile: []*descriptorpb.FileDescriptorProto{
					{
						Name:    proto.String("testdata/foo/foo.proto"),
						Syntax:  proto.String(protoreflect.Proto3.String()),
						Package: proto.String("testdata.foo"),
						Options: &descriptorpb.FileOptions{
							GoPackage: proto.String("testdata.foo"),
						},
						Service: []*descriptorpb.ServiceDescriptorProto{
							{
								Name: proto.String("DoMagic"),
							},
						},
					},
				},
			},
		}, {
			name:     "as plugin",
			isPlugin: true,
			request: &pluginpb.CodeGeneratorRequest{
				ProtoFile: []*descriptorpb.FileDescriptorProto{
					{
						Name:    proto.String("testdata/foo/foo.proto"),
						Syntax:  proto.String(protoreflect.Proto3.String()),
						Package: proto.String("testdata.foo"),
						Options: &descriptorpb.FileOptions{
							GoPackage: proto.String("testdata.foo"),
						},
						Service: []*descriptorpb.ServiceDescriptorProto{
							{
								Name: proto.String("DoMagic"),
							},
						},
					},
				},
			},
		},
	}

	// run the tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			gen, err := protogen.Options{}.New(tc.request)
			if err != nil {
				t.Fatal(err)
			}

			// Act
			_, err = generateFile(gen, gen.Files[0], tc.isPlugin)

			// Assert
			if tc.wantError && err != nil { // expected error and got it
				return
			} else if tc.wantError && err == nil { // expected error but did not get it
				t.Errorf("expected error value: %v, got: %s", tc.wantError, err)
			} else if !tc.wantError && err != nil { // got unexpected error
				t.Errorf("expected error value: %v, got: %s", tc.wantError, err)
			} else if !tc.wantError && err == nil {
				return
			}
		})
	}
}

func TestGenerateServiceBridges(t *testing.T) {
	// Arrange
	gen, err := protogen.Options{}.New(&pluginpb.CodeGeneratorRequest{})
	if err != nil {
		t.Fatal(err)
	}

	// create test cases
	tests := []struct {
		name              string
		serviceGoName     string
		serviceFullName   string
		isPlugin          bool
		wantGenerateError bool
		wantContentError  bool
	}{
		{
			name:              "zero values",
			wantGenerateError: true,
		}, {
			name:              "empty service name",
			serviceGoName:     "",
			serviceFullName:   "github.com/openkcm/foo.DoMagic",
			wantGenerateError: true,
		}, {
			name:              "zero values",
			serviceGoName:     "DoMagic",
			serviceFullName:   "",
			wantGenerateError: true,
		}, {
			name:              "valid service",
			serviceGoName:     "DoMagic",
			serviceFullName:   "github.com/openkcm/foo.DoMagic",
			wantGenerateError: false,
			wantContentError:  false,
		}, {
			name:              "valid plugin",
			serviceGoName:     "DoMagic",
			serviceFullName:   "github.com/openkcm/foo.DoMagic",
			isPlugin:          true,
			wantGenerateError: false,
			wantContentError:  false,
		},
	}

	// run the tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			g := gen.NewGeneratedFile("foo.go", "github.com/openkcm/foo")
			g.P("package ", "foo")

			// Act 1
			err = generateServiceBridges(g, tc.serviceGoName, tc.serviceFullName, tc.isPlugin)

			// Assert 1
			if tc.wantGenerateError && err != nil { // expected error and got it
				return
			} else if tc.wantGenerateError && err == nil { // expected error but did not get it
				t.Errorf("expected error value: %v, got: %s", tc.wantGenerateError, err)
			} else if !tc.wantGenerateError && err != nil { // got unexpected error
				t.Errorf("expected error value: %v, got: %s", tc.wantGenerateError, err)
			} else if !tc.wantGenerateError && err == nil {
				// Act 2
				_, err = g.Content()

				// Assert 2
				if tc.wantContentError && err != nil { // expected error and got it
					return
				} else if tc.wantContentError && err == nil { // expected error but did not get it
					t.Errorf("expected error value: %v, got: %s", tc.wantContentError, err)
				} else if !tc.wantContentError && err != nil { // got unexpected error
					t.Errorf("expected error value: %v, got: %s", tc.wantContentError, err)
				} else if !tc.wantContentError && err == nil {
					return
				}
			}
		})
	}
}

func TestUnexport(t *testing.T) {
	// create test cases
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name: "zero values",
		}, {
			name:  "lowercase",
			input: "fooBar",
			want:  "fooBar",
		}, {
			name:  "uppercase",
			input: "FooBar",
			want:  "fooBar",
		},
	}

	// run the tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			got := unexport(tc.input)

			// Assert
			if got != tc.want {
				t.Errorf("unexport() = %v, want %v", got, tc.want)
			}
		})
	}
}
