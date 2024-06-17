package models

type WechatNotificationResponse struct {
	ErrorMessage string `json:"errmsg"`
	ErrorCode    int    `json:"errcode"`
}

type WechatAlarmTemplate struct {
	Touser     string          `json:"touser"`
	URL        string          `json:"url"`
	TemplateID string          `json:"template_id"`
	Data       WechatAlarmData `json:"data"`
}

type WechatAlarmData struct {
	//告警状态
	Status ValueColor `json:"status"`
	//告警主题
	Summary ValueColor `json:"summary"`
	//告警类型
	Alertname ValueColor `json:"alertname"`
	//告警级别
	Severity ValueColor `json:"severity"`
	//告警实例
	Instance ValueColor `json:"instance"`
	//告警内容
	Message ValueColor `json:"message"`
	//告警时间
	StartsAt ValueColor `json:"startsat"`
	//告警链接
	GeneratorURL ValueColor `json:"generatorurl"`
}

type ValueColor struct {
	Value string `json:"value"`
	Color string `json:"color"`
}
