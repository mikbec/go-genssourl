//
// ideas from:
//  * https://gist.github.com/sohamkamani/08377222d5e3e6bc130827f83b0c073e
//  * https://stackoverflow.com/questions/10701874/generating-the-sha-hash-of-a-string-using-golang
//

package app

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	//"fmt"
	"hash"
	"io/ioutil"
	"log"
)

func ParseRsaPublicKeyFromPemStr(pubPEM []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pubPEM)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	var cert *x509.Certificate
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := cert.PublicKey.(*rsa.PublicKey)
	return pub, nil

	// pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	//
	//	if err != nil {
	//		return nil, err
	//	}
	//
	// switch pub := pub.(type) {
	// case *rsa.PublicKey:
	//
	//	return pub, nil
	//
	// default:
	//
	//		break // fall through
	//	}
	//
	// return nil, errors.New("Key type is not RSA")
}

func HexStringOfHashValue(inputMessage string, hashAlgo string) (string, error) {
	var myHash hash.Hash
	myStringHash := ""

	switch {
	case hashAlgo == "md5":
		myHash = md5.New()
	case hashAlgo == "sha1":
		myHash = sha1.New()
	case hashAlgo == "sha256":
		myHash = sha256.New()
	case hashAlgo == "sha512":
		myHash = sha512.New()
	default:
		return "", errors.New("Hash algorithm is not supported.")
	}
	myHash.Write([]byte(inputMessage))
	myStringHash = hex.EncodeToString(myHash.Sum(nil))

	return myStringHash, nil
}

func HexStringOfEncryptedHashValue(inputMessage, hashAlgo, pubPemFileName string) (string, error) {
	var myHash hash.Hash
	myStringHash := ""

	switch {
	case hashAlgo == "md5":
		myHash = md5.New()
	case hashAlgo == "sha1":
		myHash = sha1.New()
	case hashAlgo == "sha256":
		myHash = sha256.New()
	case hashAlgo == "sha512":
		myHash = sha512.New()
	default:
		return "", errors.New("Hash algorithm is not supported.")
	}

	pubPemFileContent, err := ioutil.ReadFile(pubPemFileName)
	if err != nil {
		log.Fatal(err)
		return "", err
	} else {
		//log.Print("reading done of PEM File " + pubPemFileName)
		//log.Print("got " + string(pubPemFileContent))
		//fmt.Print(string(pubPemFileContent))
	}

	publicKey, err := ParseRsaPublicKeyFromPemStr(pubPemFileContent)
	if err != nil {
		log.Fatal(err)
		return "", err
	} else {
		//log.Print("parsing done of PEM File " + pubPemFileName)
	}

	encryptedBytes, err := rsa.EncryptOAEP(
		myHash,
		rand.Reader,
		publicKey,
		[]byte(inputMessage),
		nil)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	myStringHash = hex.EncodeToString(encryptedBytes)

	return myStringHash, nil
}
