//
// ideas from:
//  * https://gist.github.com/sohamkamani/08377222d5e3e6bc130827f83b0c073e
//  * https://stackoverflow.com/questions/10701874/generating-the-sha-hash-of-a-string-using-golang
//  * https://github.com/zaffka/rsa
//

package app

import (
	"crypto"
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

func ParseRsaPrivateKeyFromPemStr(privPEM []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privPEM)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the PrivateKey")
	}

	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	asRsa, ok := priv.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("failed to parse PrivateKey as RSA key")
	}

	return asRsa, nil
}

func ParseRsaPrivateKeyFromPemStrToPublicKey(privPEM []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(privPEM)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the PublicKey as PrivateKey")
	}

	priv, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	asRsa, ok := priv.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("failed to convert PrivateKey to PublicKey")
	}

	return asRsa, nil
}

func ParseRsaPublicKeyFromPemStr(pubPEM []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pubPEM)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the PublicKey")
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

func HexStringOfEncryptedHashValue(inputMessage, hashAlgo, keyPemFileName string, useSigning bool, optDebug ...int) (string, error) {
	var myHash hash.Hash
	var myCHash crypto.Hash
	myStringHash := ""
	var rsaPrivKey *rsa.PrivateKey
	var rsaPubKey *rsa.PublicKey
	var encryptedBytes []byte
	var err error

	var debug int = 0
	if len(optDebug) > 0 {
		debug = optDebug[0]
	}

	switch {
	case hashAlgo == "unhashed":
		if useSigning == true {
			myHash = nil
		} else {
			return "", errors.New("Hash algorithm only supported with PrivateKey.")
		}
		myCHash = crypto.Hash(0) // Note: crypto.Hash(0), unhashed payload
	case hashAlgo == "md5":
		myHash = md5.New()
		myCHash = crypto.MD5
	case hashAlgo == "sha1":
		myHash = sha1.New()
		myCHash = crypto.SHA1
	case hashAlgo == "sha256":
		myHash = sha256.New()
		myCHash = crypto.SHA256
	case hashAlgo == "sha512":
		myHash = sha512.New()
		myCHash = crypto.SHA512
	default:
		return "", errors.New("Hash algorithm is not supported.")
	}

	keyPemFileContent, err := ioutil.ReadFile(keyPemFileName)
	if err != nil {
		log.Fatal(err)
		return "", err
	} else {
		if debug >= 4 {
			log.Print("reading done of PEM File " + keyPemFileName)
			log.Print("got " + string(keyPemFileContent))
		}
	}

	if debug >= 2 {
		log.Print("Using PEM File: " + keyPemFileName)
		log.Print("Using HashAlgo: " + hashAlgo)
		log.Print("InputString   : " + inputMessage)
		log.Print("got " + string(keyPemFileContent))
	}

	if useSigning == true {
		//rsaPrivKey, err = ParseRsaPrivateKeyFromPemStrToPublicKey(keyPemFileContent)
		rsaPrivKey, err = ParseRsaPrivateKeyFromPemStr(keyPemFileContent)
	} else {
		rsaPubKey, err = ParseRsaPublicKeyFromPemStr(keyPemFileContent)
	}
	if err != nil {
		log.Fatal(err)
		return "", err
	} else {
		//log.Print("parsing done of PEM File " + keyPemFileName)
	}

	// generate encrypted hash
	var payload []byte = []byte(inputMessage)
	if useSigning == true && myHash != nil {
		// generate hash of input
		myHash.Write([]byte(inputMessage))
		payload = myHash.Sum(nil)
		if debug >= 2 {
			log.Print("HashSum(in Hex): " + hex.EncodeToString(payload))
		}
	} else {
		// use input as unhashed payload
		if debug >= 2 {
			log.Print("No Hash(plain): " + inputMessage)
		}
	}

	// sign or encrypt hash
	if useSigning == true {
		encryptedBytes, err = rsa.SignPSS(
			rand.Reader,
			rsaPrivKey,
			myCHash,
			payload,
			nil)
	} else {
		encryptedBytes, err = rsa.EncryptOAEP(
			myHash,
			rand.Reader,
			rsaPubKey,
			payload,
			nil)
	}
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	myStringHash = hex.EncodeToString(encryptedBytes)

	if debug >= 2 {
		log.Print("StringHash: " + myStringHash)
	}
	return myStringHash, nil
}
