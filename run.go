/**
 * @Author : NewtSun
 * @Date : 2023/4/17 16:42
 * @Description :
 **/

package main

/*
$ curl "http://localhost:9999/api?key=Tom"
630
$ curl "http://localhost:9999/api?key=kkk"
kkk not exist
*/

import (
	"GoCache/controller"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *controller.Group {
	return controller.NewGroup("scores", 2<<10, controller.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

func startCacheServer(addr string, addrs []string, gee *controller.Group) {
	httpPoolPeersPicker := controller.NewHTTPPool(addr) // run local client
	httpPoolPeersPicker.Set(addrs...)
	gee.RegisterPeersPicker(httpPoolPeersPicker)
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], httpPoolPeersPicker))
}

func startAPIServer(apiAddr string, gee *controller.Group) {
	// 这里的想法是 根据不同的命名空间，在 /api 后面加上 group 的名字
	// 比如 9001 端口开放之后，允许访问不同的命名空间
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := gee.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())

		}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))

}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "Geecache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	gee := createGroup()
	//if api {
	//	go startAPIServer(apiAddr, gee)
	//}
	go startAPIServer(apiAddr, gee)
	startCacheServer(addrMap[port], addrs, gee)
}
