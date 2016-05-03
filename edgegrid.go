// Package edgegrid provides the Akamai {OPEN} Edgegrid Authentication scheme
package edgegrid

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tuvistavie/securerandom"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"
	"unicode"
)

// Base struct provides all the neccessary fields to
// create authorization header, debug is optional
type Base struct {
	ClientToken  string
	ClientSecret string
	AccessToken  string
	HeaderToSign []string
	MaxBody      int
	Debug        bool
}

// Must be assigned the UTC time when the request is signed.
// Format of “yyyyMMddTHH:mm:ss+0000”
func makeEdgeTimeStamp() string {
	local := time.FixedZone("GMT", 0)
	t := time.Now().In(local)
	return fmt.Sprintf("%d%02d%02dT%02d:%02d:%02d+0000",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

// Must be assigned a nonce (number used once) for the request.
// It is a random string used to detect replayed request messages.
// A GUID is recommended.
func createNonce() string {
	uuid, err := securerandom.Uuid()
	if err != nil {
		log.Errorf("Generate Uuid failed, %s", err)
		return ""
	}
	return uuid
}

func stringMinifier(in string) (out string) {
	white := false
	for _, c := range in {
		if unicode.IsSpace(c) {
			if !white {
				out = out + " "
			}
			white = true
		} else {
			out = out + string(c)
			white = false
		}
	}
	return
}

// createSignature is the base64-encoding of the SHA–256 HMAC of the data to sign with the signing key.
func createSignature(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func createHash(data string) string {
	h := sha256.Sum256([]byte(data))
	return base64.StdEncoding.EncodeToString(h[:])
}

func (b *Base) canonicalizeHeaders(req *http.Request) string {
	var unsortedHeader []string
	var sortedHeader []string
	for k := range req.Header {
		unsortedHeader = append(unsortedHeader, k)
	}
	sort.Strings(unsortedHeader)
	for _, k := range unsortedHeader {
		for _, sign := range b.HeaderToSign {
			if sign == k {
				v := strings.TrimSpace(req.Header.Get(k))
				sortedHeader = append(sortedHeader, fmt.Sprintf("%s:%s", strings.ToLower(k), strings.ToLower(stringMinifier(v))))
			}
		}
	}
	return strings.Join(sortedHeader, "\t")

}

// signingKey is derived from the client secret.
// The signing key is computed as the base64 encoding of the SHA–256 HMAC of the timestamp string
// (the field value included in the HTTP authorization header described above) with the client secret as the key.
func (b *Base) signingKey(timestamp string) string {
	key := createSignature(timestamp, b.ClientSecret)
	return key
}

// The content hash is the base64-encoded SHA–256 hash of the POST body.
// For any other request methods, this field is empty. But the tab separator (\t) must be included.
// The size of the POST body must be less than or equal to the value specified by the service.
// Any request that does not meet this criteria SHOULD be rejected during the signing process,
// as the request will be rejected by EdgeGrid.
func (b *Base) createContentHash(req *http.Request) string {
	var (
		contentHash  string
		preparedBody string
		bodyBytes    []byte
	)
	if req.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(req.Body)
		req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		preparedBody = string(bodyBytes)
	}

	log.Debugf("Body is %s", preparedBody)
	if req.Method == "POST" && len(preparedBody) > 0 {
		log.Debugf("Signing content: %s", preparedBody)
		if len(preparedBody) > b.MaxBody {
			log.Debugf("Data length %d is larger than maximum %d",
				len(preparedBody), b.MaxBody)

			preparedBody = preparedBody[0:b.MaxBody]
			log.Debugf("Data truncated to %d for computing the hash", len(preparedBody))
		}
		contentHash = createHash(preparedBody)
	}
	log.Debugf("Content hash is '%s'", contentHash)
	return contentHash
}

// The data to sign includes the information from the HTTP request that is relevant to ensuring that the request is authentic.
// This data set comprised of the request data combined with the authorization header value (excluding the signature field,
// but including the ; right before the signature field).
func (b *Base) signingData(req *http.Request, authHeader string) string {

	dataSign := []string{
		req.Method,
		req.URL.Scheme,
		req.URL.Host,
		req.URL.Path + req.URL.RawQuery,
		b.canonicalizeHeaders(req),
		b.createContentHash(req),
		authHeader,
	}
	log.Debugf("Data to sign %s", fmt.Sprintf(strings.Join(dataSign, "\t")))
	return fmt.Sprintf(strings.Join(dataSign, "\t"))
}

func (b *Base) signingRequest(req *http.Request, authHeader string, timestamp string) string {
	return createSignature(b.signingData(req, authHeader),
		b.signingKey(timestamp))
}

// The Authorization header starts with the signing algorithm moniker (name of the algorithm) used to sign the request.
// The moniker below identifies EdgeGrid V1, hash message authentication code, SHA–256 as the hash standard.
// This moniker is then followed by a space and an ordered list of name value pairs with each field separated by a semicolon.
func (b *Base) createAuthHeader(req *http.Request, timestamp string, nonce string) string {
	authHeader := fmt.Sprintf("EG1-HMAC-SHA256 client_token=%s;access_token=%s;timestamp=%s;nonce=%s;",
		b.ClientToken,
		b.AccessToken,
		timestamp,
		nonce,
	)
	log.Debugf("Unsigned authorization header: '%s'", authHeader)

	signedAuthHeader := fmt.Sprintf("%ssignature=%s", authHeader, b.signingRequest(req, authHeader, timestamp))

	log.Debugf("Signed authorization header: '%s'", signedAuthHeader)
	return signedAuthHeader
}

// MakeHeader sets the authorization header to use Akamai Open API
func MakeHeader(b Base, req *http.Request) *http.Request {
	if b.Debug {
		log.SetLevel(log.DebugLevel)
	}
	timestamp := makeEdgeTimeStamp()
	nonce := createNonce()

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", b.createAuthHeader(req, timestamp, nonce))
	return req
}

// InitConfig initializes configuration file
func InitConfig(file string) Base {
	var b Base
	viper.SetConfigType("yaml")
	viper.SetConfigFile(file)
	err := viper.ReadInConfig()
	if err != nil {
		log.Panicf("Fatal error config file: %s \n", err)
	}
	err = viper.Unmarshal(&b)
	if err != nil {
		log.Panicf("Fatal error unable to unmarshal config file: %s \n", err)
	}
	return b
}
