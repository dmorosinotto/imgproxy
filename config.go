package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
)

func intEnvConfig(i *int, name string) {
	if env, err := strconv.Atoi(os.Getenv(name)); err == nil {
		*i = env
	}
}

func strEnvConfig(s *string, name string) {
	if env := os.Getenv(name); len(env) > 0 {
		*s = env
	}
}

func hexEnvConfig(b *[]byte, name string) {
	var err error

	if env := os.Getenv(name); len(env) > 0 {
		if *b, err = hex.DecodeString(env); err != nil {
			log.Fatalf("%s expected to be hex-encoded string\n", name)
		}
	}
}

func hexFileConfig(b *[]byte, filepath string) {
	if len(filepath) == 0 {
		return
	}

	f, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Can't open file %s\n", filepath)
	}

	src, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalln(err)
	}

	src = bytes.TrimSpace(src)

	dst := make([]byte, hex.DecodedLen(len(src)))
	n, err := hex.Decode(dst, src)
	if err != nil {
		log.Fatalf("%s expected to contain hex-encoded string\n", filepath)
	}

	*b = dst[:n]
}

type config struct {
	Bind            string
	ReadTimeout     int
	WriteTimeout    int
	DownloadTimeout int
	Concurrency     int
	MaxClients      int
	TTL             int

	MaxSrcDimension int

	Quality         int
	GZipCompression int

	Key  []byte
	Salt []byte

	Secret string
}

var conf = config{
	Bind:            ":8080",
	ReadTimeout:     10,
	WriteTimeout:    10,
	DownloadTimeout: 5,
	Concurrency:     runtime.NumCPU() * 2,
	TTL:             3600,
	MaxSrcDimension: 4096,
	Quality:         80,
	GZipCompression: 5,
}

func init() {
	keypath := flag.String("keypath", "", "path of the file with hex-encoded key")
	saltpath := flag.String("saltpath", "", "path of the file with hex-encoded salt")
	flag.Parse()

	strEnvConfig(&conf.Bind, "IMGPROXY_BIND")
	intEnvConfig(&conf.ReadTimeout, "IMGPROXY_READ_TIMEOUT")
	intEnvConfig(&conf.WriteTimeout, "IMGPROXY_WRITE_TIMEOUT")
	intEnvConfig(&conf.DownloadTimeout, "IMGPROXY_DOWNLOAD_TIMEOUT")
	intEnvConfig(&conf.Concurrency, "IMGPROXY_CONCURRENCY")
	intEnvConfig(&conf.MaxClients, "IMGPROXY_MAX_CLIENTS")

	intEnvConfig(&conf.TTL, "IMGPROXY_TTL")

	intEnvConfig(&conf.MaxSrcDimension, "IMGPROXY_MAX_SRC_DIMENSION")

	intEnvConfig(&conf.Quality, "IMGPROXY_QUALITY")
	intEnvConfig(&conf.GZipCompression, "IMGPROXY_GZIP_COMPRESSION")

	hexEnvConfig(&conf.Key, "IMGPROXY_KEY")
	hexEnvConfig(&conf.Salt, "IMGPROXY_SALT")

	hexFileConfig(&conf.Key, *keypath)
	hexFileConfig(&conf.Salt, *saltpath)

	strEnvConfig(&conf.Secret, "IMGPROXY_SECRET")

	if len(conf.Key) == 0 {
		log.Fatalln("Key is not defined")
	}
	if len(conf.Salt) == 0 {
		log.Fatalln("Salt is not defined")
	}

	if len(conf.Bind) == 0 {
		log.Fatalln("Bind address is not defined")
	}

	if conf.ReadTimeout <= 0 {
		log.Fatalf("Read timeout should be greater than 0, now - %d\n", conf.ReadTimeout)
	}

	if conf.WriteTimeout <= 0 {
		log.Fatalf("Write timeout should be greater than 0, now - %d\n", conf.WriteTimeout)
	}

	if conf.DownloadTimeout <= 0 {
		log.Fatalf("Download timeout should be greater than 0, now - %d\n", conf.DownloadTimeout)
	}

	if conf.Concurrency <= 0 {
		log.Fatalf("Concurrency should be greater than 0, now - %d\n", conf.Concurrency)
	}

	if conf.MaxClients <= 0 {
		conf.MaxClients = conf.Concurrency * 5
	}

	if conf.TTL <= 0 {
		log.Fatalf("TTL should be greater than 0, now - %d\n", conf.TTL)
	}

	if conf.MaxSrcDimension <= 0 {
		log.Fatalf("Max src dimension should be greater than 0, now - %d\n", conf.MaxSrcDimension)
	}

	if conf.Quality <= 0 {
		log.Fatalf("Quality should be greater than 0, now - %d\n", conf.Quality)
	} else if conf.Quality > 100 {
		log.Fatalf("Quality can't be greater than 100, now - %d\n", conf.Quality)
	}

	if conf.GZipCompression < 0 {
		log.Fatalf("GZip compression should be greater than or quual to 0, now - %d\n", conf.GZipCompression)
	} else if conf.GZipCompression > 9 {
		log.Fatalf("GZip compression can't be greater than 9, now - %d\n", conf.GZipCompression)
	}

	initVips()
}
