package tls

const testCAcert = `-----BEGIN CERTIFICATE-----
MIIFojCCA4qgAwIBAgICK2cwDQYJKoZIhvcNAQELBQAwcjELMAkGA1UEBhMCRUMx
EjAQBgNVBAgTCVBpY2hpbmNoYTEOMAwGA1UEBxMFUXVpdG8xFjAUBgNVBAkTDUZh
a2Ugc3RyZWV0IDExDzANBgNVBBETBjIyMjIyMjEWMBQGA1UEChMNQ295SU0gVGVz
dCBDQTAeFw0yMTAyMjgxNDEyMzZaFw0zNjAyMjgxNDEyMzZaMHIxCzAJBgNVBAYT
AkVDMRIwEAYDVQQIEwlQaWNoaW5jaGExDjAMBgNVBAcTBVF1aXRvMRYwFAYDVQQJ
Ew1GYWtlIHN0cmVldCAxMQ8wDQYDVQQREwYyMjIyMjIxFjAUBgNVBAoTDUNveUlN
IFRlc3QgQ0EwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQCwwmFR3d9k
Hv7eTW8NGOrI7BPdVZmp5wNQPi7FGwq+SbpF2KLvNIFXxQifnX/NuauTlko5oYTS
Q2RC7vraW6q05v9DMdBewyGtZug+Xvd1S0D5YDkew6OlPfxp3Sc0h10zqEnP1Un9
dWzEovPdJp1BCh0OiFcaXTVWWZfrdoY1U1/j3sRXxWrM37bmnT1rzTeiGW1qaq7p
UdWDFRGWQhXvBGryG4Nof/BblHGs6mA2SG+BsZYQYEEAJOgk5Inq4F8HvltR3K2z
F9X41PMpANivQPB1DDrR46XGoUjOKaSyh5w0kn1vED6eFRJT+Yx9s090lf4xduAV
9xbCYS/UiJFiCZcYPqlS42vIYDAiojFVZJhP2/csYl6oOGxOfGIxTl1OFGZlhkOK
LJTY6h9wabOyg23Xng8haUa9rvl1UdgRi2CNS7CxKqd8HxsI92nJuqQReopticUw
CRBCz7/qCbS01lgSJ1LHQP2J6CAzcR74RKUD83OmK0Ktpo/bdxvLb6YrvBU4Qh8T
69DFq+Lpj5Ya+wzVKrq0AykgiaYXLnynrE7p5CiEtPjtQMPQdCS0uRDb8QB04tt8
MJXezIQwgdH7jgHfQYf/dW1j4aAMCcGSZhykjigXzqFC7zxPYP1cajlQrVLj2IXJ
9cueAHyRxzCZCeMSL1nLsPr1ZgPOw5SLQQIDAQABo0IwQDAOBgNVHQ8BAf8EBAMC
AoQwHQYDVR0lBBYwFAYIKwYBBQUHAwIGCCsGAQUFBwMBMA8GA1UdEwEB/wQFMAMB
Af8wDQYJKoZIhvcNAQELBQADggIBABlP5vr1yRsI5trhvCT3MMUpgyI6BcePeJZo
MzFtGsYOMr6F2sPF80iUjJ/zykzb+JSmsDgDrhGOtrOl+8yyNMW0kHycse6/jPl5
ApCpvTkdR6RSiXb0nvHa+JSCjPXO75ur1KokFaI3u1Z7RKk5gLBAMYy3viGenhhA
CGIye2EG0fZ1bmLhZQwgRHzHYgDmPR67slyMlg5+ImXS102kf2E5Bm8Ahw36lobm
5GJeDGZ++tXtemR4skdgLCuwZ0JtFolpS/vrKCl7OgEhb/ZXX94om7JRQLh6raM1
EkUD7e10vgnGKrUp5v4TWNjow0wUgT+caNVT6WJBI07PPWmDiRXg+AoYhaKLc1pX
9LBohL5Y5WKK/zgHp1+wCEb+fRMOSAjq92+i20cKqBEHyHC55q15rYiwzx+ir9py
+n4xdYO7FPrS01B7jOpk5IwRgprgH0EIBI9cTjbTUzWCFHMbrqc/1Hf0BX79jSy3
oGht6WQuFgm1b8/CGxKiMDtGSSnwCzH49GWIKH2pe28YMj/36QLAxHaU4hGvRmDU
SVFluKgUkJEMrUNFIQ/IhYjaQeBvQ7sfi6Kb2CENdNzVspyV/St5xVEbKQKrBSLs
yHX1cHt0L+r1JhINJc/LIBvD/d1w+zuVNaqwXe3uJG0O5m8inLrfiRm4flehcVnC
tGQj9IWB
-----END CERTIFICATE-----
`

const testCAprivKey = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
`

const testCert = `-----BEGIN CERTIFICATE-----
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
NkWyCPyUQ8lK34BD6w==
-----END CERTIFICATE-----
`

const testPrivKey = `-----BEGIN RSA PRIVATE KEY-----
MIIJJwIBAAKCAgEAn2qXP8su9pdcT6s6I9cxVOszZ2nNyML2K+4EA+PIYJN3/r2l
DU92kMyrAB0eOBPgMqE2R0wBhqoqw0sQJnHbptqM1SMSD7hFXbqBMB9pujoEW2TG
Vn+BBfjuy6A1u/HR85ybBbaO4z57cQmI2TP4JCLEor40oCJxtiFT7RIKmdJ2SGqC
75SSByzLOWbg/s7exyFnQy5k9bSJUD9fCLppBsGoX2hZP1v2teYwcHMM6a+USVaK
nE9hQy/DDLJ66FyfUHRlriaD54orahOINNURM09OIByzQd/qW3VQUMpBYIDO457c
Hcm+vyH+knCDTWSdMG8rGmxC59kbhgUn6jCu/PT+y4wFUMjmCxsf75vRYGaZqwcq
r4uHmul7JXTUXfD9EenXLqfMzMTD1fnLDrdGve+in9Tt3jXfq7hlZDaj3qiZ7R4k
DtJUjk9K8zxy/xl03o4sWV9ws8sIUV3x3tbhFnexQ+97GwG+EPNIdcbRpHo4Z1AE
3deo6OeEY7MhBxOVs6+q7CwfKKVTn3aPxDfop0Xqf9Qq7aCH7f+O/2Yo3NbqiiLx
Ou1fwSWGC7DhtE4e34ellcImPRS5Vwb0JS39KdffCjece+Fc1X/2soQOr2G/JrWI
dc9eHsyWpaZ8Sot8uP4fNM5QTejDxJHIofUxFvMlgs4zmPS9PLpMJYOLJI0CAwEA
AQKCAgACU+8jelcUOL+bVjfCIDlTMSAOCYh8vwQTPiWG3QOnDWA6MxC+8gMcODDj
DonLbdbfRmVhgyWejsuTEHyK4yy+8gAOeLWhzyIMLVYHmt3TX1eC8iTHTJNYv/rU
tGE0fmJ/eTD2U2Ugwl/RFb+O1GhyNqPCcJ6aHAanDzOHibTn7B/YDN4em3/KZQgO
rYbpkaHFLKKyY3IL+Hfs2RANM5OnCprn0cFD4JborxTT/4oXu32h2Iaro6ka7w6d
F9odnISjCyAU+/D/J5Bcuy5I/zeCFU1hwKmJc7ibX0ot89Yij571yfMS6EhFyDxM
bSIttiNpeqYZe606b3wsZ9TeYZmciBvi7OcyJMVNsmoOdi4/rNPLa8A0FNFW+iRu
JHxpbz4dizYqS4jcSXFiXQFaoMVfwVDCm1G0AYCHawbZIiEBwBYWvoD2DM9wpFUP
gJGa4kWoiANovH8Vgf+e7nrzJ5CJ3JYMjOvN+sfR4MpEQ9DSHdHvq4wLGPn52Leu
g9buwfj5K+ONdY3XvTSLdVcemdNBiiEBY22ecALIsPJQZIAvYXDPIh9J65aBq78E
PUTkkH5vEzgpwiejSFbz98m3UwYLEtzJhNxORlOEG770wDujhqhV3IxfpfSzViK1
vVP6WQce+xRdZhBMh6wPoFIJHFsAc0i4OYWcqbvMlcJRqdjIwQKCAQEA0bYvoaiE
+aY9sTjyQawih/J8iofDwm2oHWIl3x0R1mvjr091aRHsZjG6oSpjPy82FPMi57S8
m5N6Qarc6tNbTHVcaCVGHYyfyQ9qDEP0032kPjeS0s/UhUbbFuusv0iPdNKOOJMY
oQWzi8bc9E4c1Qq5eFcC37djRMrpssd1LZQhD/gZg5K5hXvXsyUs1A6VQ4an+3ed
ZRA9sIL5GENOIE11OnsLgRQ+KNCBbbJrBUzj36SdFBbfgJVDxsjy8D79cegzqsia
LtwpPrrGt1egClTQGUB9QUehGVRbaGg4a3o7n/oI62ZG1J+rvbA9bV2UU5XH3OuH
phFFjs0fHH55EQKCAQEAwpp0ph3rpryZ3HWaC5PCQ2A6HG+a/n0/LI6fkuSGx36D
NP2OhmELFrNUOFXd/mCvhxY61RCmw54SBG8MyTxZVk0A50tou99Z2pNSLozP/ZhR
I0yVWqdKhNoXSXak+w4dUHg4Q4SGXhlKyJIq21GGZ18JMvli2sQYCqRIGz92Xv1B
Eoft6kYAogptpU0uizaHVwrk1o74DjNJqO0rzeA9OaOdnskZhVN3rzeM3G2ore/K
h04EoeqZxdVY3uHe2Zk7Z1bE/TnDkWNmGVvjdN62/IBm9ZBpWl3WBRoysS4a83qY
mXA2NFDbe8YfrnfGmXOHXKPsoiT5rLWCqC2XWDeTvQKCAQAB64EUIc7V2kfGT5co
MsM+K2IogoWwSgC4BCYEnOeE5wf2muugQqG/bcUfpJu0AGKmXnN7W5Q+eGMuJrpP
DBBR6uElsvGpY5gy5wk5g4XCSewvBaM6etyfO77VvuKd/bQShbr3maEoGD1EklWD
hxOMf8Si7WkBU1R9VL4+/MR93lVPKB5Trgw0xKV85mI6rsd/DsSK8NVoD3YBH7HY
HwWgFhV0q5u3WtAW35HPx0pjigisC33EqVDyhGtSbpSKzojTeiS+84c11p4qDNu0
4gB9F7mwAX8kEdvPt43+rrWVhlD1bfyW6yDK4YtY+TwWvDyXZ0+lHiLnylCwtgAK
6r2BAoIBAFSl0NWtMCLf6OFfejlM9XRPOBfEaBwIqOEdzMWdiA7gtfvnywYi0ir2
qEy09RJARjmxbrfdPVzbtiSdlWc3S/jhF+KEB7Oo7LHJ4TaEY7iAd9Kt7k13dU+i
efynkg3uTswA7yBXVgc6YzApfGDX7mmqihrVJa3ZHEgMu5y2lyusZ5DC9bcw6feS
J61+jB9cAbTX9UBrAfVTU9gaCjLMNnWK+PXnraUz8FyUAj6jqHq4UlVWl2dC386R
Bc41W7U1FQTXVmp7pNjp7rBbKu5cLiZZR+/K+Dipln2zrpcpYenEyvn7OGi7Py1w
ubkvOoDnItsmJrlE8iGw9ntnEWz7B9UCggEAHpkFsKZZbTi1JvTng69wCFdjsOLS
JJwH8n0/RGom4EKbja/icotHTS8t9RVRiFrhkBfgewnAb2HkrNfFdxuSQalIt+47
nm6iJJb7tJ5bQ6+wVywUYWXg67+D+Ga/AqAtzv995ZrvfMCEQF7rGpPADIAXP4gF
oxer4bQdwuOPyE3km/tnELQuimEcduSgVaivIidZm11Vjao3bfgIiEVJd2JdQDuf
cYtdD1oj5sujvYmeNbi/A5Wd7wIZg3/bLy8j65dihKWApG2BKrAroOAncwJuu32c
y4Z6a74qt65usxVZPA2yLAQRom7hp5zW421kbqskGRFEtW+XSFhhWVtfJA==
-----END RSA PRIVATE KEY-----
`

// THIS IS THE CODE THAT WAS USED TO GENERATE THE ABOVE DATA:

// 	ca := &x509.Certificate{
// 		SerialNumber: big.NewInt(11111),
// 		Subject: pkix.Name{
// 			Organization:  []string{"CoyIM Test CA"},
// 			Country:       []string{"EC"},
// 			Province:      []string{"Pichincha"},
// 			Locality:      []string{"Quito"},
// 			StreetAddress: []string{"Fake street 1"},
// 			PostalCode:    []string{"222222"},
// 		},
// 		NotBefore:             time.Now(),
// 		NotAfter:              time.Now().AddDate(15, 0, 0),
// 		IsCA:                  true,
// 		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
// 		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
// 		BasicConstraintsValid: true,
// 	}

// 	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
// 	if err != nil {
// 		panic(err)
// 	}

// 	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
// 	if err != nil {
// 		panic(err)
// 	}

// 	caPEM := new(bytes.Buffer)
// 	pem.Encode(caPEM, &pem.Block{
// 		Type:  "CERTIFICATE",
// 		Bytes: caBytes,
// 	})

// 	caPrivKeyPEM := new(bytes.Buffer)
// 	pem.Encode(caPrivKeyPEM, &pem.Block{
// 		Type:  "RSA PRIVATE KEY",
// 		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
// 	})

// 	fmt.Printf("%s\n", caPEM)
// 	fmt.Printf("%s\n", caPrivKeyPEM)

// 	block, _ := pem.Decode([]byte(testCAcert))
// 	ca, _ := x509.ParseCertificate(block.Bytes)

// 	block, _ = pem.Decode([]byte(testCAprivKey))
// 	caPrivKey, _ := x509.ParsePKCS1PrivateKey(block.Bytes)

// 	cert := &x509.Certificate{
// 		SerialNumber: big.NewInt(2222),
// 		Subject: pkix.Name{
// 			Organization:  []string{"CoyIM test cert"},
// 			Country:       []string{"EC"},
// 			Province:      []string{"Pichincha"},
// 			Locality:      []string{"Quito"},
// 			StreetAddress: []string{"Somewhere 123"},
// 			PostalCode:    []string{"9999"},
// 		},
// 		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
// 		DNSNames:     []string{"foo.bar.com", "coy.im"},
// 		NotBefore:    time.Now(),
// 		NotAfter:     time.Now().AddDate(15, 0, 0),
// 		SubjectKeyId: []byte{1, 2, 3, 4, 6},
// 		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
// 		KeyUsage:     x509.KeyUsageDigitalSignature,
// 	}

// 	certPrivKey, e4 := rsa.GenerateKey(rand.Reader, 4096)
// 	if e4 != nil {
// 		fmt.Printf("e: %v\n", e4)
// 		panic(e4)
// 	}
// 	certBytes, e5 := x509.CreateCertificate(rand.Reader, cert, ca, &certPrivKey.PublicKey, caPrivKey)
// 	if e5 != nil {
// 		fmt.Printf("e: %v\n", e5)
// 		panic(e5)
// 	}

// 	certPEM := new(bytes.Buffer)
// 	pem.Encode(certPEM, &pem.Block{
// 		Type:  "CERTIFICATE",
// 		Bytes: certBytes,
// 	})

// 	certPrivKeyPEM := new(bytes.Buffer)
// 	pem.Encode(certPrivKeyPEM, &pem.Block{
// 		Type:  "RSA PRIVATE KEY",
// 		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
// 	})

// 	fmt.Printf("%s\n", certPEM)
// 	fmt.Printf("%s\n", certPrivKeyPEM)
