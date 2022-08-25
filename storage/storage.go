package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/openshift/osin"
)

type Storage struct {
	db *sql.DB
}

// New returns a new postgres storage instance.
func New(db *sql.DB) *Storage {
	return &Storage{db}
}

// Clone the storage if needed. For example, using mgo, you can clone the session with session.Clone
// to avoid concurrent access problems.
// This is to avoid cloning the connection at each method access.
// Can return itself if not a problem.
func (s *Storage) Clone() osin.Storage {
	return s
}

// Close the resources the Storage potentially holds (using Clone for example)
func (s *Storage) Close() {
}

// GetClient loads the client by id
func (s *Storage) GetClient(id string) (osin.Client, error) {
	var c osin.DefaultClient

	row := s.db.QueryRow("SELECT id, secret, redirect_uri, extra FROM client WHERE id=$1", id)

	err := row.Scan(&c.Id, &c.Secret, &c.RedirectUri)
	// TODO: check if no rows and return appropriate error
	if err != nil {
		return &osin.DefaultClient{}, errors.New(err.Error())
	}
	// TODO: User data isn't populated, check where it's needed
	return &c, nil
}

// UpdateClient updates the client (identified by its id) and replaces the values with the values of client.
func (s *Storage) UpdateClient(c osin.Client) error {
	// TODO: updating won't be supported immediately
	return nil
}

// CreateClient stores the client in the database and returns an error, if something went wrong.
func (s *Storage) CreateClient(c osin.Client) error {
	// TODO: User data isn't populated, check where it's needed
	_, err := s.db.Exec("INSERT INTO client (id, secret, redirect_uri) VALUES ($1, $2, $3)",
		c.GetId(), c.GetSecret(), c.GetRedirectUri())

	return err
}

// RemoveClient removes a client (identified by id) from the database. Returns an error if something went wrong.
func (s *Storage) RemoveClient(id string) (err error) {
	if _, err = s.db.Exec("DELETE FROM client WHERE id=$1", id); err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// SaveAuthorize saves authorize data.
func (s *Storage) SaveAuthorize(data *osin.AuthorizeData) (err error) {
	if _, err = s.db.Exec(
		"INSERT INTO authorize (client, code, expires_in, scope, redirect_uri, state, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		data.Client.GetId(),
		data.Code,
		data.ExpiresIn,
		data.Scope,
		data.RedirectUri,
		data.State,
		data.CreatedAt,
	); err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// LoadAuthorize looks up AuthorizeData by a code.
// Client information MUST be loaded together.
// Optionally can return error if expired.
func (s *Storage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	var data osin.AuthorizeData
	var cid string
	if err := s.db.QueryRow("SELECT client, code, expires_in, scope, redirect_uri, state, created_at FROM authorize WHERE code=$1 LIMIT 1", code).Scan(&cid, &data.Code, &data.ExpiresIn, &data.Scope, &data.RedirectUri, &data.State, &data.CreatedAt); err == sql.ErrNoRows {
		return nil, errors.New("not found")
	} else if err != nil {
		return nil, errors.New(err.Error())
	}

	c, err := s.GetClient(cid)
	if err != nil {
		return nil, err
	}

	if data.ExpireAt().Before(time.Now()) {
		return nil, errors.New(fmt.Sprintf("Token expired at %s.", data.ExpireAt().String()))
	}

	data.Client = c
	return &data, nil
}

// RemoveAuthorize revokes or deletes the authorization code.
func (s *Storage) RemoveAuthorize(code string) (err error) {
	if _, err = s.db.Exec("DELETE FROM authorize WHERE code=$1", code); err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// SaveAccess writes AccessData.
// If RefreshToken is not blank, it must save in a way that can be loaded using LoadRefresh.
func (s *Storage) SaveAccess(data *osin.AccessData) (err error) {

	return nil
}

// LoadAccess retrieves access data by token. Client information MUST be loaded together.
// AuthorizeData and AccessData DON'T NEED to be loaded if not easily available.
// Optionally can return error if expired.
func (s *Storage) LoadAccess(code string) (*osin.AccessData, error) {
	return nil, nil
}

// RemoveAccess revokes or deletes an AccessData.
func (s *Storage) RemoveAccess(code string) (err error) {
	return nil
}

// LoadRefresh retrieves refresh AccessData. Client information MUST be loaded together.
// AuthorizeData and AccessData DON'T NEED to be loaded if not easily available.
// Optionally can return error if expired.
func (s *Storage) LoadRefresh(code string) (*osin.AccessData, error) {
	row := s.db.QueryRow("SELECT access FROM refresh WHERE token=$1 LIMIT 1", code)
	var access string
	if err := row.Scan(&access); err == sql.ErrNoRows {
		return nil, errors.New(err.Error())
	} else if err != nil {
		return nil, errors.New(err.Error())
	}
	return s.LoadAccess(access)
}

// RemoveRefresh revokes or deletes refresh AccessData.
func (s *Storage) RemoveRefresh(code string) error {
	_, err := s.db.Exec("DELETE FROM refresh WHERE token=$1", code)
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}
