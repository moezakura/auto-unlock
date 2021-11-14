package api

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/moezakura/auto-unlock/pkg/config"
	"golang.org/x/xerrors"
)

type Api struct {
	token string
}

func NewApi(config *config.Config) *Api {
	return &Api{
		token: config.UnlockToken,
	}
}

func (a *Api) UnlockWithAll() error {
	u := "https://unlock.mox.si/exec"
	v := url.Values{}
	v.Add("token", "aa")

	resp, err := http.Get(u + v.Encode())
	if err != nil {
		return xerrors.Errorf("failed to get http response: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return xerrors.Errorf("failed to read http response: %w", err)
	}
	return nil
}
