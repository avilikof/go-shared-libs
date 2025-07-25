package alerting

import (
	"errors"
	"math/rand"
	"time"

	"github.com/alex/go-shared-libs/alerts"
)

type AlertGenerator struct {
	stream Stream
}

var ErrAlertGenerationFailed = errors.New("alert generation failed")

func NewAlertGenerator(stream Stream) *AlertGenerator {
	return &AlertGenerator{
		stream: stream,
	}
}

func (g *AlertGenerator) Generate(version AlertVersion, frequency time.Duration) error {
	var alertBytes []byte
	var err error

	switch version {
	case VersionV1:
		alert := randomAlertV1()
		if alert == nil {
			return ErrAlertGenerationFailed
		}
		alertBytes = alert.Bytes()
	case VersionV2:
		alert := randomAlertV2()
		if alert == nil {
			return ErrAlertGenerationFailed
		}
		alertBytes, err = alert.MarshalJSON()
		if err != nil {
			return err
		}
	default:
		return ErrAlertGenerationFailed
	}

	err = g.stream.Publish("test.alert", alertBytes)
	if err != nil {
		return err
	}

	time.Sleep(frequency)
	return nil
}

type AlertVersion int

const (
	VersionV1 AlertVersion = iota
	VersionV2
)

func (v AlertVersion) String() string {
	switch v {
	case VersionV1:
		return "v1"
	case VersionV2:
		return "v2"
	default:
		return "unknown"
	}
}

func randomAlertV1() *alerts.Alert {
	return alerts.NewAlert(
		getRandomHash(),
		"legacy test alert",
		"hello world!!!",
		time.Now(),
		getRandomBool(),
	)
}

func randomAlertV2() *alerts.AlertV2 {
	alert := alerts.NewAlertV2(
		getRandomHash(),
		"test_source",
		getRandomSeverity(),
		getRandomType(),
		"Simulated test alert",
		getRandomHash(),
		time.Now(),
		alerts.AlertStateActive,
	)

	alert.AddLabel("env", getRandomEnv())
	alert.AddLabel("region", getRandomRegion())
	alert.AddAnnotation("dashboard_url", "https://grafana.example.com")
	alert.AddAnnotation("note", "This is a test alert.")
	alert.AddAction(alerts.Action{
		Type:       "ticket",
		Target:     "jira",
		AutoCreate: true,
	})
	alert.AddAction(alerts.Action{
		Type:             "callout",
		Target:           "pagerduty",
		EscalationPolicy: "dev_oncall",
	})

	if getRandomBool() {
		alert.SetCorrelationID("incident-" + getRandomHash())
	}

	return alert
}

func getRandomSeverity() string {
	severities := []string{"critical", "warning", "info"}
	return severities[rand.Intn(len(severities))]
}

func getRandomType() string {
	types := []string{"cpu_high", "memory_leak", "disk_full", "network_latency"}
	return types[rand.Intn(len(types))]
}

func getRandomEnv() string {
	envs := []string{"prod", "staging", "dev"}
	return envs[rand.Intn(len(envs))]
}

func getRandomRegion() string {
	regions := []string{"eu-central", "us-west", "ap-southeast"}
	return regions[rand.Intn(len(regions))]
}

func getRandomHash() string {
	hashes := []string{
		"a2b87da92118eb5c",
		"8ee30a4181511374",
		"09c2cc18b16a4532",
		"d0c767a2124af12f",
		"39c3d3d1fde8c534",
		"21cf066573d84d45",
		"f920ca6222dd6ff7",
		"8a5dda6736c28110",
		"9f9b3ffe9f0f99fe",
		"25bdcda6a5179602",
		"405c95346de4a2b8",
		"2f0b544ed92f19eb",
		"1c6549ba972da494",
		"5f356d5964d0ad6f",
		"def4834a003d4734",
		"2dd2f4c7c1840a42",
		"447ed719b512bf6e",
		"b3039cb20d84fd3c",
		"0f65af743397e745",
		"f690e743851632b9",
		"b42e2bab19f3b9b5",
		"8ddfe0d677de0a54",
		"649129afbe8bea86",
		"cb51c3e9ba0927db",
		"599ef51b58491aee",
		"cfdad65b91937d02",
		"dc3cda45485d787d",
		"4e3b87b0b2ecc59a",
		"29dd9eaebc79d9bd",
		"1a058f1c1e2c77d7",
		"a132178c48477146",
		"7eaf46c924e58720",
		"7ca1e991c78d043f",
		"ceda9204746746b0",
		"6b10c80de8be4ccc",
		"685b74b4a34ae3a0",
		"ab710b5f691d90ba",
		"f2624b41f0e6ad8b",
		"c464ba6a88ed0d5a",
		"3c4cce2ecf0d6296",
		"42d55405850b0647",
		"3112acd2368218a6",
		"bceecdcd78cc52f9",
		"3d34c6c37552ccec",
		"2b3730683f1b486c",
		"48d247db372272db",
		"a938a6a703056708",
		"13ccc45635bccdc0",
		"8aa6fb1c1b26a67a",
		"40deaa365ce8e13d",
		"8979bffbf13a604a",
		"74bec22553a9bda7",
		"1d10babd6a840be5",
		"fb247c969773b98b",
		"f4b85a151640bc20",
		"8edce7240623807c",
		"60c3f12f208dfe43",
		"c75f323f93f35f65",
		"96454e8ef41e043c",
		"0eeb3b1b4bb1487f",
		"0c57798bb66d409e",
		"a9ffa373578ceba0",
		"27f9a2826f116072",
		"6bbd8a7817277371",
		"900951881bfc7578",
		"b45deb51f8438ecb",
		"9193269b193f0ab9",
		"fde061e41932bcc9",
		"6e51e7e38a182d9f",
		"f03b9e764eeaecea",
		"5a6ef1b9d6cf04d5",
		"9a77e4790df05be4",
		"eec1e4480b67ab8f",
		"a777ec146de4c485",
		"7a7ebdc003cdfc58",
		"b8ba0700fad23ebb",
		"2cdab40f08123bb7",
		"acbe428fcafcecd6",
		"b93de064a21103b0",
		"4625bbdb911493a1",
		"fe69372ee04c9b88",
		"9d637ce4438ee9a8",
		"303ed7fe2d995b5b",
		"5f9ad8aa7511df49",
		"e0a762c5771b661c",
		"2635e0739d285937",
		"972a7319cb0771e7",
		"25f84dc266048aff",
		"7637bb1aaabf9f10",
		"115daf486cc9e293",
		"1774bc73a4111a08",
		"e3136f3410faac57",
		"8eabd4f3cf4c98a1",
		"9d51501fede0f641",
		"26a92a4818fc8c39",
		"dc328402766d5189",
		"5befe6f03d6928d0",
		"e997bf1b9c336774",
		"3d64dc6f9d3a9137",
		"d61bd39e222ab8e2",
		"7a3b0e21c8d06eec", "314c84010e7d920f", "9aafdfc884f54abd",
		"52b4c36b6e1d0e98", "1a2f3e3e4bdb09cb", "e56c98b760e3e5cd",
		"95d5ac5eae8c00f7", "4f5f35d9dc1ec474", "d028a4cb4d826317",
		"2a1c2f7ec0bd749e", "1a3b0e21c8d06eec", "114c84010e7d920f",
		"1aafdfc884f54abd", "12b4c36b6e1d0e98", "2a2f3e3e4bdb09cb",
		"156c98b760e3e5cd", "15d5ac5eae8c00f7", "1f5f35d9dc1ec474",
		"1028a4cb4d826317", "1a1c2f7ec0bd749e",
	}
	return hashes[rand.Intn(len(hashes))]
}

func getRandomBool() bool {
	return rand.Intn(2) == 1
}
