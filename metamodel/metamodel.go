package metamodel

import   "github.com/minio/minio-go"

type Access struct {
	Endpoint	string
    AccessKey	string
    SecretKey	string
	UseSSL		bool
}

// Localisation : bucket activities
// Type : objet
// {"Sphere":{"musique":"musique"}, {"sport":"sport"}}
type Activities struct {
	Sphere map[string]Sphere	`json:"sphere"`
}

// Localisation : bucket activities
// Type : objet
// {"Evts":{"concert du 25":{"Desc":"Description du match"}}}
type Activity struct {
	Evts []Event
}

// Activité de l'utilisateur
// { sphere: {spherename: [lastevent1, lastevent2, ...]}}
//
type Cache struct {
	Cache 			Spheres		`json:"cache"`
}

// Paramètres de l'interface, complémentaire avec le cache
// { presentation: False, }
//
type Config struct {
	Config 			map[string][]bool		`json:"config"`
}

type Event struct {
	Name string					`json:"name"`
	Desc string					`json:"desc"`
}

// Localisation : bucket <Event name>
// {"Desc":"","Medias":null,"Members":null,"Grants":{"user1":"READ","user2":"FULLCONTROL","user3":"READ","user4":"FULLCONTROL","user6":"READ"}}
type EventDetails struct {
	Desc string
	Medias map[string]string
	Members []string
	Grants map[string]string
}

type Medias struct {
	Medias	[]minio.ObjectInfo	`json:"medias"`
}

type PostActivity struct {
	AddEvent 	Event		`json:"addEvent"`
	DeleteEvent Event 		`json:"deleteEvent"`
}

type Sphere struct {
	Allow	[]string			`json:"allow"`
	Events	[]string			`json:"events"`
}

type Spheres struct {
	LastSpheres		[]string				`json:"lastspheres"`
	Spheres			map[string][]string		`json:"spheres"`
}

type User struct {
	Info  UserInfo			`json:"info"`
	Keys Keys				`json:"keys"`
}

type Identification struct {
	UserName	string		`json:"username"`
	Password	string		`json:"password"`
	NewPassword string		`json:"newpassword"`
	Mail		string		`json:"mail"`
}

type Login struct {
	Login	Identification	`json:"login"`
}

type Register struct {
	Register	Identification	`json:"register"`
}

type SetPassword struct {
	SetPassword	Identification	`json:"setpassword"`
}

type UserInfo struct {
	UserName	string		`json:"username"`
	Mail		string		`json:"mail"`
	PasswdHash	string		`json:"passwordhash"`
}

type Keys struct {
	AccessKey	string		`json:"accesskey"`
	SecretKey	string		`json:"secretkey"`
}