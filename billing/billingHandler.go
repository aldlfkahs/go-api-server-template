package billing

import (
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"k8s.io/klog"
)

func Get(res http.ResponseWriter, req *http.Request) {

	fileName := "/root/.aws/credentials"
	dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		klog.Errorln("Cannot open credential file")
		klog.Errorln(err)
	}

	lines := strings.Split(string(dat), "\n")
	reg := regexp.MustCompile("\\[.*\\]")

	var output *costexplorer.GetCostAndUsageOutput
	for _, account := range lines {
		if reg.MatchString(account) {
			account = strings.TrimLeft(account, "[")
			account = strings.TrimRight(account, "]")
			klog.Infoln(account)

			output, err = makeCost(req, account)

			if err != nil {
				klog.Errorln(err)
			} else {
				klog.Infoln(output)
			}
		}
	}

	//klog.Infoln(js)
	// res.Header().Set("Content-Type", "application/json")
	// js, err := json.Marshal(result)
	// if err != nil {
	// 	klog.Errorln(err)
	// }
	// res.Write(js)

	// klog.Infoln("Cost Report:", result.ResultsByTime)

}

func makeCost(req *http.Request, account string) (*costexplorer.GetCostAndUsageOutput, error) {

	queryParams := req.URL.Query()

	/*** GET QUERY PARAMS ***/
	//Must be in YYYY-MM-DD Format
	var startTime int64
	var endTime int64
	startUnix := queryParams.Get("startTime")
	endUnix := queryParams.Get("endTime")
	if startUnix == "" || endUnix == "" {
		klog.Errorln("Must pass both of startTime and endTime")
		return nil, errors.New("Time parameter error")
	}
	startTime, _ = strconv.ParseInt(startUnix, 10, 64)
	endTime, _ = strconv.ParseInt(endUnix, 10, 64)

	granularity := queryParams.Get("granularity") // "MONTHLY"
	if granularity == "" {
		granularity = "MONTHLY"
	}

	// "AmortizedCost", "NetAmortizedCost", "BlendedCost", "UnblendedCost", "NetUnblendedCost", "UsageQuantity", "NormalizedUsageAmount",
	metrics := queryParams["metrics"]
	if len(metrics) == 0 {
		metrics = []string{
			"BlendedCost",
		}
	}

	dimension := queryParams.Get("dimension")
	if dimension == "" {
		dimension = "INSTANCE_TYPE"
	}

	/*** GET CREDENTIALS BY READING /root/.aws/credentials ***/
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: account,
		//SharedConfigState: session.SharedConfigEnable,
	})
	svc := costexplorer.New(sess)

	/*** GET COST FROM AWS ***/
	result, err := svc.GetCostAndUsage(&costexplorer.GetCostAndUsageInput{
		TimePeriod: &costexplorer.DateInterval{
			Start: aws.String(time.Unix(startTime, 0).Format("2006-01-02")),
			End:   aws.String(time.Unix(endTime, 0).Format("2006-01-02")),
		},
		Granularity: aws.String(granularity),
		GroupBy: []*costexplorer.GroupDefinition{
			&costexplorer.GroupDefinition{
				Type: aws.String("DIMENSION"),
				Key:  aws.String(dimension),
			},
		},
		Metrics: aws.StringSlice(metrics),
	})
	if err != nil {
		klog.Errorln(err)
	}

	return result, err
}
