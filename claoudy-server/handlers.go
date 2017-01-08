package main

import(
	"net/http"
    "github.com/gorilla/mux"
	"github.com/dcasier/claoudia/metamodel"
	"encoding/base64"
	//"io/ioutil"
	"fmt"
	"io"
	"strconv"
        "net"
        "time"
		
)

func DeleteActivityHandler(w http.ResponseWriter, r *http.Request) {
	GetAllowOrigin(w, r)
    // Stop here if its Preflighted OPTIONS request
    if r.Method == "OPTIONS" {
        return
    }
	keys := GetKeys(r)
	activity := mux.Vars(r)["activity"]
	DeleteActivity(activity, *keys)
	act, err := GetActivities(*keys)
    if err != nil {
        fmt.Println("- DeleteActivityHandler")
        fmt.Println(err)
    }
	w.Write(InterfaceToJson(act))
}

func GetAllowOrigin(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	fmt.Println(origin)
    if origin == "http://192.168.1.26:3000" || origin == "http://192.168.1.11:3000" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
        w.Header().Set("Access-Control-Allow-Headers",
            "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Allow-Origin, Private-Token, responsetype, X-Auth-Token, AccessKey, SecretKey")
		
		//if r.Header.Get("AccessKey") != "" && r.Header.Get("SecretKey") != "" {
		//	w.Header().Set("AccessKey", r.Header.Get("AccessKey"))
		//	w.Header().Set("SecretKey", r.Header.Get("SecretKey"))
		//}
    }
}

func GetKeys(r *http.Request) *metamodel.Keys {
	keys := new(metamodel.Keys)
	keys.AccessKey = r.Header.Get("AccessKey")
	keys.SecretKey = r.Header.Get("SecretKey")
	return keys
}

func GetActivitiesHandler(w http.ResponseWriter, r *http.Request) {
	GetAllowOrigin(w, r)
    // Stop here if its Preflighted OPTIONS request
    if r.Method == "OPTIONS" {
        return
    }
	keys := GetKeys(r)
	acts, err := GetActivities(*keys)
    if err != nil {
        fmt.Println("- GetActivitiesHandler")
        fmt.Println(err)
    }
	w.Write(InterfaceToJson(acts))
}

func GetActivityHandler(w http.ResponseWriter, r *http.Request) {
	GetAllowOrigin(w, r)
    // Stop here if its Preflighted OPTIONS request
    if r.Method == "OPTIONS" {
        return
    }
	keys := GetKeys(r)
	act, err := GetActivity(mux.Vars(r)["activity"], *keys)
    if err != nil {
        fmt.Println("- GetActivityHandler")
        fmt.Println(err)
    }
	w.Write(InterfaceToJson(act))
}

func GetCacheHandler(w http.ResponseWriter, r *http.Request) {
	GetAllowOrigin(w, r)
    // Stop here if its Preflighted OPTIONS request
    if r.Method == "OPTIONS" {
        return
    }
	keys := GetKeys(r)
	cache, err := GetCache(*keys)
    if err != nil {
        fmt.Println("- GetCacheHandler")
        fmt.Println(err)
    }
	w.Write(InterfaceToJson(cache))
}

func GetConfigHandler(w http.ResponseWriter, r *http.Request) {
	GetAllowOrigin(w, r)
    // Stop here if its Preflighted OPTIONS request
    if r.Method == "OPTIONS" {
        return
    }
	keys := GetKeys(r)
	config, err := GetConfig(*keys)
    if err != nil {
        fmt.Println("- GetConfigHandler")
        fmt.Println(err)
    }
	w.Write(InterfaceToJson(config))
}

func GetEventHandler(w http.ResponseWriter, r *http.Request) {
	GetAllowOrigin(w, r)
    // Stop here if its Preflighted OPTIONS request
    if r.Method == "OPTIONS" {
        return
    }
	keys := GetKeys(r)
	vars := mux.Vars(r)
	evt, err := GetEvent(vars["event"], *keys)
    if err != nil {
        fmt.Println("- GetEventHandler")
        fmt.Println(err)
    }
	w.Write(InterfaceToJson(evt))
}

func GetMediaHandler(w http.ResponseWriter, r *http.Request) {
	GetAllowOrigin(w, r)
    // Stop here if its Preflighted OPTIONS request
    if r.Method == "OPTIONS" {
        return
    }

	//keys := &metamodel.Keys{AccessKey:  "08ATNB835CYQ504UXM0N", SecretKey: "kRcBHv5dyhtfVS3jlIKBQ97RWrrm0fBUSUsaxBAP"}
	keys := GetKeys(r)
	fmt.Println(keys)
	vars := mux.Vars(r)
	object := GetObject(vars["event"], vars["media"], *keys)
	//b, err := ioutil.ReadAll(object)
	//fmt.Println(b[:32])
    //if err != nil {
    //    fmt.Println("- GetMediaHandler")
    //    fmt.Println(err)
    //}	
	//w.Write(b)
	w.Header().Set("Content-Disposition", "attachment; filename="+vars["media"])
    stat, _ := object.Stat()
	fmt.Println(stat)
	w.Header().Set("Content-Type", stat.ContentType)
    w.Header().Set("Content-Length", strconv.FormatInt(stat.Size, 10))
	w.Header().Set("responseType", "Blob")
	io.Copy(w, object)
}

func Gz(w http.ResponseWriter, r *http.Request) {
        url := "http://upload.wikimedia.org/wikipedia/en/b/bc/Wiki.png"

        timeout := time.Duration(5) * time.Second
        transport := &http.Transport{
                ResponseHeaderTimeout: timeout,
                Dial: func(network, addr string) (net.Conn, error) {
                        return net.DialTimeout(network, addr, timeout)
                },
                DisableKeepAlives: true,
        }
        client := &http.Client{
                Transport: transport,
        }
        resp, err := client.Get(url)
        if err != nil {
                fmt.Println(err)
        }
        defer resp.Body.Close()

        //copy the relevant headers. If you want to preserve the downloaded file name, extract it with go's url parser.
        w.Header().Set("Content-Disposition", "attachment; filename=ThumbnailHPIM3396.JPG")
        w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
        w.Header().Set("Content-Length", r.Header.Get("Content-Length"))

        //stream the body to the client without fully loading it into memory
		//w.Write(resp.Body)
        io.Copy(w, resp.Body)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	GetAllowOrigin(w, r)
    // Stop here if its Preflighted OPTIONS request
    if r.Method == "OPTIONS" {
        return
    }
	token, err := Login(r.Body)
	if err != nil {
        fmt.Println("- LoginHandler")
        fmt.Println(err)
    }
	w.Write(token)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	GetAllowOrigin(w, r)
    // Stop here if its Preflighted OPTIONS request
    if r.Method == "OPTIONS" {
        return
    }
	//vars := mux.Vars(r)
	err := Logout("token")
	if err != nil {
        fmt.Println("- LogoutHandler")
        fmt.Println(err)
    }
}

func ListMediaHandler(w http.ResponseWriter, r *http.Request) {
	GetAllowOrigin(w, r)
    // Stop here if its Preflighted OPTIONS request
    if r.Method == "OPTIONS" {
        return
    }
	keys := GetKeys(r)
	vars := mux.Vars(r)
	w.Write(ListMedia(vars["activity"], vars["event"], *keys))
}

func PostActivityHandler(w http.ResponseWriter, r *http.Request) {
	GetAllowOrigin(w, r)
    // Stop here if its Preflighted OPTIONS request
    if r.Method == "OPTIONS" {
        return
    }
	keys := GetKeys(r)
	activity := mux.Vars(r)["activity"]
	err := PostActivity(activity, r.Body, *keys)
	if err != nil {
        fmt.Println("- PostActivityHandler - PostActivity")
        fmt.Println(err)
    }
	act, err := GetActivity(activity, *keys)
    if err != nil {
        fmt.Println("- PostActivityHandler - GetActivity")
        fmt.Println(err)
    }
	w.Write(InterfaceToJson(act))
}

func PostEventHandler(w http.ResponseWriter, r *http.Request) {
	GetAllowOrigin(w, r)
    // Stop here if its Preflighted OPTIONS request
    if r.Method == "OPTIONS" {
        return
    }
	keys := GetKeys(r)
	//formFile := r.MultipartForm
	fmt.Println(r.Header)
	vars := mux.Vars(r)
	mr, err := r.MultipartReader()
	if err != nil {
        fmt.Println("- PostEventHandler - r.MultipartReader")
        fmt.Println(err)
    }
	for {
		np, err := mr.NextPart()
		if err != nil {
			if err ==  io.EOF {
				//No next part
				break
			}
	        fmt.Println("- PostEventHandler - NextPart")
			fmt.Println(err)
		}
		err = PostMedia(vars["activity"], vars["event"], np.FileName(), base64.NewDecoder(base64.StdEncoding, np), *keys)
	    if err != nil {
	        fmt.Println("- PostEventHandler - PostMedia")
			fmt.Println(err)
		}
	}
}

func PostMediaHandler(w http.ResponseWriter, r *http.Request) {
	GetAllowOrigin(w, r)
    // Stop here if its Preflighted OPTIONS request
    if r.Method == "OPTIONS" {
        return
    }
	keys := GetKeys(r)
	vars := mux.Vars(r)
	//dec  := base64.NewDecoder(base64.StdEncoding, r.Body)
	err := PostMedia(vars["activity"], vars["event"], vars["media"], r.Body, *keys)
    if err != nil {
        fmt.Println("- PostMediaHandler")
        fmt.Println(err)
    }
}

func PutActivityHandler(w http.ResponseWriter, r *http.Request) {
	GetAllowOrigin(w, r)
    // Stop here if its Preflighted OPTIONS request
    if r.Method == "OPTIONS" {
        return
    }
	keys := GetKeys(r)
	activity := mux.Vars(r)["activity"]
	err := PutActivity(activity, *keys)
	if err != nil {
        fmt.Println("- PutActivityHandler - PutActivity")
        fmt.Println(err)
    }
	act, err := GetActivity(activity, *keys)
    if err != nil {
        fmt.Println("- PutActivityHandler - GetActivity")
        fmt.Println(err)
    }
	w.Write(InterfaceToJson(act))
}

func PutCacheHandler(w http.ResponseWriter, r *http.Request) {
	GetAllowOrigin(w, r)
    // Stop here if its Preflighted OPTIONS request
    if r.Method == "OPTIONS" {
        return
    }
	keys := GetKeys(r)
	err := PutCache(r.Body, *keys)
	cache, err := GetCache(*keys)
    if err != nil {
        fmt.Println("- PutCacheHandler")
        fmt.Println(err)
    }
	w.Write(InterfaceToJson(cache))
}

func PutEventHandler(w http.ResponseWriter, r *http.Request) {
	GetAllowOrigin(w, r)
    // Stop here if its Preflighted OPTIONS request
    if r.Method == "OPTIONS" {
        return
    }
	keys := GetKeys(r)
	evt := new(metamodel.Event)
	vars := mux.Vars(r)
	evt.Name = vars["event"]
	err := PutEvent(vars["activity"], evt, *keys)
    if err != nil {
        fmt.Println("- PutEventHandler - PutEvent")
        fmt.Println(err)
    }
	event, err := GetEvent(vars["event"], *keys)
    if err != nil {
        fmt.Println("- PutEventHandler - GetEvent")
        fmt.Println(err)
    }
	w.Write(InterfaceToJson(event))
}

func PutMediaHandler(w http.ResponseWriter, r *http.Request) {
	GetAllowOrigin(w, r)
    // Stop here if its Preflighted OPTIONS request
    if r.Method == "OPTIONS" {
        return
    }
	keys := GetKeys(r)
	vars := mux.Vars(r)
	err := PutMedia(vars["event"], vars["media"], r.Body, *keys)
    if err != nil {
        fmt.Println("- PutMediaHandler")
        fmt.Println(err)
    }
}

func SetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	GetAllowOrigin(w, r)
    // Stop here if its Preflighted OPTIONS request
    if r.Method == "OPTIONS" {
        return
    }
	keys := GetKeys(r)
	token, err := SetPassword(r.Body, *keys)
	if err != nil {
        fmt.Println("- SetPasswordHandler")
        fmt.Println(err)
    }
	w.Write(token)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	GetAllowOrigin(w, r)
    // Stop here if its Preflighted OPTIONS request
    if r.Method == "OPTIONS" {
        return
    }
	token, err := Register(r.Body)
	if err != nil {
        fmt.Println("- LoginHandler")
        fmt.Println(err)
    }
	w.Write(token)
}
