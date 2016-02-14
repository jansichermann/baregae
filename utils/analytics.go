package utils

import (
	"appengine"
	"appengine/urlfetch"
	"net/http"
	"net/url"
	"strings"
)

const gaEndpoint = "https://ssl.google-analytics.com/collect"

func GoogleAnalyticsPageViewMap(clientId, userId, userAgent, referer, documentPath string) map[string]string {
	params := make(map[string]string)
	params["cid"] = clientId
	params["uid"] = userId
	params["ua"] = userAgent
	params["dr"] = referer
	params["dp"] = documentPath
	params["t"] = "pageview" // type

	return params
}

func SendGoogleAnalyticsTrackingData(ctx appengine.Context, projectId string, params map[string]string) {
	client := urlfetch.Client(ctx)

	values := url.Values{}
	values.Set("tid", projectId)
	for k, v := range params {
		values.Set(k, v)
	}

	resp, err := client.Post(gaEndpoint, "x-www-form-urlencoded", strings.NewReader(values.Encode()))
	if resp.StatusCode != http.StatusOK {
		ctx.Errorf("Tracking Call Returned %s", resp.Status)
	}
	if err != nil {
		ctx.Errorf("%s", err.Error())
		return
	}
}
