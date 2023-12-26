package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"testing"
)

func TestAES(t *testing.T) {
	// Llave de encriptación (debe tener 16, 24 o 32 bytes para AES-128, AES-192 o AES-256 respectivamente)
	key := []byte("aaaaaaaaaaaaaaaa")
	// Datos a encriptar
	data := []byte("Hola, mundo!")

	// Generar un nuevo bloque de cifrado AES
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	// Obtener el tamaño del bloque
	blockSize := block.BlockSize()

	// Calcular la cantidad de bytes de relleno necesarios
	padding := blockSize - (len(data) % blockSize)

	// Crear un slice con los datos a encriptar y el relleno
	plaintext := make([]byte, len(data)+padding)
	copy(plaintext, data)

	// Rellenar con bytes del valor del padding
	for i := len(data); i < len(plaintext); i++ {
		plaintext[i] = byte(padding)
	}

	// Generar un vector de inicialización único aleatorio (IV)
	iv := make([]byte, blockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err.Error())
	}
	fmt.Println()
	// Modo de operación de cifrado en bloque con relleno PKCS#7
	mode := cipher.NewCBCEncrypter(block, iv)

	// Agregar el IV al inicio de los datos encriptados
	encryptedData := make([]byte, blockSize+len(plaintext))
	copy(encryptedData[:blockSize], iv)

	// Encriptar los datos
	mode.CryptBlocks(encryptedData[blockSize:], plaintext)
	fmt.Println("KEY: ", key)
	fmt.Println("IV: ", iv)
	fmt.Println("Datos encriptados: ", encryptedData)

	// Convertir los datos encriptados a formato base64 para impresión
	encodedData := base64.StdEncoding.EncodeToString(encryptedData)

	fmt.Println("Datos encriptados (en base64): " + encodedData)
}
