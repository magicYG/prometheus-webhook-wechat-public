package webrouter

import (
	"encoding/json"
	"log"
	"net/http"

	"git.ucloudadmin.com/uk8s/prometheus-webhook-wechat-public/pkg/models"
	"git.ucloudadmin.com/uk8s/prometheus-webhook-wechat-public/pkg/notifier"
	"git.ucloudadmin.com/uk8s/prometheus-webhook-wechat-public/pkg/request"
	"github.com/go-chi/chi"
)

type WechatResource struct {
	Profileurl string
	HttpClient *http.Client
	Chatids    map[string][]string
	Corpid     string
	Corpsecret string
	TemplateID string
}

func (rs *WechatResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/{profile}/send", rs.SendNotification)
	return r
}

func (rs *WechatResource) SendNotification(w http.ResponseWriter, r *http.Request) {
	profile := chi.URLParam(r, "profile")
	getTokenResp, err := request.Get(rs.Corpid, rs.Corpsecret, rs.Profileurl)
	if err != nil {
		log.Printf("Failed to request token: %s", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	webhookURL := "https://" + rs.Profileurl + "/cgi-bin/message/template/send?access_token=" + getTokenResp

	var promMessage models.WebhookMessage
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&promMessage); err != nil {
		log.Printf("Cannot decode prometheus webhook JSON request: %s", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	for _, alert := range promMessage.Alerts {
		for _, chatid := range rs.Chatids[profile] {
			notification := notifier.BuildWechatMsg(rs.TemplateID, &promMessage, alert, chatid)
			robotResp, err := notifier.SendWechatNotification(rs.HttpClient, webhookURL, notification)
			if err != nil {
				log.Printf("Failed to send notification to %s: %s", chatid, err)
				http.Error(w, "Bad Request", http.StatusBadRequest)
				continue
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("send to " + chatid + " OK\n"))
			if robotResp.ErrorCode != 0 {
				log.Printf("Failed to send notification to wechat: [%d] %s", robotResp.ErrorCode, robotResp.ErrorMessage)
				continue
			}
			log.Printf("Successfully send notification to %s", chatid)
		}
	}

}
