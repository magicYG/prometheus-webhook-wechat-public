package notifier

import (
	"bytes"
	"encoding/json"
	"net/http"

	"git.ucloudadmin.com/uk8s/prometheus-webhook-wechat-public/pkg/models"
	"github.com/pkg/errors"
)

const (
	// Template field default color
	wctemplatecolor = "#173177"
)

func BuildWechatMsg(templateid string, promMessage *models.WebhookMessage, alert models.Alert, openID string) *models.WechatAlarmTemplate {
	tpl := &models.WechatAlarmTemplate{
		Touser:     openID,
		URL:        alert.GeneratorURL,
		TemplateID: templateid,
		Data: models.WechatAlarmData{
			Status: models.ValueColor{
				Value: alert.Status,
				Color: wctemplatecolor,
			},
			Summary: models.ValueColor{
				Value: alert.Annotations.Summary,
				Color: wctemplatecolor,
			},
			Alertname: models.ValueColor{
				Value: alert.Labels.Alertname,
				Color: wctemplatecolor,
			},

			Severity: models.ValueColor{
				Value: alert.Labels.Severity,
				Color: wctemplatecolor,
			},
			Instance: models.ValueColor{
				Value: alert.Labels.Instance,
				Color: wctemplatecolor,
			},
			Message: models.ValueColor{
				Value: alert.Annotations.Message,
				Color: wctemplatecolor,
			},
			StartsAt: models.ValueColor{
				Value: alert.StartsAt.Format("2006-01-02 15:04:05"),
				Color: wctemplatecolor,
			},
			GeneratorURL: models.ValueColor{
				Value: alert.GeneratorURL,
				Color: wctemplatecolor,
			},
		},
	}
	return tpl
}

func SendWechatNotification(httpClient *http.Client, ApiURL string, notification *models.WechatAlarmTemplate) (*models.WechatNotificationResponse, error) {
	body, err := json.Marshal(&notification)
	if err != nil {
		return nil, errors.Wrap(err, "error encoding Wechat request")
	}

	httpReq, err := http.NewRequest("POST", ApiURL, bytes.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "error building Wechat request")
	}
	httpReq.Header.Set("Content-Type", "application/json")

	req, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(err, "error sending notification to Wechat")
	}
	defer req.Body.Close()

	if req.StatusCode != 200 {
		return nil, errors.Errorf("unacceptable response code %d", req.StatusCode)
	}

	var robotResp models.WechatNotificationResponse
	enc := json.NewDecoder(req.Body)
	if err := enc.Decode(&robotResp); err != nil {
		return nil, errors.Wrap(err, "error decoding response from DingTalk")
	}

	return &robotResp, nil
}
