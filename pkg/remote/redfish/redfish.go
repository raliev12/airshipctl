package redfish

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"

	redfishApi "opendev.org/airship/go-redfish/api"
	redfishClient "opendev.org/airship/go-redfish/client"

	alog "opendev.org/airship/airshipctl/pkg/log"
)

type RemoteDirect struct {

	// Context
	Context context.Context

	// remote URL
	RemoteURL url.URL

	// ephemeral Host ID
	EphemeralNodeID string

	// ISO URL
	IsoPath string

	// Redfish Client implementation
	RedfishAPI redfishApi.RedfishAPI

	// optional Username Authentication
	Username string

	// optional Password
	Password string
}

// Top level function to handle Redfish remote direct
func (cfg RemoteDirect) DoRemoteDirect() error {
	alog.Debugf("Using Redfish Endpoint: '%s'", cfg.RemoteURL.String())

	/* Get system details */
	systemID := cfg.EphemeralNodeID
	system, _, err := cfg.RedfishAPI.GetSystem(cfg.Context, systemID)
	if err != nil {
		return ErrRedfishClient{Message: fmt.Sprintf("Get System[%s] failed with err: %v", systemID, err)}
	}
	alog.Debugf("Ephemeral Node System ID: '%s'", systemID)

	/* get manager for system */
	managerID := GetResourceIDFromURL(system.Links.ManagedBy[0].OdataId)
	alog.Debugf("Ephemeral node managerID: '%s'", managerID)

	/* Get manager's Cd or DVD virtual media ID */
	vMediaID, vMediaType, err := GetVirtualMediaID(cfg.Context, cfg.RedfishAPI, managerID)
	if err != nil {
		return err
	}
	alog.Debugf("Ephemeral Node Virtual Media Id: '%s'", vMediaID)

	/* Load ISO in manager's virtual media */
	err = SetVirtualMedia(cfg.Context, cfg.RedfishAPI, managerID, vMediaID, cfg.IsoPath)
	if err != nil {
		return err
	}
	alog.Debugf("Successfully loaded virtual media: '%s'", cfg.IsoPath)

	/* Set system's bootsource to selected media */
	err = SetSystemBootSourceForMediaType(cfg.Context, cfg.RedfishAPI, systemID, vMediaType)
	if err != nil {
		return err
	}

	/* Reboot system */
	err = RebootSystem(cfg.Context, cfg.RedfishAPI, systemID)
	if err != nil {
		return err
	}
	alog.Debug("Restarted ephemeral host")

	return nil
}

// NewRedfishRemoteDirectClient creates a new Redfish remote direct client.
func NewRedfishRemoteDirectClient(
	remoteURL string,
	ephNodeID string,
	username string,
	password string,
	isoPath string,
	insecure bool,
	useproxy bool,
) (RemoteDirect, error) {
	if remoteURL == "" {
		return RemoteDirect{},
			ErrRedfishMissingConfig{
				What: "redfish remote url empty",
			}
	}

	if ephNodeID == "" {
		return RemoteDirect{},
			ErrRedfishMissingConfig{
				What: "redfish ephemeral node id empty",
			}
	}

	var ctx context.Context
	if username != "" && password != "" {
		ctx = context.WithValue(
			context.Background(),
			redfishClient.ContextBasicAuth,
			redfishClient.BasicAuth{UserName: username, Password: password},
		)
	} else {
		ctx = context.Background()
	}

	if isoPath == "" {
		return RemoteDirect{},
			ErrRedfishMissingConfig{
				What: "redfish ephemeral node iso Path empty",
			}
	}

	cfg := &redfishClient.Configuration{
		BasePath:      remoteURL,
		DefaultHeader: make(map[string]string),
		UserAgent:     "airshipctl/client",
	}

	// see https://github.com/golang/go/issues/26013
	// We clone the default transport to ensure when we customize the transport
	// that we are providing it sane timeouts and other defaults that we would
	// normally get when not overriding the transport
	defaultTransportCopy := (http.DefaultTransport.(*http.Transport))
	transport := defaultTransportCopy.Clone()

	if insecure {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true, //nolint:gosec
		}
	}
	if !useproxy {
		transport.Proxy = nil
	}

	cfg.HTTPClient = &http.Client{
		Transport: transport,
	}

	var api redfishApi.RedfishAPI = redfishClient.NewAPIClient(cfg).DefaultApi

	parsedURL, err := url.Parse(remoteURL)
	if err != nil {
		return RemoteDirect{},
			ErrRedfishMissingConfig{
				What: fmt.Sprintf("invalid url format: %v", err),
			}
	}

	client := RemoteDirect{
		Context:         ctx,
		RemoteURL:       *parsedURL,
		EphemeralNodeID: ephNodeID,
		IsoPath:         isoPath,
		RedfishAPI:      api,
	}

	return client, nil
}
