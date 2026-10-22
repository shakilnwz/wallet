package sync

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/shakilnwz/wallet/core/models"
	"github.com/shakilnwz/wallet/db"
)

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func normalizeAddress(addr string) string {
	addr = strings.TrimSpace(addr)
	if !strings.HasPrefix(addr, "http://") && !strings.HasPrefix(addr, "https://") {
		addr = "http://" + addr
	}
	return strings.TrimRight(addr, "/")
}

// errPeerOutdated is returned when the peer doesn't serve the current sync API.
var errPeerOutdated = errors.New("peer is running an older Wallet version; update both devices")

func checkStatus(resp *http.Response) error {
	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return errPeerOutdated
	default:
		return fmt.Errorf("peer returned HTTP %d", resp.StatusCode)
	}
}

func (c *Client) Ping(peerAddr string) (string, error) {
	url := normalizeAddress(peerAddr) + apiPrefix + "/ping"
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if err := checkStatus(resp); err != nil {
		return "", err
	}

	var data struct {
		DeviceID string `json:"device_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	return data.DeviceID, nil
}

// FullTwoWaySync pulls remote records and merges into local DB, then pushes local records to remote peer.
func (c *Client) FullTwoWaySync(peerAddr string, database *db.DB) (pulledIns, pulledUpd, pushedIns, pushedUpd int, err error) {
	base := normalizeAddress(peerAddr)

	// 1. PULL: Fetch remote package
	pullResp, err := c.httpClient.Get(base + apiPrefix + "/package")
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("failed to fetch remote package: %w", err)
	}
	defer pullResp.Body.Close()

	if err := checkStatus(pullResp); err != nil {
		return 0, 0, 0, 0, fmt.Errorf("remote pull: %w", err)
	}

	var remotePkg models.SyncPackage
	if err := json.NewDecoder(pullResp.Body).Decode(&remotePkg); err != nil {
		return 0, 0, 0, 0, fmt.Errorf("failed to decode remote package: %w", err)
	}

	// Merge remote package into local DB
	pulledIns, pulledUpd, _, err = database.ImportAndMergeSyncPackage(&remotePkg)
	if err != nil {
		return pulledIns, pulledUpd, 0, 0, fmt.Errorf("failed to merge remote package locally: %w", err)
	}

	// 2. PUSH: Export merged local package and send to remote peer
	localPkg, err := database.ExportSyncPackage()
	if err != nil {
		return pulledIns, pulledUpd, 0, 0, fmt.Errorf("failed to export local package: %w", err)
	}

	payload, err := json.Marshal(localPkg)
	if err != nil {
		return pulledIns, pulledUpd, 0, 0, err
	}

	pushResp, err := c.httpClient.Post(base+apiPrefix+"/merge", "application/json", bytes.NewReader(payload))
	if err != nil {
		return pulledIns, pulledUpd, 0, 0, fmt.Errorf("failed to push to remote: %w", err)
	}
	defer pushResp.Body.Close()

	if err := checkStatus(pushResp); err != nil {
		return pulledIns, pulledUpd, 0, 0, fmt.Errorf("remote merge: %w", err)
	}

	var pushResult struct {
		Inserted int `json:"inserted"`
		Updated  int `json:"updated"`
		Skipped  int `json:"skipped"`
	}
	_ = json.NewDecoder(pushResp.Body).Decode(&pushResult)

	return pulledIns, pulledUpd, pushResult.Inserted, pushResult.Updated, nil
}
