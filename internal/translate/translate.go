package translate

import (
	"TweakItDocs/internal/translate/gen"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"log"
	"os"
	"regexp"
)

func GetProto() *gen.Data {
	b, err := os.ReadFile("data.protobin")
	if err != nil {
		log.Fatal(err)
	}
	var data gen.Data
	err = proto.Unmarshal(b, &data)
	if err != nil {
		log.Fatal(err)
	}
	return &data
}

func Bro() {
	err := toProtoJSON("data.json")
	if err != nil {
		log.Fatal(err)
	}

	var r gen.Data
	b, err := os.ReadFile("proto-data.json")
	if err != nil {
		log.Fatal(err)
	}
	err = protojson.Unmarshal(b, &r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("yeeeee baby")
	out, err := proto.Marshal(&r)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile("data.protobin", out, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func toProtoJSON(filename string) error {
	b, err := os.ReadFile(filename)
	if err != nil {
		return errors.Wrap(err, "could not read the file")
	}
	pkgEXP := regexp.MustCompile(`(_index":\s*\{\s*"index":\s*.*\s*")reference(":\s*\{\s*"class_package)`)
	idxEXP := regexp.MustCompile(`(_index":\s*\{\s*"index":\s*.*\s*")reference(":\s*\{\s*"class_index)`)
	b = pkgEXP.ReplaceAll(b, []byte("${1}norm_reference${2}"))
	b = idxEXP.ReplaceAll(b, []byte("${1}object_reference${2}"))
	out := make([]byte, 0, len(b)+100)
	out = append(out, []byte(`{"packages":`)...)
	out = append(out, b...)
	out = append(out, []byte("}")...)
	err = os.WriteFile("proto-"+filename, out, 0644)
	return errors.Wrap(err, "could not write to the file")
}
