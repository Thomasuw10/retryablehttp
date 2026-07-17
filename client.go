package retryablehttp

import (
	"net/http"
)

// ... (existing code)

func (c *Client) Do(req *Request) (*http.Response, error) {
	// ... (existing code)
	for i := 0; i <= c.RetryMax; i++ {
		// Create a new request for this attempt
		newReq := req.clone()
		
		// Apply PrepareRetry hook if present
		if c.PrepareRetry != nil {
			c.PrepareRetry(newReq)
		}

		// Ensure headers are set correctly to avoid accumulation
		// The clone() method should ideally copy the map, but we ensure 
		// that modifications via PrepareRetry don't append to existing slices.
		// If the user uses req.Header.Set in PrepareRetry, it is safe.
		
		resp, err := c.HTTPClient.Do(newReq.Request)
		// ... (rest of the loop)
	}
	return nil, nil
}

// Ensure Request.clone() performs a deep copy of the Header map
func (r *Request) clone() *Request {
	np := *r
	np.Request = r.Request.Clone(r.Request.Context())
	// Deep copy headers
	np.Header = make(http.Header)
	for k, v := range r.Header {
		np.Header[k] = append([]string(nil), v...)
	}
	return &np
}