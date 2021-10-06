package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	"github.com/go-redis/redis/v7"
	"github.com/panjf2000/ants/v2"
)

const (
	intervalInDays             = 10
	REDIS_CONNECTION_POOL_SIZE = 150
	REDIS_IDLE_CONNECTION      = 80
	WORKER_THREAD_SIZE         = 100
)

type Kind int

const (
	Ldev Kind = iota
	Pool
	DkcAlert
	Ctl1Alert
	Ctl2Alert
	UnknownAlert
	PfStorage
	Storage
	SystemInfo
	MpsBusyRate
	MpOwnerBusyRate
	VolumePerf
	Cache
	Unknown
)

func (k Kind) ToString() string {
	switch k {
	case Ldev:
		return "ldev"
	case Pool:
		return "pool"
	case DkcAlert:
		return "dkcalert"
	case Ctl1Alert:
		return "ctl1alert"
	case Ctl2Alert:
		return "ctl2alert"
	case UnknownAlert:
		return "unknownalert"
	case PfStorage:
		return "pfstorage"
	case Storage:
		return "storage"
	case SystemInfo:
		return "systeminfo"
	case MpsBusyRate:
		return "mpsbusyrate"
	case MpOwnerBusyRate:
		return "mpownerbusyrate"
	case VolumePerf:
		return "volumeperf"
	case Cache:
		return "cache"
	}
	return "Unknown"
}

func main() {
	fmt.Println("begin schedulehandler")
	defer fmt.Println("end schedulehandler")

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	//registers the channel
	signal.Notify(sigs, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println("Caught SIGTERM, shutting down")
		// Finish any outstanding requests, then...
		done <- true
	}()

	fmt.Println("Starting application")
	// Main logic goes here

	startTs := time.Now().Unix()

	//uploadKind(Pool)
	uploadKind(Storage)
	/* uploadKind(VolumePerf)
	uploadKind(MpsBusyRate)
	uploadKind(MpOwnerBusyRate)
	uploadKind(Cache) */

	endTs := time.Now().Unix()
	fmt.Println("Total time: ", (endTs - startTs))

	fmt.Println("exiting")

	/* 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	   		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	   	})

	   	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
	   		fmt.Fprintf(w, "Hi")
	   	})

	   	log.Fatal(http.ListenAndServe(":8081", nil)) */

}

func uploadKind(kind Kind) {
	fmt.Println("begin uploadKind")
	defer fmt.Println("end uploadKind")

	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:     []string{"clustercfg.comjct12rnkm1sy.tygu6t.usw2.cache.amazonaws.com:6379"},
		Password:  "",                                    // no password set
		TLSConfig: &tls.Config{InsecureSkipVerify: true}, // TLS required when TransitEncryptionEnabled: true
	})

	var keys []string
	for i := 0; i < intervalInDays; i++ {
		day := time.Now().UTC().AddDate(0, 0, -1*i).Format("2006/01/02")
		pattern := fmt.Sprintf("%s/*/%s/*", kind.ToString(), day)
		val, err := client.Keys(pattern).Result()
		if err != nil {
			fmt.Println("Redis Err: ", err)
			continue
		}
		if len(val) != 0 {
			keys = append(keys, val...)
		}

	}
	fmt.Println("Total keys for kind - is -", kind, len(keys))
	if len(keys) != 0 {
		UploadS3(keys)
	}
}

var RedisClient = func() redis.Cmdable {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        []string{"clustercfg.comjct12rnkm1sy.tygu6t.usw2.cache.amazonaws.com:6379"},
		Password:     "",                                    // no password set
		TLSConfig:    &tls.Config{InsecureSkipVerify: true}, // TLS required when TransitEncryptionEnabled: true
		PoolSize:     REDIS_CONNECTION_POOL_SIZE,
		MinIdleConns: REDIS_IDLE_CONNECTION,
	})
	return client
}

type service struct {
	uploader s3manageriface.UploaderAPI
}

func NewService() (service, error) {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("Session Err:", err)
		return service{}, err
	}

	uploader := s3manager.NewUploader(sess)

	return service{uploader: uploader}, nil
}

func (svc service) upload(redisUpload RedisToS3Upload) {
	client := redisUpload.client
	key := redisUpload.key
	val, err := client.Get(key).Result()
	if err != nil {
		fmt.Println("Err: ", err)
		return
	}

	uploader := redisUpload.uploader
	reader := strings.NewReader(val)

	key += ".csv"

	fmt.Println("Key before upload: ", &key)
	fmt.Println("Body before upload: ", reader)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("olympus-metrics-archive-dev"),
		Key:    &key,
		Body:   reader,
	})

	if err != nil {
		fmt.Println("Err in upload: ", err)
		return
	}

}

func UploadS3(keys []string) {
	client := RedisClient()
	svc, err := NewService()
	if err != nil {
		fmt.Println("Err: ", err)
		return
	}
	//uploader := svc.uploader

	// Use the common pool.
	var wg sync.WaitGroup

	// Use the pool with a function,
	p, _ := ants.NewPoolWithFunc(WORKER_THREAD_SIZE, func(redisToS3Upload interface{}) {
		svc.upload(redisToS3Upload.(RedisToS3Upload))
		wg.Done()
	})

	defer p.Release()
	// Submit tasks one by one.
	for _, key := range keys {
		wg.Add(1)
		_ = p.Invoke(RedisToS3Upload{client: client, key: key, uploader: svc.uploader})
	}
	wg.Wait()

	_, ok := client.(*redis.ClusterClient)

	if ok {
		err := client.(*redis.ClusterClient).Close()
		if err != nil {
			fmt.Println("redis Close error:", err)
		}
	}

	fmt.Println("finish all tasks:", keys[0], keys[len(keys)-1])
}

type RedisToS3Upload struct {
	client   redis.Cmdable
	uploader s3manageriface.UploaderAPI
	key      string
}
