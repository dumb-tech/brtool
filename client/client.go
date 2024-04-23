package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
)

var (
	apiRequestURL = func(useTLS bool, host, endpoint string) string {
		scheme := "http://"
		if useTLS {
			scheme = "https://"
		}
		return scheme + host + endpoint
	}

	endpointApi = "/api"

	endpointApiSettings       = endpointApi + "/settings/"
	endpointApiDatabase       = endpointApi + "/database"
	endpointDatabaseRows      = endpointApiDatabase + "/rows"
	endpointDatabaseRowsTable = endpointDatabaseRows + "/table"

	endpointUpdateRow = func(tableID, rowID string) string {
		return fmt.Sprintf(endpointDatabaseRowsTable+"/%s/%s/?user_field_names=true", tableID, rowID)
	}
)

type config struct {
	host  string
	token string
	debug bool
}

type BaserowClient struct {
	cl     *http.Client
	cfg    config
	log    *slog.Logger
	useTLS bool
}

func New(host string, token string) *BaserowClient {
	return &BaserowClient{
		cl:  &http.Client{},
		cfg: config{host: host, token: token},
		log: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
	}
}

func (bc *BaserowClient) SetLogger(log *slog.Logger) *BaserowClient {
	if log != nil {
		bc.log = log
	}
	return bc
}

func (bc *BaserowClient) Debug(v bool) *BaserowClient {
	bc.cfg.debug = v
	return bc
}

func (bc *BaserowClient) UseTLS(v bool) *BaserowClient {
	bc.cl.Transport = http.DefaultTransport
	if v {
		bc.cl.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	bc.useTLS = v

	return bc
}

func (bc *BaserowClient) Ping() error {
	req, err := http.NewRequest(http.MethodGet, apiRequestURL(bc.useTLS, bc.cfg.host, ""), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Token "+bc.cfg.token)

	if bc.cfg.debug {
		dump, _ := httputil.DumpRequest(req, true)
		fmt.Printf("=== REQUEST =================\n\n %s\n\n", string(dump))
		fmt.Println()
	}

	resp, err := bc.cl.Get(apiRequestURL(bc.useTLS, bc.cfg.host, endpointApiSettings))
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("failed to close body via ping baserow host")
		}
	}(resp.Body)

	if bc.cfg.debug {
		dump, _ := httputil.DumpResponse(resp, true)
		fmt.Printf("=== RESPONSE =================\n\n %s\n\n", string(dump))
		fmt.Println()
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to ping baserows host %q", bc.cfg.host)
	}

	type Settings struct {
		AllowNewSignups                     bool        `json:"allow_new_signups"`
		AllowSignupsViaWorkspaceInvitations bool        `json:"allow_signups_via_workspace_invitations"`
		AllowSignupsViaGroupInvitations     bool        `json:"allow_signups_via_group_invitations"`
		AllowResetPassword                  bool        `json:"allow_reset_password"`
		AllowGlobalWorkspaceCreation        bool        `json:"allow_global_workspace_creation"`
		AllowGlobalGroupCreation            bool        `json:"allow_global_group_creation"`
		AccountDeletionGraceDelay           int         `json:"account_deletion_grace_delay"`
		ShowAdminSignupPage                 bool        `json:"show_admin_signup_page"`
		TrackWorkspaceUsage                 bool        `json:"track_workspace_usage"`
		ShowBaserowHelpRequest              bool        `json:"show_baserow_help_request"`
		CoBrandingLogo                      interface{} `json:"co_branding_logo"`
		InstanceWideLicenses                struct {
		} `json:"instance_wide_licenses"`
	}

	var settings Settings

	if err := json.NewDecoder(resp.Body).Decode(&settings); err != nil {
		return fmt.Errorf("failed to get settings by baserow api via host %q", bc.cfg.host)
	}

	return nil
}

func (bc *BaserowClient) UpdateRowField(tableID int, rowID int, field string, new string) error {
	url := apiRequestURL(bc.useTLS, bc.cfg.host, endpointUpdateRow(strconv.Itoa(tableID), strconv.Itoa(rowID)))

	payload, err := json.Marshal(map[string]string{field: new})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Token "+bc.cfg.token)
	req.Header.Set("Content-Type", "application/json")

	if bc.cfg.debug {
		dump, _ := httputil.DumpRequest(req, true)
		fmt.Printf("=== REQUEST =================\n\n %s\n\n", string(dump))
		fmt.Println()
	}

	resp, err := bc.cl.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("failed to close response body via update row field")
		}
	}(resp.Body)

	if bc.cfg.debug {
		dump, _ := httputil.DumpResponse(resp, true)
		fmt.Printf("=== RESPONSE =================\n\n %s\n\n", string(dump))
		fmt.Println()
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update baserow row field. status - %d\n",
			resp.StatusCode)
	}

	return nil
}
