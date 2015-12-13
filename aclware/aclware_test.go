package authware

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/InteractiveLecture/serviceclient/test"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func mockOptions(expectedPermissions []string, realPermissions []string) Options {
	emptyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	result := DefaultOptions(emptyHandler, "media", expectedPermissions...)
	//Mock extractor and fetcher
	result.Extractor = func(r *http.Request) (id string, sid string) {
		return "1", "admin"
	}
	result.Fetcher = func(id string, sid string) (map[string]interface{}, error) {
		var permissions = make(map[string]interface{})
		if realPermissions == nil {
			return nil, errors.New("mock error")
		}
		for _, v := range realPermissions {
			permissions[v] = true
		}
		return permissions, nil
	}
	return result
}

func TestPermissions(t *testing.T) {
	permissions := pa("read", "write", "delete", "create")
	handler := New(mockOptions(permissions, permissions))
	mediaServer := httptest.NewServer(handler)
	defer mediaServer.Close()
	resp, err := http.Get(mediaServer.URL)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestFailingPermissions(t *testing.T) {
	handler := New(mockOptions(pa("read", "write"), pa("read")))
	mediaServer := httptest.NewServer(handler)
	defer mediaServer.Close()
	resp, err := http.Get(mediaServer.URL)
	assert.Nil(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 401, resp.StatusCode)
}

func TestFailGettingPermissions(t *testing.T) {
	handler := New(mockOptions(pa("read", "write"), nil))
	mediaServer := httptest.NewServer(handler)
	defer mediaServer.Close()
	resp, err := http.Get(mediaServer.URL)
	assert.Nil(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 500, resp.StatusCode)
}

func TestDefaultFetcher(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	authServer := servicetest.Service(mockCtrl, "acl-service", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		permissions := map[string]interface{}{
			"read":   true,
			"create": true,
			"delete": false,
			"update": false,
		}
		encoder := json.NewEncoder(w)
		encoder.Encode(permissions)
	}))
	defer mockCtrl.Finish()
	defer authServer.Close()
	emptyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	options := DefaultOptions(emptyHandler, "media", "read")
	permissions, err := options.Fetcher("1", "admin")
	assert.Nil(t, err)
	assert.True(t, permissions["read"].(bool))
	assert.True(t, permissions["create"].(bool))
	assert.False(t, permissions["delete"].(bool))
}

func pa(permissions ...string) []string {
	return permissions
}
