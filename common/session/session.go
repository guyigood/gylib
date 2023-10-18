package sessions

import (
	"encoding/json"
	"errors"
	"github.com/garyburd/redigo/redis"
	"github.com/guyigood/gylib/common"
	"github.com/guyigood/gylib/common/redispack"
	"github.com/satori/go.uuid"
	"net/http"
	"sync"
	"time"
)

// Errors
var ErrNoConnection = errors.New("connection to redispack has not been established")

var Redis_Pool *redis.Pool
var isConnected bool
var Gsession *Session
var opts SessionOptions

type SessionOptions struct {
	SessionKey string
	Timeout    string
}

func session_init() {
	Redis_Pool = redispack.Get_redis_pool()
	isConnected = true
	data := common.Getini("conf/app.ini", "session", map[string]string{"sessionkey": "", "timeout": ""})
	opts.SessionKey = data["sessionkey"]
	opts.Timeout = data["timeout"]

}

type Session struct {
	sync.Mutex
	ID     string                 `json:"sessionId"`
	Values map[string]interface{} `json:"values"`
}

func Start_session(w http.ResponseWriter, r *http.Request) error {
	if !isConnected {
		session_init()
	}
	var err error
	Gsession, err = Open(w, r)
	//fmt.Println("gsession", Gsession, err)
	if err != nil {
		return err
	}
	return nil
}

// Get will return the value from an existing session
func (s *Session) Get(param string) interface{} {
	client := Redis_Pool.Get()
	defer client.Close()
	if !isConnected {
		return nil
	}
	//client.Do("GET",param)
	s.Lock()
	defer s.Unlock()

	return s.Values[param]
}

// Set a value on a session and store it to redispack
func (s *Session) Set(param string, val interface{}) error {

	s.Lock()
	defer s.Unlock()
	client := Redis_Pool.Get()
	defer client.Close()
	if !isConnected {
		return ErrNoConnection
	}

	if s.Values == nil {
		s.Values = make(map[string]interface{})
	}

	s.Values[param] = val

	raw, _ := json.Marshal(s)
	client.Do("SETEX", s.ID, opts.Timeout, raw)
	//dostr,err:=client.Do("SETEX", s.ID, opts.Timeout,raw)
	//fmt.Println("set redis", s.ID,raw,dostr,err)
	return nil
}

// Clear an existing session
func (s *Session) Clear() {

	s.Lock()
	defer s.Unlock()

	s.Values = make(map[string]interface{})

	raw, _ := json.Marshal(s)
	client := Redis_Pool.Get()
	defer client.Close()
	client.Do("SETEX", s.ID, opts.Timeout, raw)
}

// Connect to the Redis server we'll be using
// for session storage.
func Connect(options SessionOptions) error {
	opts = options
	return nil
}

// End closes the connection with redispack
func End() error {
	if !isConnected {
		return ErrNoConnection
	}
	err := Redis_Pool.Close()
	if err != nil {
		return err
	}

	isConnected = false

	return nil
}

// HasSession lets you check an incoming request
// for an existing session.
func HasSession(r *http.Request) (bool, error) {

	if !isConnected {
		return false, ErrNoConnection
	}

	sess, err := r.Cookie(opts.SessionKey)
	if err != nil {
		return false, err
	}

	if sess != nil {
		return true, nil
	}

	return false, nil
}

// Open will either get a session from an existing ID
// or if the cookie cannot be found a new session
// will be returned
func Open(w http.ResponseWriter, r *http.Request) (*Session, error) {
	var session Session
	client := Redis_Pool.Get()
	defer client.Close()
	if !isConnected {
		return &session, ErrNoConnection
	}
	sess, err := r.Cookie(opts.SessionKey)
	if err == nil {
		// Cookie was found; let's look it up
		reply, err := client.Do("GET", sess.Value)
		if reply != nil {
			raw := reply.([]byte)
			err = json.Unmarshal(raw, &session)
			if err == nil {
				//fmt.Println("getjson",raw)
				return &session, err
			}

		}
		//return &session, err
	}
	// No session found. Let's make one
	session.ID = generateSessionID()
	raw, _ := json.Marshal(session)
	client.Do("SETEX", session.ID, opts.Timeout, raw)
	cookie := &http.Cookie{
		Name:  opts.SessionKey,
		Value: session.ID,
	}
	cookie.Expires = time.Now().AddDate(0, 0, 1)
	cookie.Path = "/"
	http.SetCookie(w, cookie)
	return &session, nil
}

func generateSessionID() string {
	uuid := uuid.NewV4()
	u_str := "ses_" + uuid.String() + common.GetRangStr(999999)
	return u_str
}

//func generateSessionID() string {
//
//	raw := make([]byte, 30)
//
//	_, err := rand.Read(raw)
//	if err != nil {
//		return generateSessionID()
//	}
//
//	return hex.EncodeToString(raw)
//}
