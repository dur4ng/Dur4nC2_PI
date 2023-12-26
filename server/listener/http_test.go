package listener

import (
	_implantCrypto "Dur4nC2/implant/cryptography"
	_implantHTTPClient "Dur4nC2/implant/transport"
	_serverCrypto "Dur4nC2/misc/crypto"
	"Dur4nC2/misc/protobuf/implantpb"
	"Dur4nC2/server/db"
	beaconRepository "Dur4nC2/server/domain/beacon/repository/postgres"
	_implantRepository "Dur4nC2/server/domain/implant/respository/postgres"
	"Dur4nC2/server/domain/models"
	taskRepository "Dur4nC2/server/domain/task/repository/postgres"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/nacl/box"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"
)

var _ = Describe("Listeners", func() {
	BeforeEach(func() {

	})
	Context("Connections", func() {
		It("Initial connection", func() {
			//cryptography.ECCServerKeyPair()
			/*
				keypair := cryptography.ECCServerKeyPair()
				implantCrypto.SetSecrets("", "", "", keypair.PublicBase64(), "", "")
				server, err := StartHTTPListener(&HTTPListenerConfig{
					Addr: "127.0.0.1:8000",
				})
				if err != nil {
					fmt.Println("Listener failed to start ", err)
					return
				}
				Ω(err).To(Succeed())
				go server.HTTPServer.ListenAndServe()
				if err != nil {
					fmt.Println("Listener failed to serve ", err)
					return
				}
				Ω(err).To(Succeed())
				baseURL := &url.URL{
					Scheme: "transport",
					Host:   "127.0.0.1:8000",
					Path:   fmt.Sprintf("/login"),
				}
				fmt.Println("[b] ", baseURL.String())
				_, err = _implantHTTPClient.HTTPStartConnection(baseURL.String(), "")
				Ω(err).To(Succeed())
				if err != nil {
					return
				}
				server.Cleanup()

			*/
		})
		It("Gin test", func() {
			r := gin.Default()

			// Ping test
			r.GET("/ping", func(c *gin.Context) {
				c.String(http.StatusOK, "pong")
			})

			// Get user value
			r.GET("/user/:name", func(c *gin.Context) {
				user := c.Params.ByName("name")
				c.JSON(http.StatusOK, gin.H{"user": user, "status": "bbb"})
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/ping", nil)

			r.ServeHTTP(w, req)
			Ω(w.Code).To(Equal(http.StatusOK))
		})
		It("mux test", func() {
			server, err := StartHTTPListener(&HTTPListenerConfig{
				Addr: "127.0.0.1:8000",
			})
			if err != nil {
				fmt.Println("Listener failed to start ", err)
				return
			}
			Ω(err).To(Succeed())
			go server.HTTPServer.ListenAndServe()
			if err != nil {
				fmt.Println("Listener failed to serve ", err)
				return
			}
			time.Sleep(time.Second)

			req, reqErr := http.NewRequest(http.MethodGet, "http://localhost:8000/test", nil)
			if reqErr != nil {
				fmt.Errorf("there was an error building the HTTP request:\r\n%s", reqErr.Error())
			}
			req.Close = true
			var transport http.RoundTripper
			transport = &http.Transport{
				MaxIdleConns:    10,
				IdleConnTimeout: 1 * time.Nanosecond,
			}
			s := &http.Client{Transport: transport}
			resp, err := s.Do(req)
			if err != nil {
				fmt.Println("Error in HTTPClient.Do(req): ", err)
			}
			Ω(resp.StatusCode).Should(Equal(http.StatusOK))
		})
		It("Connection initialization", func() {
			server_keypair := _serverCrypto.ECCServerKeyPair()
			implant_keypair, err := _implantRepository.NewPosgresImplantRepository(db.Session()).Create()
			fmt.Println("digest creaded = ", implant_keypair.ECCPublicKeyDigest)
			_implantCrypto.SetSecrets(implant_keypair.ECCPublicKey, implant_keypair.ECCPrivateKey, implant_keypair.ECCPublicKeyDigest, server_keypair.PublicBase64(), "", "", server_keypair.PrivateBase64())
			server, err := StartHTTPListener(&HTTPListenerConfig{
				Addr:   "127.0.0.1:8000",
				Domain: "127.0.0.1",
			})
			if err != nil {
				fmt.Println("Listener failed to start ", err)
				return
			}
			Ω(err).To(Succeed())
			go server.HTTPServer.ListenAndServe()
			if err != nil {
				fmt.Println("Listener failed to serve ", err)
				return
			}
			time.Sleep(1 * time.Second)

			baseURL := &url.URL{
				Scheme: "transport",
				Host:   "127.0.0.1:8000",
				Path:   fmt.Sprintf("/test"),
			}
			fmt.Println("[b] ", baseURL.String())
			config := &models.ImplantConfig{
				Domain:     "127.0.0.1",
				URL:        "http://127.0.0.1:8000",
				PathPrefix: "",
			}
			client, err := _implantHTTPClient.HTTPStartConnection(*config)
			Ω(err).To(Succeed())
			if err != nil {
				return
			}

			fmt.Println("[test] Sending BeaconRegister...")
			register := implantpb.Register{
				Name:     "test",
				Hostname: "host1234",
				Username: "dur4n",
				Os:       "Linux",
			}
			beaconRegister := implantpb.BeaconRegister{
				Jitter:   5,
				Interval: 3,
				Register: &register,
			}
			registerData, err := proto.Marshal(&beaconRegister)
			Ω(err).To(Succeed())
			envelope := implantpb.Envelope{
				ID:   1234,
				Type: implantpb.MsgRegister,
				Data: registerData,
			}
			responseEnvelope, err := client.WriteEnvelope(&envelope)
			client.ID = string(responseEnvelope.Data)
			Ω(err).To(Succeed())

			//create Whoami task flow ---------------
			beaconRepo := beaconRepository.NewPostgresBeaconRepository(db.Session())
			taskRepo := taskRepository.NewPostgresBeaconTaskRepository(db.Session())
			beacons, err := beaconRepo.List()
			index := len(beacons) - 1
			envelope = implantpb.Envelope{
				ID:   1444,
				Type: implantpb.MsgWhoamiReq,
			}
			data, err := proto.Marshal(&envelope)
			task := models.BeaconTask{
				EnvelopeID: envelope.ID,
				BeaconID:   beacons[index].ID,
				State:      models.SENT,
				Request:    data,
			}
			err = taskRepo.Create(beacons[index], task)
			Ω(err).To(Succeed())
			envelope = implantpb.Envelope{
				ID:   1235,
				Type: implantpb.MsgBeaconTasks,
			}
			tasksData, err := client.WriteEnvelope(&envelope)
			Ω(err).To(Succeed())
			tasks := &implantpb.BeaconTasks{}
			proto.Unmarshal(tasksData.Data, tasks)
			fmt.Printf("Envelope type %d", tasks.Tasks[0].Type)

			server.Cleanup()
		})
		It("ECC crypto", func() {
			server_keypair := _serverCrypto.ECCServerKeyPair()
			implant_keypair, _ := _implantRepository.NewPosgresImplantRepository(db.Session()).Create()

			//msg := "hello"
			plaintext := _implantCrypto.RandomKey()
			httpSessionInit := &implantpb.HTTPSessionInit{Key: plaintext[:]}
			httpSessionInitData, _ := proto.Marshal(httpSessionInit)

			implant_privatekeyraw, _ := base64.RawStdEncoding.DecodeString(implant_keypair.ECCPrivateKey)
			implant_public_key_raw, _ := base64.RawStdEncoding.DecodeString(implant_keypair.ECCPublicKey)
			var implant_public_key [32]byte
			copy(implant_public_key[:], implant_public_key_raw)

			// encrypt
			var nonce [24]byte
			if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
				fmt.Println("Error in reader")
			}
			var implant_privatekey [32]byte
			copy(implant_privatekey[:], implant_privatekeyraw)
			ciphertext := box.Seal(nonce[:], httpSessionInitData, &nonce, server_keypair.Public, &implant_privatekey)

			digest := sha256.Sum256((implant_public_key)[:])
			msgEncrypted := make([]byte, 32+len(ciphertext))
			copy(msgEncrypted, digest[:])
			copy(msgEncrypted[32:], ciphertext)
			var hexData = make([]byte, hex.EncodedLen(len(msgEncrypted)))
			hex.Encode(hexData, msgEncrypted)

			// decrypt
			if len(ciphertext) < 24 {
				fmt.Println(errors.New("ciphertext too short"))
			}
			var decryptNonce [24]byte
			copy(decryptNonce[:], ciphertext[:24])
			plaintextMsg, ok := box.Open(nil, ciphertext[24:], &decryptNonce, &implant_public_key, server_keypair.Private)
			if !ok {
				fmt.Println("Decryption failed")
			}
			//Ω(string(plaintext)).Should(Equal(msg))

			ciphertext = msgEncrypted[32:]
			copy(decryptNonce[:], ciphertext[:24])
			plaintextMsg, ok = box.Open(nil, ciphertext[24:], &decryptNonce, &implant_public_key, server_keypair.Private)
			if !ok {
				fmt.Println("Decryption failed")
			}
			sessionInit := &implantpb.HTTPSessionInit{}
			proto.Unmarshal(plaintextMsg, sessionInit)
			//Ω(string(sessionInit.GetKey())).Should(Equal(msg))

			//var rawData = make([]byte, hex.EncodedLen(len(hexData)))
			//_, err := hex.Decode(rawData, hexData)
			//Ω(err).To(Succeed())
			//ciphertext = rawData[32:]
			ciphertext = msgEncrypted[32:]
			copy(decryptNonce[:], ciphertext[:24])
			plaintextMsg, ok = box.Open(nil, ciphertext[24:], &decryptNonce, &implant_public_key, server_keypair.Private)
			if !ok {
				fmt.Println("Decryption failed")
			}
			sessionInit = &implantpb.HTTPSessionInit{}
			err := proto.Unmarshal(plaintextMsg, sessionInit)
			Ω(err).To(Succeed())
			Ω(sessionInit.GetKey()).Should(Equal(httpSessionInit.GetKey()))
		})
		/*
			It("Connection handler", func() {
				server_keypair := _serverCrypto.ECCServerKeyPair()
				implant_keypair, err := _implantRepository.NewPosgresImplantRepository(db.Session()).Create()
				fmt.Println("digest creaded = ", implant_keypair.ECCPublicKeyDigest)
				_implantCrypto.SetSecrets(implant_keypair.ECCPublicKey, implant_keypair.ECCPrivateKey, implant_keypair.ECCPublicKeyDigest, server_keypair.PublicBase64(), "", "", server_keypair.PrivateBase64())
				server, err := StartHTTPListener(&HTTPListenerConfig{
					Addr: "127.0.0.1:8000",
				})
				if err != nil {
					fmt.Println("Listener failed to start ", err)
					return
				}
				Ω(err).To(Succeed())
				go server.HTTPServer.ListenAndServe()
				if err != nil {
					fmt.Println("Listener failed to serve ", err)
					return
				}
				time.Sleep(1 * time.Second)

				baseURL := &url.URL{
					Scheme: "transport",
					Host:   "127.0.0.1:8000",
					Path:   fmt.Sprintf("/login.html"),
				}
				fmt.Println("[b] ", baseURL.String())
				_, err = _implantHTTPClient.HTTPStartConnection(baseURL.String(), "")
				Ω(err).To(Succeed())

				server.Cleanup()
			})

		*/
	})
})
