package proxy

import (
	"context"
	"database/sql"
	"net/http"
	"net/url"

	"golang.org/x/net/http/httpproxy"

	"yunion.io/x/jsonutils"

	proxyapi "yunion.io/x/onecloud/pkg/apis/cloudcommon/proxy"
	"yunion.io/x/onecloud/pkg/cloudcommon/db"
	"yunion.io/x/onecloud/pkg/httperrors"
	"yunion.io/x/onecloud/pkg/mcclient"
	"yunion.io/x/onecloud/pkg/util/httputils"
)

type SProxySettingManager struct {
	db.SStandaloneResourceBaseManager
	db.SDomainizedResourceBaseManager
}

var ProxySettingManager *SProxySettingManager

func init() {
	ProxySettingManager = &SProxySettingManager{
		SStandaloneResourceBaseManager: db.NewStandaloneResourceBaseManager(
			SProxySetting{},
			"proxysettings_tbl",
			"proxysetting",
			"proxysettings",
		),
	}
	ProxySettingManager.SetVirtualObject(ProxySettingManager)
}

type SProxySetting struct {
	db.SStandaloneResourceBase
	db.SDomainizedResourceBase

	HTTPProxy  string `create:"admin_optional" list:"admin" update:"admin"`
	HTTPSProxy string `create:"admin_optional" list:"admin" update:"admin"`
	NoProxy    string `create:"admin_optional" list:"admin" update:"admin"`
}

func (man *SProxySettingManager) ValidateCreateData(ctx context.Context, userCred mcclient.TokenCredential, ownerId mcclient.IIdentityProvider, query jsonutils.JSONObject, data proxyapi.ProxySettingCreateInput) (proxyapi.ProxySettingCreateInput, error) {
	if err := data.ProxySetting.Sanitize(); err != nil {
		return data, httperrors.NewInputParameterError("%s", err)
	}
	return data, nil
}

func (ps *SProxySetting) CustomizeCreate(ctx context.Context, userCred mcclient.TokenCredential, ownerId mcclient.IIdentityProvider, query jsonutils.JSONObject, data jsonutils.JSONObject) error {
	ps.DomainId = ownerId.GetProjectDomainId()
	return ps.SStandaloneResourceBase.CustomizeCreate(ctx, userCred, ownerId, query, data)
}

func (ps *SProxySetting) ValidateUpdateData(ctx context.Context, userCred mcclient.TokenCredential, query jsonutils.JSONObject, data proxyapi.ProxySettingUpdateInput) (proxyapi.ProxySettingUpdateInput, error) {
	if ps.Id == proxyapi.ProxySettingId_DIRECT {
		return data, httperrors.NewConflictError("DIRECT setting cannot be changed")
	}
	if err := data.ProxySetting.Sanitize(); err != nil {
		return data, httperrors.NewInputParameterError("%s", err)
	}
	return data, nil
}

func (ps *SProxySetting) HttpTransportProxyFunc() httputils.TransportProxyFunc {
	cfg := &httpproxy.Config{
		HTTPProxy:  ps.HTTPProxy,
		HTTPSProxy: ps.HTTPSProxy,
		NoProxy:    ps.NoProxy,
	}
	proxyFunc := cfg.ProxyFunc()
	return func(req *http.Request) (*url.URL, error) {
		return proxyFunc(req.URL)
	}
}

func (ps *SProxySetting) ValidateDeleteCondition(ctx context.Context) error {
	if ps.Id == proxyapi.ProxySettingId_DIRECT {
		return httperrors.NewConflictError("DIRECT setting cannot be deleted")
	}
	for _, man := range referrersMen {
		t := man.TableSpec().Instance()
		n, err := t.Query().
			Equals("proxy_setting_id", ps.Id).
			CountWithError()
		if err != nil {
			return httperrors.NewInternalServerError("get proxysetting refcount fail %s", err)
		}
		if n > 0 {
			return httperrors.NewResourceBusyError("proxysetting %s is still referred to by %d %s",
				ps.Id, n, man.KeywordPlural())
		}
	}
	return nil
}

func (man *SProxySettingManager) InitializeData() error {
	_, err := man.FetchById(proxyapi.ProxySettingId_DIRECT)
	if err == nil {
		return nil
	}
	if err != sql.ErrNoRows {
		return err
	}

	m, err := db.NewModelObject(man)
	if err != nil {
		return err
	}
	ps := m.(*SProxySetting)
	ps.Id = proxyapi.ProxySettingId_DIRECT
	ps.Name = proxyapi.ProxySettingId_DIRECT
	ps.Description = "Connect directly"
	if err := man.TableSpec().Insert(ps); err != nil {
		return err
	}
	return nil
}

var referrersMen []db.IModelManager

func RegisterReferrer(man db.IModelManager) {
	referrersMen = append(referrersMen, man)
}