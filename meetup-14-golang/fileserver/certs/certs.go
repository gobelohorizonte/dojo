/*
* Ukkbox Go Library (C) 2017 Inc.
*
* @project     Ukkbox
* @package     main
* @author      @jeffotoni
* @size        01/06/2017

* @description Our main auth will be responsible for validating our
* handlers, validating users and will also be in charge of creating users,
* removing them and doing their validation of access.
* We are using jwt to generate the tokens and validate our handlers.
* The logins and passwords will be in the AWS Dynamond database.
*
* $ openssl genrsa -out private.rsa 1024
* $ openssl rsa -in private.rsa -pubout > public.rsa.pub
*
 */

package certs

import (
	"crypto/rsa"
)

var (
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
)

const (
	RSA_PRIVATE = `-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQCbhxO3T2mwbScO5Jt+m63rjq2v53NeKUpxrsY9LsbpQawoqWDX
jjn1HX4MeXQsZFCv65/jhACqLD4O7m8K6J+TpDxoaFehIjRE9U5r4Z41MYj8qOBY
bkR39kTwN9Cx0FA3Ny42k2ui8FNLWPzIK3wWO4G6RNW3Xj7jkaWNiXpPFQIDAQAB
AoGBAINAv7bX4g3uUCQVcdSrdV9yDcqBva8dkaHXKZ3AuEVqEuxN5ViEwwzFUvcc
GJrOHfoZE9piMF1s8QKQ3k2KfABApU+T+vfdrdWUypR/RYjJSseBTl7369yNF8So
3QtqjFhxHMQQmRZsinObTCuQWR0y282YHeD76n58jvn3BEYRAkEAyzXj/XH0IVU1
3bBSnW1r4j7W2ed/DBRJyjOmxYB5JIQnbM4fyx9t4INx/eXOjjtI0mkODITr+bwX
8SBvFLfHqwJBAMPuIfOT9BiZWvPYyKCMzZxLIBGPXusECiN08T2TgHQkxPK7Ov2O
/lzFy4Rzwak6l8Yuwij4MpDMziY+U4vkhD8CQBEkAo8mRYlqPpjsfot451i2JDlN
JZJHQ9IieTa/l3hVDV5IJLZleEcvzzWzZLDqn0HgSqcTrzPpgbt9GOGOfvECQEIt
lNopBzGn2siyWFGiPXClD1ffDThkTOhc/37E64ZPRRaXlv095zx+spcyYh8+4zTV
Zk9gRfQSuS7BroZ50RECQDaeCOfF88sV8dwKp3ykhkDuuGRoGz2mr5MxJt3Ub7kf
3VrHTkNsQTr3VO3AORukA4NqAQnYsx13Ps6y5AHsNfM=
-----END RSA PRIVATE KEY-----
`

	RSA_PUBLIC = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCbhxO3T2mwbScO5Jt+m63rjq2v
53NeKUpxrsY9LsbpQawoqWDXjjn1HX4MeXQsZFCv65/jhACqLD4O7m8K6J+TpDxo
aFehIjRE9U5r4Z41MYj8qOBYbkR39kTwN9Cx0FA3Ny42k2ui8FNLWPzIK3wWO4G6
RNW3Xj7jkaWNiXpPFQIDAQAB
-----END PUBLIC KEY-----
`
)
