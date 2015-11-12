package authware

import (
	"encoding/json"
	"errors"
	"github.com/InteractiveLecture/serviceclient"
	"github.com/InteractiveLecture/serviceclient/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func mockOptions(expectedPermissions []string, realPermissions []string) Options {
	emptyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	result := DefaultOptions(emptyHandler, "media", expectedPermissions...)
	//Mock extractor and fetcher
	result.Extractor = func(r *http.Request) (id string, sid string) {
		return "1", "admin"
	}
	result.Fetcher = func(id string, sid string, objectClass string) (map[string]interface{}, error) {
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
	mock := mocks.NewMockBackendAdapter(mockCtrl)
	defer mockCtrl.Finish()
	authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		permissions := map[string]interface{}{
			"read":   true,
			"write":  true,
			"delete": false,
		}
		encoder := json.NewEncoder(w)
		encoder.Encode(permissions)
	}))
	defer authServer.Close()
	mock.EXPECT().Configure("acl-service").Return(nil)
	url := strings.Split(authServer.URL, "http://")[1]
	mock.EXPECT().Resolve("acl-service").Return(url, nil)
	//mock.EXPECT().Refresh()
	serviceclient.Configure(mock, "acl-service")
	emptyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	options := DefaultOptions(emptyHandler, "media", "read")
	permissions, err := options.Fetcher("1", "admin", "media")
	assert.Nil(t, err)
	assert.True(t, permissions["read"].(bool))
	assert.True(t, permissions["write"].(bool))
	assert.False(t, permissions["delete"].(bool))

	//authServer := httptest.NewServer()

}

func pa(permissions ...string) []string {
	return permissions
}
