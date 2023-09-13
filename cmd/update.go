package cmd

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"runtime"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/minio/selfupdate"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

const (
	releaseURL = "https://api.github.com/repos/autoproof/cli/releases/latest"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates Autoproof CLI to the latest version",
	RunE: func(cmd *cobra.Command, args []string) error {
		currentVersion := "v" + cliVersion
		if !semver.IsValid(currentVersion) {
			return fmt.Errorf("current application version is not semver compatible: %s", currentVersion)
		}

		rt := getUpdateTransport(30 * time.Second)
		latestUpstreamRelease, err := getLatestRelease(rt)
		if err != nil {
			return err
		}

		if !semver.IsValid(latestUpstreamRelease) {
			return fmt.Errorf("remote version is not semver compatible: %s", latestUpstreamRelease)
		}

		if semver.Compare(currentVersion, latestUpstreamRelease) > 0 {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Autoproof CLI is up-to-date.")
			return nil
		}

		cliBinURL := fmt.Sprintf(
			"https://github.com/autoproof/cli/releases/download/%s/autoproofcli_%s_%s.tar.gz",
			latestUpstreamRelease,
			runtime.GOOS,
			runtime.GOARCH,
		)

		reader, length, err := getUpdateReaderFromURL(cliBinURL, rt)
		defer reader.Close()

		pbar := pb.New64(length).
			SetWriter(cmd.OutOrStdout()).
			SetTemplateString(`{{ red "Downloading:" }} {{bar . (red "[") (green "=") (red "]")}} {{speed . | rndcolor }}`).
			Start()
		defer pbar.Finish()

		gzipReader, err := gzip.NewReader(pbar.NewProxyReader(reader))
		if err != nil {
			return err
		}
		defer gzipReader.Close()

		tarReader := tar.NewReader(gzipReader)
		for {
			header, err := tarReader.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}

			if header.Name == "autoproofcli" {
				if err := selfupdate.Apply(tarReader, selfupdate.Options{}); err != nil {
					if rollbackErr := selfupdate.RollbackError(err); rollbackErr != nil {
						return fmt.Errorf("self update rollaback error: %w", rollbackErr)
					}
					return fmt.Errorf("self update: %w", err)
				}

				pbar.Finish()
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Update Autoproof CLI to latest release %s\n", latestUpstreamRelease)
				return nil
			}
		}

		return nil
	},
}

func getUpdateReaderFromURL(url string, rt http.RoundTripper) (io.ReadCloser, int64, error) {
	httpClient := http.Client{
		Transport: rt,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, -1, fmt.Errorf("new request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, -1, fmt.Errorf("request release endpoint: %w", err)
	}

	return resp.Body, resp.ContentLength, nil
}

func getUpdateTransport(timeout time.Duration) http.RoundTripper {
	var updateTransport http.RoundTripper = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: timeout,
		}).DialContext,
		IdleConnTimeout:       timeout,
		TLSHandshakeTimeout:   timeout,
		ExpectContinueTimeout: timeout,
		DisableCompression:    true,
	}
	return updateTransport
}

func getLatestRelease(rt http.RoundTripper) (string, error) {
	httpClient := http.Client{
		Transport: rt,
	}

	req, err := http.NewRequest(http.MethodGet, releaseURL, nil)
	if err != nil {
		return "", err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request release endpoint: %w", err)
	}
	defer resp.Body.Close()

	var releaseResponse struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&releaseResponse); err != nil {
		return "", fmt.Errorf("decode release response: %w", err)
	}

	return releaseResponse.TagName, nil
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
