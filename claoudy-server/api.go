package main

import (
    "fmt"
    "io"
	//"os"
    "encoding/json"
    //"io/ioutil"
    "github.com/minio/minio-go"
	"github.com/dcasier/claoudy/metamodel"
	"bytes"
	"strings"
	"unicode"
	"golang.org/x/text/runes"
    "golang.org/x/text/transform"
    "golang.org/x/text/unicode/norm"
    "golang.org/x/crypto/bcrypt"
)

// TODO :
// - Garantir la consistance des datas (ACID) : si erreur, rollback sur les opérations précédentes

func AddGrantToEvent(event string, member string, grants string, keys metamodel.Keys) error {
	//todo : 
	// - Ajout d'ACL, changer l'api minio-go vers alz-v3/s3 ?
	// - Ou adapter minio-go
	evtDetails := new(metamodel.EventDetails)
	err := GetMeta(event, "details", evtDetails, keys)
    if err != nil {
        fmt.Println("- AddGrantToEvent")
        fmt.Println(err)
    }
	if evtDetails.Grants == nil {
		evtDetails.Grants = make(map[string]string)
	}
	evtDetails.Grants[member] = grants
	return PutMeta(event, "details", evtDetails, keys)
}

func BucketName(name string) string {

    t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
    s, _, _ := transform.String(t, name)
	fmt.Println(s)
	return strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.ToLower(s), "_", "--", -1), " ", "---", -1), "-.", "-----", -1), ".-", "----", -1)
}

func Client(keys metamodel.Keys) *minio.Client {
    //file, e := ioutil.ReadFile("c:/Go/work/bin/env/config.json")
    //if e != nil {
    //    fmt.Printf("File error: %v\n", e)
    //    os.Exit(1)
    //}
    //var acc metamodel.Access
    //json.Unmarshal(file, &acc)
	client, err := minio.New("127.0.0.1:9000", keys.AccessKey, keys.SecretKey, false)
    if err != nil {
        fmt.Println("- Client")
		fmt.Println(err)
    }
	return client
}

func DeleteActivity(activity string, keys metamodel.Keys) {
    meta := new(metamodel.Activities)
	err := GetMeta("activities", "activities", meta, keys)
    if err != nil {
        fmt.Println("- DeleteActivity - GetMeta")
        fmt.Println(err)
    }
	delete(meta.Sphere, activity)
	err = PutMeta("activities", "activities", meta, keys)
    if err != nil {
        fmt.Println("- DeleteActivity - PutMeta")
        fmt.Println(err)
    }
	fmt.Println(activity)
	err = DeleteMeta("activities", activity, keys)
    if err != nil {
        fmt.Println("- DeleteActivity - DeleteMeta")
        fmt.Println(err)
    }
}

func DeleteMeta(bucket string, object string, keys metamodel.Keys) error {
	return DeleteObject(bucket, object, keys)
}

func DeleteObject(bucketname string, name string, keys metamodel.Keys) error {
	err := Client(keys).RemoveObject(BucketName(bucketname), name)
	return err
}

func DelGrantToEvent(event string, member string, keys metamodel.Keys) error {
	evtDetails := new(metamodel.EventDetails)
	err := GetMeta(event, "details", evtDetails, keys)
    if err != nil {
        fmt.Println("- DelGrantToEvent")
        fmt.Println(err)
    }
	delete(evtDetails.Grants, member)
	return PutMeta(event, "details", evtDetails, keys)
}

func GetActivities(keys metamodel.Keys) (*metamodel.Activities, error) {
    meta := new(metamodel.Activities)
	err := GetMeta("activities", "activities", meta, keys)
	return meta, err
}

func GetActivity(activity string, keys metamodel.Keys) (metamodel.Sphere, error) {
    meta := new(metamodel.Activities)
	err := GetMeta("activities", "activities", meta, keys)
	return meta.Sphere[activity], err
}

func GetCache(keys metamodel.Keys) (*metamodel.Cache, error) {
	meta := new(metamodel.Cache)
	err := GetMeta("activities", "cache", meta, keys)
	return meta, err
}

func GetConfig(keys metamodel.Keys) (*metamodel.Config, error) {
	meta := new(metamodel.Config)
	err := GetMeta("activities", "config", meta, keys)
	return meta, err
}

func GetEvent(event string, keys metamodel.Keys) (*metamodel.EventDetails, error) {
	evtDetails := new(metamodel.EventDetails)
	err := GetMeta(event, "details", evtDetails, keys)
	return evtDetails, err
}

func GetMedia(event string, name string, keys metamodel.Keys) *minio.Object {
	return GetObject(event, name, keys)
}

/*
	Peuple l'interface "data" avec le contenu de bucket/object
*/
func GetMeta(bucket string, object string, data interface{}, keys metamodel.Keys) error {
	obj := GetObject(bucket, object, keys)
	return SetMeta(obj, data)
}

func getOrMakeBucket(bucket string, keys metamodel.Keys) error {
	found, err := Client(keys).BucketExists(BucketName(bucket))
    if err != nil {
        fmt.Println("- getOrMakeBucket")
        fmt.Println(err)
    }
	if found {
		return nil
	}else {
		return Client(keys).MakeBucket(BucketName(bucket), "us-east-1")
	}
}

func GetObject(bucket string, object string, keys metamodel.Keys) *minio.Object {
	obj, err := Client(keys).GetObject(BucketName(bucket), object)
    if err != nil {
        fmt.Println("- GetObject")
        fmt.Println(err)
    }
	return obj
}

func InterfaceToJson(data interface{}) []byte {
	b, err := json.Marshal(data)
    if err != nil {
        fmt.Println("- InterfaceToJson")
        fmt.Println(err)
    }
	return b
}

func Login(body io.Reader) ([]byte, error) {
	login := new(metamodel.Login)
	err   := SetMeta(body, login)
	if err != nil {
        fmt.Println("- Login - SetMeta")
        fmt.Println(err)
    }
	user  := new(metamodel.User)
	keys := new(metamodel.Keys)
	keys.AccessKey = "08ATNB835CYQ504UXM0N"
	keys.SecretKey = "kRcBHv5dyhtfVS3jlIKBQ97RWrrm0fBUSUsaxBAP"
	err   = GetMeta("activities", login.Login.UserName, user, *keys)
    if err != nil {
        fmt.Println("- Login - GetMeta")
        fmt.Println(err)
    }
	err = bcrypt.CompareHashAndPassword([]byte(user.Info.PasswdHash), []byte(login.Login.Password))
	if err != nil {
		var nilKeys metamodel.Keys
		return InterfaceToJson(nilKeys), err
    }
	return InterfaceToJson(user.Keys), nil
}

func Logout(token string) error {
	return nil
}

func ListMedia(activity string, event string, keys metamodel.Keys) []byte {
	// Create a done channel to control 'ListObjects' go routine.
	doneCh := make(chan struct{})

	// Indicate to our routine to exit cleanly upon return.
	defer close(doneCh)
	medias := new(metamodel.Medias)
	for object := range Client(keys).ListObjectsV2(BucketName(event), "", true, doneCh) {
		medias.Medias = append(medias.Medias, object)
	}
	return InterfaceToJson(medias)
}

func PostActivity(activity string, body io.Reader, keys metamodel.Keys) error {
	meta := new(metamodel.PostActivity)
	err := SetMeta(body, meta)
    if err != nil {
        fmt.Println("- PostActivity")
        fmt.Println(err)
    }
	if meta.AddEvent != (metamodel.Event{}) {
		event := new(metamodel.Event)
		event.Name = meta.AddEvent.Name
		err = PutEvent(activity, event, keys)
	}
	return err
}

func PostMedia(activity string, event string, name string, file io.Reader, keys metamodel.Keys) error {
	return PutObject(event, name, file, "image/jpeg", keys)
}

func PutActivity(activity string, keys metamodel.Keys) error {
	err := getOrMakeBucket("activities", keys)
    if err != nil {
        fmt.Println("- PutActivity")
        fmt.Println(err)
    }
	meta := new(metamodel.Activities)
	err = GetMeta("activities", "activities", meta, keys)
    if err != nil {
        fmt.Println("- PutActivity - GetMeta")
        fmt.Println(err)
    }
    if meta.Sphere == nil {
		meta.Sphere = make(map[string]metamodel.Sphere)
    }
	meta.Sphere[activity] = metamodel.Sphere{Allow: make([]string, 0), Events: make([]string,0)}
	return PutMeta("activities", "activities", meta, keys)
}

func PutCache(obj io.Reader, keys metamodel.Keys) error {
	meta := new(metamodel.Cache)
	err  := SetMeta(obj, meta)
	if err != nil {
        fmt.Println("- PutCache")
        fmt.Println(err)
    }
	return PutMeta("activities", "cache", meta, keys)
}

func PutEvent(activity string, event *metamodel.Event, keys metamodel.Keys) error {
	meta := new(metamodel.Activities)
	err := GetMeta("activities", "activities", meta, keys)
    if err != nil {
        fmt.Println("- PutEvent")
        fmt.Println(err)
    }
	events := append(meta.Sphere[activity].Events, event.Name)	
	meta.Sphere[activity] = metamodel.Sphere{Allow: meta.Sphere[activity].Allow, Events: events}
	err = PutMeta("activities", "activities", meta, keys)
	getOrMakeBucket(event.Name, keys)
	return err
}

func PutMedia(event string, name string, file io.Reader, keys metamodel.Keys) error {
	return PutObject(event, name, file, "image/jpeg", keys)
}

/*
	Stocke l'interface "data" dans bucket/object
*/
func PutMeta(bucket string, object string, data interface{}, keys metamodel.Keys) error {
	b := InterfaceToJson(data)
	return PutObject(bucket, object, bytes.NewReader(b), "application/octet-stream", keys)
}

func PutObject(bucketname string, name string, file io.Reader,fileType string, keys metamodel.Keys) error {
	_, err := Client(keys).PutObject(BucketName(bucketname), name, file, fileType)
	fmt.Println(keys)
	fmt.Println(err)
	return err
}

func Register(body io.Reader) ([]byte, error) {
	register := new(metamodel.Register)
	err   := SetMeta(body, register)
	if err != nil {
        fmt.Println("- Register - SetMeta")
        fmt.Println(err)
    }
	user := new(metamodel.User)
	user.Info.UserName = register.Register.UserName
	user.Info.Mail	   = register.Register.Mail
	bytePasswdHash, err := bcrypt.GenerateFromPassword([]byte(register.Register.Password), bcrypt.DefaultCost)	
	if err != nil {
        fmt.Println("- Register - GenerateFromPassword")
        fmt.Println(err)
    }	
	user.Info.PasswdHash = string(bytePasswdHash)
	
	user.Keys.AccessKey = "08ATNB835CYQ504UXM0N"
	user.Keys.SecretKey = "kRcBHv5dyhtfVS3jlIKBQ97RWrrm0fBUSUsaxBAP"
	
	err = PutMeta("activities", register.Register.UserName, user, user.Keys)
    if err != nil {
        fmt.Println("- SetPassword - PutMeta")
        fmt.Println(err)
    }		
	return InterfaceToJson(user.Keys), err
}

func SetMeta(obj io.Reader, data interface{}) error {
	return json.NewDecoder(obj).Decode(data)
}

func SetPassword(body io.Reader, keys metamodel.Keys) ([]byte, error) {
	setpassword := new(metamodel.SetPassword)
	err   := SetMeta(body, setpassword)
    _, err = Login(body)
    if err != nil {
		var nilToken metamodel.Keys
        return InterfaceToJson(nilToken), err
    }
	user := new(metamodel.User)
	err = GetMeta("activities", setpassword.SetPassword.UserName, user, keys)
    if err != nil {
        fmt.Println("- SetPassword - GetMeta")
        fmt.Println(err)
    }	
	bytePasswdHash, err := bcrypt.GenerateFromPassword([]byte(user.Info.PasswdHash), bcrypt.DefaultCost)
	user.Info.PasswdHash = string(bytePasswdHash)
	err = PutMeta("activities", setpassword.SetPassword.UserName, user, keys)
    if err != nil {
        fmt.Println("- SetPassword - PutMeta")
        fmt.Println(err)
    }		
	return InterfaceToJson(user.Keys), err
}
