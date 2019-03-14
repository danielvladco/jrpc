package main

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/gogo/protobuf/vanity"
)

func main() {
	gen := generator.New()

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		gen.Error(err, "reading input")
	}

	if err := proto.Unmarshal(data, gen.Request); err != nil {
		gen.Error(err, "parsing input proto")
	}

	if len(gen.Request.FileToGenerate) == 0 {
		gen.Fail("no files to generate")
	}

	useGogoImport := false
	// Match parsing algorithm from Generator.CommandLineParameters
	for _, parameter := range strings.Split(gen.Request.GetParameter(), ",") {
		kvp := strings.SplitN(parameter, "=", 2)
		// We only care about key-value pairs where the key is "gogoimport"
		if len(kvp) != 2 || kvp[0] != "gogoimport" {
			continue
		}
		useGogoImport, err = strconv.ParseBool(kvp[1])
		if err != nil {
			gen.Error(err, "parsing gogoimport option")
		}
	}

	gen.CommandLineParameters(gen.Request.GetParameter())
	gen.WrapTypes()
	gen.SetPackageNames()
	gen.BuildTypeNameMap()
	gen.GeneratePlugin(&plugin{useGogoImport: useGogoImport, enums: make(map[string]struct{}), messages: make(map[string]struct{})})

	for i := 0; i < len(gen.Response.File); i++ {
		gen.Response.File[i].Name = proto.String(strings.Replace(*gen.Response.File[i].Name, ".pb.go", ".jrpc.pb.go", -1))
	}

	// Send back the results.
	data, err = proto.Marshal(gen.Response)
	if err != nil {
		gen.Error(err, "failed to marshal output proto")
	}
	_, err = os.Stdout.Write(data)
	if err != nil {
		gen.Error(err, "failed to write output proto")
	}
}

type plugin struct {
	*generator.Generator
	generator.PluginImports

	enums    map[string]struct{}
	messages map[string]struct{}

	useGogoImport bool
}

func (p *plugin) Name() string                { return "gqlenum" }
func (p *plugin) Init(g *generator.Generator) { p.Generator = g }
func (p *plugin) Generate(file *generator.FileDescriptor) {
	if !p.useGogoImport {
		vanity.TurnOffGogoImport(file.FileDescriptorProto)
	}
	p.PluginImports = generator.NewPluginImports(p.Generator)
	http := p.NewImport("net/http")
	fmt := p.NewImport("fmt")
	json := p.NewImport("encoding/json")
	status := p.NewImport("google.golang.org/grpc/status")
	codes := p.NewImport("google.golang.org/grpc/codes")

	for _, r := range p.Request.FileToGenerate {
		if r == file.GetName() {
			for _, svc := range file.GetService() {

				p.P(`
func `, generator.CamelCase(svc.GetName()), `HTTPServer(svc `, generator.CamelCase(svc.GetName()), `Server) `, http.Use(), `.Handler {
	mux := `, http.Use(), `.NewServeMux()`)
				for _, rpc := range svc.GetMethod() {
					obj := p.TypeNameByObject(rpc.GetInputType())
					pkg := p.DefaultPackageName(obj)
					if pkg != "" {
						pkg = p.NewImport(string(obj.GoImportPath())).Use() + "."
					}
					goType := pkg + generator.CamelCaseSlice(obj.TypeName())
					reqPath := "/" + svc.GetName() + "/" + rpc.GetName()
					if rpc.GetServerStreaming() || rpc.GetClientStreaming() {
						continue
					}
					p.P(`mux.HandleFunc("`, reqPath, `", func (w `, http.Use(), `.ResponseWriter, r *`, http.Use(), `.Request) {
	req := &`, goType, `{}
	defer r.Body.Close()
	w.Header().Set("Content-Type","application/json")
	if err := `, json.Use(), `.NewDecoder(r.Body).Decode(req); err != nil {
		fmt.Fprintf(w, "{%q:%q,%q:%q}", "error", err.Error(), "status", `, http.Use(), `.StatusText(`, http.Use(), `.StatusBadRequest))
		w.WriteHeader(`, http.Use(), `.StatusBadRequest)
		return
	}


	res, err := svc.`, generator.CamelCase(rpc.GetName()), `(r.Context(), req)
	if err != nil {
		stt, _ := `, status.Use(), `.FromError(err)
		st := map[`, codes.Use(), `.Code] int {
			`, codes.Use(), `.Canceled: 400,
			`, codes.Use(), `.Unknown: 500,
			`, codes.Use(), `.InvalidArgument: 400,
			`, codes.Use(), `.DeadlineExceeded: 503,
			`, codes.Use(), `.NotFound: 404,
			`, codes.Use(), `.AlreadyExists: 400,
			`, codes.Use(), `.PermissionDenied: 403,
			`, codes.Use(), `.ResourceExhausted: 503,
			`, codes.Use(), `.FailedPrecondition: 400,
			`, codes.Use(), `.Aborted: 400,
			`, codes.Use(), `.OutOfRange: 500,
			`, codes.Use(), `.Unimplemented: 404,
			`, codes.Use(), `.Internal: 500,
			`, codes.Use(), `.Unavailable: 503,
			`, codes.Use(), `.DataLoss: 500,
			`, codes.Use(), `.Unauthenticated: 401,
		}[stt.Code()]
		b, _ := `, json.Use(), `.Marshal(stt.Details())
		`, fmt.Use(), `.Fprintf(w, "{%q:%q,%q:%q, %q: %s}", "error", err.Error(), "status", stt.Code().String(), "details", b)
		w.WriteHeader(st)
		return
	}

	if err := `, json.Use(), `.NewEncoder(w).Encode(res); err != nil {
		`, fmt.Use(), `.Fprintf(w, "{%q:%q,%q:%q}", "error", err.Error(), "status", `, http.Use(), `.StatusText(`, http.Use(), `.StatusInternalServerError))
		w.WriteHeader(`, http.Use(), `.StatusInternalServerError)
		return
	}
})`)
				}
				p.P("\nreturn mux }")
			}
		}
	}
}
