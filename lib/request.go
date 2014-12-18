package lib

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/crypto/openpgp"
)

const (
	pgpNameHeader   = "X-SIGNER"
	signatureHeader = "X-SIGNATURE"
)

type Request struct {
	UserName      string
	IdentityName  string
	EntityList    openpgp.EntityList
	EncryptedBody []byte
	*http.Request
}

func NewRequest(method, urlStr string, body []byte) (*Request, error) {
	r, err := http.NewRequest(method, urlStr, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	return &Request{"", "", nil, body, r}, nil
}

func AuthenticatedRequest(r *http.Request, es *EntityStorage) (*Request, error) {
	identityName := r.Header.Get(pgpNameHeader)
	if identityName == "" {
		return nil, errors.New("Missing PGPName Header")
	}

	entity := es.EntityForIdentity(identityName)
	if entity == nil {
		return nil, errors.New("Unable to find EntityList for provided IdentityName")
	}

	var body []byte
	if r.Body != nil {
		body, _ = ioutil.ReadAll(r.Body)
	}

	request := &Request{entity.User, identityName, entity.EntityList, body, r}
	err := request.CheckSignature()
	if err != nil {
		return nil, err
	}

	return request, nil
}

func (r *Request) CheckSignature() error {
	signature := r.Header.Get(signatureHeader)
	if signature == "" {
		return errors.New("Missing Signature Header")
	}

	data, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return errors.New("Invalid Signature")
	}

	signature = string(data)
	message := r.SignatureMessage()
	signer, err := openpgp.CheckDetachedSignature(
		r.EntityList,
		bytes.NewBufferString(message),
		bytes.NewBufferString(signature))

	if err != nil {
		return err
	}

	matchedIdentity := false
	for identityName := range signer.Identities {
		if identityName == r.IdentityName {
			matchedIdentity = true
			break
		}
	}

	if !matchedIdentity {
		return errors.New("Signer Identity does not matched provided value")
	}

	return nil
}

func (r *Request) Sign(entity *Entity) error {
	r.IdentityName = entity.GetIdentity()

	signature, err := r.Signature(entity)
	if err != nil {
		return err
	}

	r.Header.Add(pgpNameHeader, r.IdentityName)
	r.Header.Add(signatureHeader, signature)

	return nil
}

func (r *Request) Signature(entity *Entity) (string, error) {
	buf := new(bytes.Buffer)
	message := bytes.NewBufferString(r.SignatureMessage())
	if err := openpgp.DetachSign(buf, entity.PGPEntity(), message, nil); err != nil {
		return "", err
	}

	signature, err := ioutil.ReadAll(buf)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

func (r *Request) SignatureMessage() string {
	message := fmt.Sprintf("%s %s %s=%s %s", r.Method, r.URL.Path,
		pgpNameHeader, r.IdentityName, r.EncryptedBody)

	return strings.TrimRight(message, " ")
}
