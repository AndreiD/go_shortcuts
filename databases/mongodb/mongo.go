package database

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	mongoURL       = "localhost:27017"
	DbName         = "xxxxx"
	CollectionName = "xxxxxxxxxxxx"
)

// Session is the mongodb session
type Session struct {
	session *mgo.Session
}

// NewSession creates a new session to work with. Connection url should be moved to  some config files
func NewSession() (*Session, error) {
	//var err error
	session, err := mgo.Dial(mongoURL)
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)
	return &Session{session}, err
}

// Copy copies an existing session
func (s *Session) Copy() *mgo.Session {
	return s.session.Copy()
}

// Close ends the session
func (s *Session) Close() {
	if s.session != nil {
		s.session.Close()
	}
}

// DeleteCollection delets a collection (for testing)
func (s *Session) DeleteCollection() error {
	if s.session != nil {
		_, err := s.session.DB(DbName).C(CollectionName).RemoveAll(bson.M{})
		if err != nil {
			return err
		}
	}
	return nil
}
