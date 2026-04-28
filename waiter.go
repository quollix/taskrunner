package taskrunner

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"
)

func (t *TaskRunner) WaitUntilPortIsReady(port string) {
	t.retryOperation(func() (bool, error) {
		conn, err := net.DialTimeout("tcp", "localhost:"+port, 1*time.Second)
		if err == nil {
			conn.Close()
			return true, nil
		}
		return false, err
	}, "port", "localhost:"+port, t.Config.DefaultWaitTimeout)
}

func (t *TaskRunner) retryOperation(operation func() (bool, error), description, target string, maxAttempts int) {
	attempt := 0
	for attempt < maxAttempts {
		success, err := operation()
		if success && err == nil {
			t.Log.Info("%s was requested successfully at %s", description, target)
			return
		} else {
			if attempt%5 == 0 {
				t.Log.Info("attempt %v/%v: %s is not yet reachable at %s. error: %v. Trying again...", attempt, maxAttempts, description, target, err)
			}
			attempt++
			time.Sleep(1 * time.Second)
		}
	}
	t.Log.Error("error: %s could not be reached in time at %s. Cleanup and exit...", description, target)
	t.ExitWithError()
}

func (t *TaskRunner) WaitForWebPageToBeReady(targetUrl string) {
	t.retryOperation(func() (bool, error) {
		parsedUrl, err := url.Parse(targetUrl)
		if err != nil {
			return false, err
		}

		request, err := http.NewRequest(http.MethodGet, targetUrl, nil)
		if err != nil {
			return false, err
		}
		request.Header.Set("Origin", parsedUrl.Scheme+"://"+parsedUrl.Host)

		httpClient := http.DefaultClient
		if parsedUrl.Scheme == "https" {
			httpClient = &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
					},
				},
			}
		}

		response, err := httpClient.Do(request)
		if err == nil && response.StatusCode == http.StatusOK {
			return true, nil
		}
		return false, err
	}, "Index page", targetUrl, t.Config.DefaultWaitTimeout)
}
