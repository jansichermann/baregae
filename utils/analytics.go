package utils

import (
	"appengine"
	"appengine/delay"
	"appengine/urlfetch"
	"net/http"
	"net/url"
	"strings"
)

const gaEndpoint = "https://ssl.google-analytics.com/collect"

func GoogleAnalyticsPageViewMap(projectId, clientId, userId, userAgent, referer, documentPath string) map[string]string {
	params := make(map[string]string)
	params["tid"] = projectId
	params["cid"] = clientId
	params["uid"] = userId
	params["ua"] = userAgent
	params["dr"] = referer
	params["dp"] = documentPath
	params["t"] = "pageview" // type

	return params
}

func SendGoogleAnalyticsTrackingData(ctx appengine.Context, params map[string]string) {
	client := urlfetch.Client(ctx)

	values := url.Values{}
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

var delayedTracking = delay.Func("utils.analytics.delayedTrack", SendGoogleAnalyticsTrackingData)

func TrackPage(r *http.Request, clientId, userId, projectId string) {
	params := GoogleAnalyticsPageViewMap(projectId, clientId, userId, r.UserAgent(), r.Referer(), r.URL.Path)
	SendGoogleAnalyticsTrackingData(appengine.NewContext(r), params)
}

func TrackPageDelayed(r *http.Request, clientId, userId, projectId string) {
	params := GoogleAnalyticsPageViewMap(projectId, clientId, userId, r.UserAgent(), r.Referer(), r.URL.Path)
	delayedTracking.Call(appengine.NewContext(r), params)
}
