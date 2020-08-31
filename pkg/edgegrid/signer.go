package edgegrid

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

// SignRequest adds a signed authorization header to the http request
func (c Config) SignRequest(r *http.Request) error {
	timestamp := Timestamp(time.Now())

	nonce, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	authHeader := fmt.Sprintf("EG1-HMAC-SHA256 client_token=%s;access_token=%s;timestamp=%s;nonce=%s;",
		c.ClientToken,
		c.AccessToken,
		timestamp,
		nonce,
	)

	msgPath := r.URL.EscapedPath()
	if r.URL.RawQuery != "" {
		msgPath = fmt.Sprintf("%s?%s", msgPath, r.URL.RawQuery)
	}

	// create the message to be signed
	msgData := []string{
		r.Method,
		r.URL.Scheme,
		r.URL.Host,
		msgPath,
		c.canonicalizeHeaders(r),
		c.createContentHash(r),
		authHeader,
	}
	msg := strings.Join(msgData, "\t")

	key := createSignature(timestamp, c.ClientSecret)

	r.Header.Set("Authorization", fmt.Sprintf("%s=%s", authHeader, createSignature(msg, key)))

	return nil
}

func (c Config) canonicalizeHeaders(r *http.Request) string {
	var unsortedHeader []string
	var sortedHeader []string
	for k := range r.Header {
		unsortedHeader = append(unsortedHeader, k)
	}
	sort.Strings(unsortedHeader)
	for _, k := range unsortedHeader {
		for _, sign := range c.HeaderToSign {
			if sign == k {
				v := strings.TrimSpace(r.Header.Get(k))
				sortedHeader = append(sortedHeader, fmt.Sprintf("%s:%s", strings.ToLower(k), strings.ToLower(stringMinifier(v))))
			}
		}
	}
	return strings.Join(sortedHeader, "\t")
}

// The content hash is the base64-encoded SHA–256 hash of the POST body.
// For any other request methods, this field is empty. But the tab separator (\t) must be included.
// The size of the POST body must be less than or equal to the value specified by the service.
// Any request that does not meet this criteria SHOULD be rejected during the signing process,
// as the request will be rejected by EdgeGrid.
func (c Config) createContentHash(r *http.Request) string {
	var (
		contentHash  string
		preparedBody string
		bodyBytes    []byte
	)

	if r.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(r.Body)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		preparedBody = string(bodyBytes)
	}

	if r.Method == "POST" && len(preparedBody) > 0 {
		if len(preparedBody) > c.MaxBody {
			preparedBody = preparedBody[0:c.MaxBody]
		}

		sum := sha256.Sum256([]byte(preparedBody))

		contentHash = base64.StdEncoding.EncodeToString(sum[:])
	}

	return contentHash
}
