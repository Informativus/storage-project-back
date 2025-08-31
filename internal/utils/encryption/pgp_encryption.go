package encryption

import (
	"bytes"
	"io"
	"mime/multipart"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

func IsPGPEncrypted(fileHeader *multipart.FileHeader) (bool, error) {
	file, err := fileHeader.Open()

	if err != nil {
		return false, err
	}

	defer file.Close()

	buf := make([]byte, 64)
	n, err := file.Read(buf)

	if err != nil && err != io.EOF {
		return false, err
	}

	buf = buf[:n]

	if bytes.HasPrefix(buf, []byte("-----BEGIN PGP MESSAGE-----")) {
		return true, nil
	}

	if len(buf) > 0 && (buf[0]&0x80 != 0) {
		return true, nil
	}

	return false, nil
}

func EncryptFile(fileHeader *multipart.FileHeader, publicKey *multipart.FileHeader) ([]byte, error) {
	file, err := fileHeader.Open()

	if err != nil {
		return nil, err
	}

	defer file.Close()

	buf := new(bytes.Buffer)

	if _, err := io.Copy(buf, file); err != nil {
		return nil, err
	}

	pubFile, err := publicKey.Open()

	if err != nil {
		return nil, err
	}

	defer pubFile.Close()

	pubBuf := new(bytes.Buffer)
	if _, err := io.Copy(pubBuf, pubFile); err != nil {
		return nil, err
	}

	keyObj, err := crypto.NewKeyFromArmored(pubBuf.String())

	if err != nil {
		return nil, err
	}

	keyRing, err := crypto.NewKeyRing(keyObj)

	if err != nil {
		return nil, err
	}

	message := crypto.NewPlainMessage(buf.Bytes())

	encrypted, err := keyRing.Encrypt(message, nil)

	if err != nil {
		return nil, err
	}

	return encrypted.GetBinary(), nil
}
