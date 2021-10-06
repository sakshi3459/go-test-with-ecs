package main

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

const intervalInDays = 10

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

	/* 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	   		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	   	})

	   	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
	   		fmt.Fprintf(w, "Hi")
	   	})

	   	log.Fatal(http.ListenAndServe(":8081", nil)) */

	startTs := time.Now().Unix()
	
	uploadKind(Pool)
	uploadKind(Storage)
	uploadKind(VolumePerf)
	uploadKind(MpsBusyRate)
	uploadKind(MpOwnerBusyRate)
	uploadKind(Cache)
	
	endTs := time.Now().Unix()
	fmt.Println("Total time: ", (endTs - startTs))

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
		fmt.Println("keys != 0", len(keys))
	}
}
