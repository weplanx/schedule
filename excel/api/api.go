package api

import (
	"bytes"
	"context"
	"excel/common"
	"fmt"
	"github.com/bytedance/sonic/decoder"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/xuri/excelize/v2"
	"io"
	"net/http"
	"strings"
	"time"
)

type API struct {
	*common.Inject
}

type M map[string]interface{}

type Result struct {
	Records []Record `json:"records"`
}

type Record struct {
	Cos   `json:"cos"`
	Event `json:"event"`
}

type Cos struct {
	CosSchemaVersion  string `json:"cosSchemaVersion"`
	CosObject         `json:"cosObject"`
	CosBucket         `json:"cosBucket"`
	CosNotificationId string `json:"cosNotificationId"`
}

type CosObject struct {
	Url  string `json:"url"`
	Meta M      `json:"meta"`
	Vid  string `json:"vid"`
	Key  string `json:"key"`
	Size int64  `json:"size"`
}

type CosBucket struct {
	Region string `json:"region"`
	Name   string `json:"name"`
	Appid  string `json:"appid"`
}

type Event struct {
	EventName         string `json:"eventName"`
	EventVersion      string `json:"eventVersion"`
	EventTime         int64  `json:"eventTime"`
	EventSource       string `json:"eventSource"`
	RequestParameters M      `json:"requestParameters"`
	EventQueue        string `json:"eventQueue"`
	ReservedInfo      string `json:"reservedInfo"`
	Reqid             int64  `json:"reqid"`
}

func (x *API) EventInvoke(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		return
	}
	ctx := req.Context()

	var result Result
	if err := decoder.
		NewStreamDecoder(req.Body).
		Decode(&result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, record := range result.Records {
		prefix := fmt.Sprintf(`/%s/%s/`, record.CosBucket.Appid, record.CosBucket.Name)
		key := strings.Replace(record.Cos.Key, prefix, "", -1)
		resp, err := x.Client.Object.Get(ctx, key, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err = x.toExcel(ctx, resp.Body); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`已触发: %s`, time.Now())))
}

type Metadata struct {
	Name  string   `msgpack:"name"`
	Parts []string `msgpack:"parts"`
}

func (x *API) toExcel(ctx context.Context, body io.Reader) (err error) {
	var metadata Metadata
	if err = msgpack.NewDecoder(body).Decode(&metadata); err != nil {
		return
	}
	file := excelize.NewFile()
	defer file.Close()
	for _, key := range metadata.Parts {
		var streamWriter *excelize.StreamWriter
		args := strings.Split(key, ".")
		if streamWriter, err = file.NewStreamWriter(args[1]); err != nil {
			return
		}
		var resp *cos.Response
		if resp, err = x.Client.Object.Get(ctx, key, nil); err != nil {
			return
		}
		dec := msgpack.NewDecoder(resp.Body)
		rowID := 1
		for {
			var data []interface{}
			if err = dec.Decode(&data); err != nil {
				if err == io.EOF {
					break
				}
				return
			}
			cell, _ := excelize.CoordinatesToCellName(1, rowID)
			if err = streamWriter.SetRow(cell, data); err != nil {
				return
			}
			rowID++
		}
		if err = streamWriter.Flush(); err != nil {
			return
		}
	}
	var buf *bytes.Buffer
	if buf, err = file.WriteToBuffer(); err != nil {
		return
	}
	key := fmt.Sprintf(`%s.xlsx`, metadata.Name)
	if _, err = x.Client.Object.Put(ctx, key, buf, nil); err != nil {
		return
	}
	return
}
