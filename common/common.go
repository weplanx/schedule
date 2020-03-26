package common

import (
	"encoding/json"
	socketio "github.com/googollee/go-socket.io"
	"github.com/streadway/amqp"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	LogOpt      *LogOption
	Socket      *socketio.Server
	SocketConn  *socketio.Conn
	AmqpConn    *amqp.Connection
	AmqpChannel *amqp.Channel
)

type (
	AppOption struct {
		Debug  bool      `yaml:"debug"`
		Listen string    `yaml:"listen"`
		Log    LogOption `yaml:"log"`
	}
	LogOption struct {
		Storage        bool   `yaml:"storage"`
		StorageDir     string `yaml:"storage_dir"`
		Socket         bool   `yaml:"socket"`
		SocketPort     string `yaml:"socket_port"`
		Amqp           bool   `yaml:"amqp"`
		AmqpUri        string `yaml:"amqp_uri"`
		AmqpExchange   string `yaml:"amqp_exchange"`
		AmqpRoutingKey string `yaml:"amqp_routing_key"`
	}
	JobOption struct {
		Identity string                  `yaml:"identity"`
		TimeZone string                  `yaml:"time_zone"`
		Start    bool                    `yaml:"start"`
		Entries  map[string]*EntryOption `yaml:"entries"`
	}
	EntryOption struct {
		CronTime string            `yaml:"cron_time"`
		Url      string            `yaml:"url"`
		Headers  map[string]string `yaml:"headers"`
		Body     interface{}       `yaml:"body"`
		NextDate time.Time         `yaml:"-"`
		LastDate time.Time         `yaml:"-"`
	}
)

func autoload(identity string) string {
	return "./config/autoload/" + identity + ".yml"
}

func ListConfig() (list []JobOption, err error) {
	var files []os.FileInfo
	files, err = ioutil.ReadDir("./config/autoload")
	if err != nil {
		return
	}
	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if ext == ".yml" {
			var in []byte
			in, err = ioutil.ReadFile("./config/autoload/" + file.Name())
			if err != nil {
				return
			}
			var config JobOption
			err = yaml.Unmarshal(in, &config)
			if err != nil {
				return
			}
			list = append(list, config)
		}
	}
	return
}

func SaveConfig(data *JobOption) (err error) {
	var out []byte
	out, err = yaml.Marshal(data)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(
		autoload(data.Identity),
		out,
		0644,
	)
	if err != nil {
		return
	}
	return
}

func RemoveConfig(identity string) error {
	return os.Remove(autoload(identity))
}

func SetLogger(option *LogOption) (err error) {
	LogOpt = option
	if _, err := os.Stat(option.StorageDir); os.IsNotExist(err) {
		os.Mkdir(option.StorageDir, os.ModeDir)
	}
	if LogOpt.Socket {
		go func() {
			Socket, err = socketio.NewServer(nil)
			if err != nil {
				return
			}
			Socket.OnConnect("/", func(s socketio.Conn) error {
				SocketConn = &s
				return nil
			})
			go Socket.Serve()
			http.Handle("/socket.io/", Socket)
			http.ListenAndServe(":"+LogOpt.SocketPort, nil)
		}()
	}
	if LogOpt.Amqp {
		AmqpConn, err = amqp.Dial(LogOpt.AmqpUri)
		if err != nil {
			return
		}
		AmqpChannel, err = AmqpConn.Channel()
		if err != nil {
			return
		}
	}
	return
}

func LoggerClose() {
	if LogOpt.Socket && Socket != nil {
		Socket.Close()
	}
	if LogOpt.Amqp {
		if AmqpChannel != nil {
			AmqpChannel.Close()
		}
		if AmqpConn != nil {
			AmqpConn.Close()
		}
	}
}

func PushLogger(v interface{}) (err error) {
	if LogOpt.Socket && SocketConn != nil {
		(*SocketConn).Emit("logger", v)
	}
	if LogOpt.Amqp && AmqpChannel != nil {
		var body []byte
		body, err = json.Marshal(v)
		if err != nil {
			return
		}
		err = AmqpChannel.Publish(
			LogOpt.AmqpExchange,
			LogOpt.AmqpRoutingKey,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			},
		)
		if err != nil {
			return
		}
	}
	return
}

func OpenStorage() bool {
	return LogOpt.Storage
}

func LogFile(identity string) (file *os.File, err error) {
	if _, err := os.Stat("./" + LogOpt.StorageDir + "/" + identity); os.IsNotExist(err) {
		os.Mkdir("./"+LogOpt.StorageDir+"/"+identity, os.ModeDir)
	}
	date := time.Now().Format("2006-01-02")
	filename := "./" + LogOpt.StorageDir + "/" + identity + "/" + date + ".log"
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		file, err = os.Create(filename)
		if err != nil {
			return
		}
	} else {
		file, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			return
		}
	}
	return
}
