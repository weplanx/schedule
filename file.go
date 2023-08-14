package workflow

import (
	"bytes"
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/weplanx/workflow/typ"
	"net/http"
	"net/url"
)

type File struct {
	Cos *cos.Client
}

func (x *Workflow) NewFile(address string, id string, key string) (f *File, err error) {
	f = new(File)
	var u *url.URL
	if u, err = url.Parse(address); err != nil {
		return
	}
	baseURL := &cos.BaseURL{BucketURL: u}
	f.Cos = cos.NewClient(baseURL, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  id,
			SecretKey: key,
		},
	})
	return
}

type ExcelMetadata struct {
	Name  string   `msgpack:"name"`
	Parts []string `msgpack:"parts"`
}

func (x *File) Excel(ctx context.Context, name string, sheets typ.ExcelSheets) (err error) {
	metadata := ExcelMetadata{
		Name:  name,
		Parts: []string{},
	}
	for sheet, data := range sheets {
		w := bytes.NewBuffer(nil)
		enc := msgpack.NewEncoder(w)
		for _, v := range data {
			if err = enc.Encode(v); err != nil {
				return
			}
		}
		key := fmt.Sprintf(`%s.%s.pack`, name, sheet)
		metadata.Parts = append(metadata.Parts, key)
		if _, err = x.Cos.Object.Put(ctx, key, w, nil); err != nil {
			return
		}
	}

	var b []byte
	if b, err = msgpack.Marshal(metadata); err != nil {
		return
	}
	key := fmt.Sprintf(`%s.excel`, name)
	if _, err = x.Cos.Object.Put(ctx, key, bytes.NewBuffer(b), nil); err != nil {
		return
	}
	return
}
