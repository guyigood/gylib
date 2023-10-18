package websession

import (
	"encoding/json"
	"errors"
	"github.com/guyigood/gylib/common"
	"github.com/guyigood/gylib/common/datatype"
	"github.com/guyigood/gylib/common/redisclient"
	"github.com/satori/go.uuid"
	"net/http"
	"sync"
	"time"
)

// Errors
var ErrNoConnection = errors.New("connection to redispack has not been established")

var Redis_Pool *redisclient.RedisClient
var isConnected bool
var Gsession *Session
var opts SessionOptions

type SessionOptions struct {
	SessionKey string
	Timeout    int64
}

func session_init() {
	Redis_Pool = redisclient.NewRedisCient()
	isConnected = true
	data := common.Getini("conf/app.ini", "session", map[string]string{"sessionkey": "", "timeout": ""})
	opts.SessionKey = data["sessionkey"]
	opts.Timeout = datatype.Str2Int64(data["timeout"])

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
	s.Lock()
	defer s.Unlock()
	return s.Values[param]
}

// Set a value on a session and store it to redispack
func (s *Session) Set(param string, val interface{}) error {
	s.Lock()
	defer s.Unlock()
	if s.Values == nil {
		s.Values = make(map[string]interface{})
	}
	s.Values[param] = val
	raw, _ := json.Marshal(s)
	Redis_Pool.SetKey(s.ID).SetExValue(raw, opts.Timeout)
	return nil
}

// Clear an existing session
func (s *Session) Clear() {
	s.Lock()
	defer s.Unlock()
	s.Values = make(map[string]interface{})
	Redis_Pool.SetKey(s.ID).DelKey()
}

// Connect to the Redis server we'll be using
// for session storage.
func Connect(options SessionOptions) error {
	opts = options
	return nil
}

// Open will either get a session from an existing ID
// or if the cookie cannot be found a new session
// will be returned
func Open(w http.ResponseWriter, r *http.Request) (*Session, error) {
	var session Session

	sess, err := r.Cookie(opts.SessionKey)
	if err == nil {
		// Cookie was found; let's look it up
		reply := Redis_Pool.SetKey(sess.Value).GetValue()
		if reply != nil {
			raw := datatype.Type2str(reply)
			err = json.Unmarshal([]byte(raw), &session)
			if err == nil {
				return &session, err
			}

		}
		//return &session, err
	}
	// No session found. Let's make one
	session.ID = generateSessionID()
	raw, _ := json.Marshal(session)
	Redis_Pool.SetKey(session.ID).SetExValue(raw, opts.Timeout)
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
