package template

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"testing"
)

type TemplateData struct {
	Template string                 `json:"template"`
	Ph       map[string]interface{} `json:"ph"`
}

// 递归渲染模板
func renderTemplate(tmplStr string, data map[string]interface{}) (string, error) {
	tmpl, err := template.New("tmpl").Parse(tmplStr)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// 递归处理 TemplateData 结构
func renderTemplateData(td *TemplateData) (string, error) {
	// 先处理 ph 里的嵌套模板
	for k, v := range td.Ph {
		// 如果值是 map，且包含 template 和 ph 字段，递归渲染
		if m, ok := v.(map[string]interface{}); ok {
			if tmplStr, ok1 := m["template"].(string); ok1 {
				if phMap, ok2 := m["ph"].(map[string]interface{}); ok2 {
					nestedTD := &TemplateData{
						Template: tmplStr,
						Ph:       phMap,
					}
					rendered, err := renderTemplateData(nestedTD)
					if err != nil {
						return "", err
					}
					td.Ph[k] = template.HTML(rendered)
				}
			}
		}
	}

	// 渲染当前模板
	return renderTemplate(td.Template, td.Ph)
}

func Test2_1(t *testing.T) {
	jsonStr := `
{
    "template": "<div style='display: inline-block'><div> <b>评估项</b>: {{.accessment_item}} </div> </br> <div> <b>评估项描述</b>: {{.accessment_desc}} </div> </br> <div> <b>风险条件</b>: {{.risk_condition}} </div> </br> <div> <b>修复建议</b>: {{.recommendation}}</div></div>",
    "ph": {
        "accessment_item": "七彩石（Rainbow）2020.年早期及未知版本使用风险",
        "accessment_desc": "不推荐使用早期版本或者直接http接口接入，早期pullconfigreq接口，无配置监听，配置发布后需要业务主动拉取配置，主要风险项：1、无配置缓存功能   2、配置拉取动作和后台高度耦合； 3、七彩石服务器入口配死，无法做到灵活调整集群容灾切换；4、pullconfigreq高频调用挤占大量的拉取请求配额，容易触发限频，导致拉取配置失败",
        "risk_condition": "pullconfigreq协议的未知版本",
        "recommendation": {
            "template": "<div>建议升级到推荐的sdk版本；七彩石最佳实践请参考：<a href='{{.url_1}}' target='_blank'><span>最佳实践</span></a> </div>",
            "ph": {
                "url_1": "https://iwiki.woa.com/p/4008378767"
            }
        }
    }
}
`
	var td TemplateData
	err := json.Unmarshal([]byte(jsonStr), &td)
	if err != nil {
		log.Fatal(err)
	}

	result, err := renderTemplateData(&td)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}
