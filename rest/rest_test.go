package rest_test

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"
)

const baseURL = "https://localhost:8443/api/v1"

type TestClient struct {
	client *http.Client
	token  string
}

func NewTestClient(token string) *TestClient {
	return &TestClient{
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
		token: token,
	}
}

func (c *TestClient) Do(method, path string, body url.Values) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewBufferString(body.Encode())
	}
	req, _ := http.NewRequest(method, baseURL+path, reqBody)
	req.Header.Set("Authorization", "Bearer "+c.token)
	if method == "POST" || method == "PUT" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	return c.client.Do(req)

}

func TestMemberLifecycle(t *testing.T) {
	jwt := os.Getenv("TEST_ADMIN_JWT")
	if jwt == "" {
		t.Fatalf("TEST_ADMIN_JWT environment not set")
	}
	adminClient := NewTestClient(jwt)
	var memberID int

	// 1. Create Member (Admin)
	t.Run("CreateMember", func(t *testing.T) {
		form := url.Values{}
		form.Add("firstname", "Test")
		form.Add("lastname", "User")
		resp, err := adminClient.Do("POST", "/member", form)

		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("Failed to create member: %v", resp.Status)
		}

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		memberID = int(res["id"].(float64))
	})

	// 2. Set Field (Manager Level via Admin Client)
	t.Run("ManagerSetField", func(t *testing.T) {
		form := url.Values{}
		form.Add("field", "Occupation")
		form.Add("value", "Refactorer")

		path := fmt.Sprintf("/member/%d", memberID)
		resp, _ := adminClient.Do("PUT", path, form)

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Manager update failed: %v", resp.Status)
		}
	})

	// 3. Self Update (requires getting a JWT for the new user)
	// Assuming a helper to generate tokens for tests
	jwt = os.Getenv("TEST_USER_JWT")
	if jwt == "" {
		t.Fatalf("TEST_USER_JWT environment not set")
	}
	userClient := NewTestClient(jwt)

	t.Run("SelfUpdateSuccess", func(t *testing.T) {
		form := url.Values{}
		form.Add("field", "PreferredName")
		form.Add("value", "Tester")

		resp, _ := userClient.Do("PUT", "/me", form)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Self update failed: %v", resp.Status)
		}
	})

	t.Run("SelfUpdateInvalidField", func(t *testing.T) {
		form := url.Values{}
		form.Add("field", "id") // Restricted in SetMeField
		form.Add("value", "999")
		resp, _ := userClient.Do("PUT", "/me", form)
		if resp.StatusCode == http.StatusOK {
			t.Error("Should have failed to update ID field")
		}
	})

	t.Run("SelfUpdateSQLInjection", func(t *testing.T) {
		form := url.Values{}
		form.Add("field", "Title`; DROP TABLE member;--")
		form.Add("value", "Chaos")
		resp, _ := userClient.Do("PUT", "/me", form)
		if resp.StatusCode != http.StatusNotAcceptable && resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("SQL injection check failed, got status: %v", resp.Status)
		}
	})

	// 4. Cleanup/Delete (if endpoint exists, or via direct DB if needed for test)
	// Note: Your router currently has delete commented out for several items.
}
