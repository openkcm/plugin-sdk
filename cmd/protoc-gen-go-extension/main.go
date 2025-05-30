package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

const (
	pluginsdkPackage = protogen.GoImportPath("github.com/openkcm/plugin-sdk/api")
	grpcPackage      = protogen.GoImportPath("google.golang.org/grpc")
)

var (
	flags     flag.FlagSet
	kind      = flags.String("kind", "plugin", `generation kind (either "plugin" or "service")`)
	submodule = flags.String("submodule", "", `package location`)
)

func main() {
	protogen.Options{ParamFunc: flags.Set}.Run(func(gen *protogen.Plugin) error {
		isPlugin := false
		switch *kind {
		case "service":
		case "plugin":
			isPlugin = true
		default:
			return fmt.Errorf(`invalid kind %q: expecting either "plugin" or "service"`, *kind)
		}
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			if !strings.Contains(f.GoImportPath.String(), *submodule) {
				continue
			}
			if _, err := generateFile(gen, f, isPlugin); err != nil {
				return err
			}
		}
		return nil
	})
}

func generateFile(gen *protogen.Plugin, file *protogen.File, isPlugin bool) (*protogen.GeneratedFile, error) {
	if len(file.Services) == 0 {
		return nil, nil
	}

	filename := file.GeneratedFilenamePrefix
	if isPlugin {
		filename += "_ext_plugin.pb.go"
	} else {
		filename += "_ext_service.pb.go"
	}

	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P("// Code generated by protoc-gen-go-extension. DO NOT EDIT.")
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()
	for _, service := range file.Services {
		serviceGoName := service.GoName
		serviceFullName := string(service.Desc.FullName())
		err := generateServiceBridges(g, serviceGoName, serviceFullName, isPlugin)
		if err != nil {
			return nil, err
		}
	}
	return g, nil
}

func generateServiceBridges(g *protogen.GeneratedFile, serviceName, serviceFullName string, isPlugin bool) error {
	if serviceName == "" {
		return fmt.Errorf("generateServiceBridges(): serviceGoName is empty")
	}
	if serviceFullName == "" {
		return fmt.Errorf("generateServiceBridges(): serviceFullName is empty")
	}

	kind := "Service"
	if isPlugin {
		kind = "Plugin"
	}

	serverIntfName := serviceName + "Server"
	pluginServerCons := serviceName + kind + "Server"
	pluginServerType := unexport(pluginServerCons)
	pluginServerIdent := g.QualifiedGoIdent(pluginsdkPackage.Ident(kind + "Server"))

	clientIntfName := serviceName + "Client"
	pluginClientType := serviceName + kind + "Client"

	g.P()
	g.P("const (")
	if isPlugin {
		g.P("	Type = ", strconv.Quote(serviceName))
	}
	g.P("	GRPCServiceFullName = ", strconv.Quote(serviceFullName))
	g.P(")")
	g.P()
	g.P("func ", pluginServerCons, "(server ", serverIntfName, ") ", pluginServerIdent, " {")
	g.P("return ", pluginServerType, "{", serverIntfName, ": server}")
	g.P("}")
	g.P()
	g.P("type ", pluginServerType, " struct {")
	g.P(serverIntfName)
	g.P("}")
	if isPlugin {
		g.P()
		g.P("func (s ", pluginServerType, ") Type() string {")
		g.P("return Type")
		g.P("}")
	}
	g.P()
	g.P("func (s ", pluginServerType, ") GRPCServiceName() string {")
	g.P("return GRPCServiceFullName")
	g.P("}")
	g.P()
	g.P("func (s ", pluginServerType, ") RegisterServer(server *", g.QualifiedGoIdent(grpcPackage.Ident("Server")), ") any {")
	g.P("Register", serverIntfName, "(server , s.", serverIntfName, ")")
	g.P("return s.", serverIntfName)
	g.P("}")

	g.P()
	g.P("type ", pluginClientType, " struct {")
	g.P(clientIntfName)
	g.P("}")
	if isPlugin {
		g.P()
		g.P("func (s ", pluginClientType, ") Type() string {")
		g.P("return Type")
		g.P("}")
	}
	g.P()
	g.P("func (c *", pluginClientType, ") IsInitialized() bool {")
	g.P("return c.", clientIntfName, " != nil")
	g.P("}")
	g.P()
	g.P("func (c *", pluginClientType, ") GRPCServiceName() string {")
	g.P("return GRPCServiceFullName")
	g.P("}")
	g.P()
	g.P("func (c *", pluginClientType, ") InitClient(conn ", g.QualifiedGoIdent(grpcPackage.Ident("ClientConnInterface")), ") any {")
	g.P("c.", clientIntfName, " = New", clientIntfName, "(conn)")
	g.P("return c.", clientIntfName)
	g.P("}")
	return nil
}

func unexport(s string) string {
	if len(s) < 1 {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}
