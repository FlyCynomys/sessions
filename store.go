package sessions

type Store interface {
	Set(key, value []byte) error //set session value
	Get(key []byte) []byte       //get session value
	Delete(key []byte) error     //delete session value
	SessionID() []byte           //back current sessionID
	Flush() error                //delete all data
}
