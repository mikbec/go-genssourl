//
// ideas from:
//  * https://gist.github.com/sohamkamani/08377222d5e3e6bc130827f83b0c073e
//  * https://stackoverflow.com/questions/10701874/generating-the-sha-hash-of-a-string-using-golang
//  * https://github.com/zaffka/rsa
//  * https://stackoverflow.com/questions/40870178/golang-rsa-decrypt-no-padding
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
	"math/big"
)

type DoCryptoTask int

const (
	// since iota starts with 0, the first value
	// defined here will be the default
	DoCryptoTaskUndefined DoCryptoTask = iota
	DoCryptoTaskEncryption
	DoCryptoTaskSigning
	DoCryptoTaskEncRsaNoPadding
)

func (s DoCryptoTask) String() string {
	switch s {
	case DoCryptoTaskEncryption:
		return "encryption"
	case DoCryptoTaskSigning:
		return "signing"
	case DoCryptoTaskEncRsaNoPadding:
		return "enc_rsa_no_padding"
	}
	return "undefined"
}

func StringToDoCryptoTask(str string) DoCryptoTask {
	switch str {
	case "encryption":
		return DoCryptoTaskEncryption
	case "signing":
		return DoCryptoTaskSigning
	case "enc_rsa_no_padding":
		return DoCryptoTaskEncRsaNoPadding
	}
	return DoCryptoTaskUndefined
}

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

	switch hashAlgo {
	case "md5":
		myHash = md5.New()
	case "sha1":
		myHash = sha1.New()
	case "sha256":
		myHash = sha256.New()
	case "sha512":
		myHash = sha512.New()
	default:
		return "", errors.New("Hash algorithm is not supported.")
	}
	myHash.Write([]byte(inputMessage))
	myStringHash = hex.EncodeToString(myHash.Sum(nil))

	return myStringHash, nil
}

func HexStringOfEncryptedHashValue(inputMessage, hashAlgo, keyPemFileName string, doCryptoTaskString string, optDebug ...int) (string, error) {
	var myHash hash.Hash
	var myCHash crypto.Hash
	myStringHash := ""
	var rsaPrivKey *rsa.PrivateKey
	var rsaPubKey *rsa.PublicKey
	var encryptedBytes []byte
	var err error
	doCryptoTask := StringToDoCryptoTask(doCryptoTaskString)

	var debug int = 0
	if len(optDebug) > 0 {
		debug = optDebug[0]
	}

	if doCryptoTask == DoCryptoTaskUndefined {
		return "", errors.New("CryptoTask '" + doCryptoTaskString + "' is not supported.")
	}

	switch hashAlgo {
	case "unhashed":
		if doCryptoTask == DoCryptoTaskSigning {
			myHash = nil
		} else {
			return "", errors.New("Hash algorithm only supported with PrivateKey.")
		}
		myCHash = crypto.Hash(0) // Note: crypto.Hash(0), unhashed payload
	case "md5":
		myHash = md5.New()
		myCHash = crypto.MD5
	case "sha1":
		myHash = sha1.New()
		myCHash = crypto.SHA1
	case "sha256":
		myHash = sha256.New()
		myCHash = crypto.SHA256
	case "sha512":
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
		log.Print("Using PEM File  : " + keyPemFileName)
		log.Print("Doing CryptoTask: " + doCryptoTaskString)
		log.Print("Using HashAlgo  : " + hashAlgo)
		log.Print("InputString     : " + inputMessage)
		log.Print("got " + string(keyPemFileContent))
	}

	switch doCryptoTask {
	case DoCryptoTaskSigning:
		//rsaPrivKey, err = ParseRsaPrivateKeyFromPemStrToPublicKey(keyPemFileContent)
		rsaPrivKey, err = ParseRsaPrivateKeyFromPemStr(keyPemFileContent)
	default:
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
	if (doCryptoTask == DoCryptoTaskSigning || doCryptoTask == DoCryptoTaskEncRsaNoPadding) && myHash != nil {
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
	switch doCryptoTask {
	case DoCryptoTaskEncryption:
		encryptedBytes, err = rsa.EncryptOAEP(
			myHash,
			rand.Reader,
			rsaPubKey,
			payload,
			nil)
	case DoCryptoTaskSigning:
		encryptedBytes, err = rsa.SignPSS(
			rand.Reader,
			rsaPrivKey,
			myCHash,
			payload,
			nil)
	case DoCryptoTaskEncRsaNoPadding:
		// Do simple RSA encryption with Nonce and No Padding
		c := new(big.Int).SetBytes([]byte(payload))
		encryptedBytes = c.Exp(c, big.NewInt(int64(rsaPubKey.E)), rsaPubKey.N).Bytes()
	default:
		err = errors.New("CryptoTask '" + doCryptoTaskString + "' is not supported ... we should not get there.")
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
