package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

func handler(request events.CloudWatchEvent) (s string, err error) {

	now := time.Now()
	from := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)

	to := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	layout := "2006-01-02"

	granularity := "MONTHLY"
	metrics := []string{
		"AmortizedCost",
	}
	// Initialize a session in us-east-1 that the SDK will use to load credentials
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	svc := costexplorer.New(sess)

	result, err := svc.GetCostAndUsage(&costexplorer.GetCostAndUsageInput{
		TimePeriod: &costexplorer.DateInterval{
			Start: aws.String(from.Format(layout)),
			End:   aws.String(to.Format(layout)),
		},
		Granularity: aws.String(granularity),
		GroupBy: []*costexplorer.GroupDefinition{
			{
				Type: aws.String("DIMENSION"),
				Key:  aws.String("SERVICE"),
			},
		},
		Metrics: aws.StringSlice(metrics),
	})
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	cost := cost(result)
	fmt.Println("cost:" + cost)
	params := Slack{
		Text:     "Total cost of this month: " + cost + "USD.",
		Username: "AWS Service Billing Monthly",
		Channel:  os.Getenv("Channel"),
		Emoji:    ":money_with_wings:",
	}
	jsonparams, _ := json.Marshal(params)
	resp, err := http.PostForm(
		os.Getenv("SlackWebhookUrl"),
		url.Values{"payload": {string(jsonparams)}},
	)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return "", nil
}

func cost(out *costexplorer.GetCostAndUsageOutput) string {
	total := out.ResultsByTime[0].Total["AmortizedCost"]
	if total == nil {
		return "0"
	}
	return *total.Amount

}

// Slack slack.
type Slack struct {
	Text     string `json:"text"`
	Username string `json:"username"`
	Channel  string `json:"channel"`
	Emoji    string `json:"icon_emoji"`
}

func main() {
	lambda.Start(handler)
}
