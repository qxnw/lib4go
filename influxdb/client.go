package influxdb

import (
	"fmt"
	"log"
	uurl "net/url"
	"time"
	//"github.com/qxnw/lib4go/influxdb"
	//"github.com/qxnw/lib4go/influxdb"
)

type InfluxClient struct {
	interval time.Duration
	url      uurl.URL
	database string
	username string
	password string
	client   *Client
	closeCh  chan struct{}
	done     bool
}

// newInfluxClient starts a InfluxDB reporter which will post the metrics from the given registry at each d interval with the specified tags
func NewInfluxClient(url, database, username, password string) (*InfluxClient, error) {
	u, err := uurl.Parse(url)
	if err != nil {
		return nil, fmt.Errorf("unable to parse InfluxDB url %s. err=%v", url, err)
	}

	rep := &InfluxClient{
		url:      *u,
		database: database,
		username: username,
		password: password,
		closeCh:  make(chan struct{}),
	}

	if err := rep.makeClient(); err != nil {
		return nil, fmt.Errorf("unable to make InfluxDB client. err=%v", err)
	}
	go rep.run()
	return rep, nil
}

func (r *InfluxClient) makeClient() (err error) {
	r.client, err = NewClient(Config{
		URL:       r.url,
		Timeout:   time.Second * 3,
		UserAgent: "hydra",
		Username:  r.username,
		Password:  r.password,
	})
	return
}

func (r *InfluxClient) run() {
	pingTicker := time.Tick(time.Second * 5)
	for {
		select {
		case <-r.closeCh:
			return
		case <-pingTicker:
			_, _, err := r.client.Ping()
			if err != nil {
				log.Printf("got error while sending a ping to InfluxDB, trying to recreate client. err=%v", err)
				if err = r.makeClient(); err != nil {
					log.Printf("unable to make InfluxDB client. err=%v", err)
				}
			}
		}
	}
}
func (r *InfluxClient) Query(sql string) (result string, err error) {
	response, err := r.client.Query(Query{Command: sql, Database: r.database})
	if err != nil {
		err = fmt.Errorf("query.error:%v", err)
		return
	}
	if err = response.Error(); err != nil {
		return "", fmt.Errorf("response.error:%v", err)
	}
	buf, err := response.MarshalJSON()
	if err != nil {
		err = fmt.Errorf("query.result.marshal.error:%v", err)
		return
	}
	result = string(buf)
	return
}
func (r *InfluxClient) Send(measurement string, tags map[string]string, fileds map[string]interface{}) error {
	var pts []Point
	pts = append(pts, Point{
		Measurement: measurement,
		Tags:        tags,
		Fields:      fileds,
		Time:        time.Now(),
	})

	bps := BatchPoints{
		Points:   pts,
		Database: r.database,
	}
	_, err := r.client.Write(bps)
	return err
}
func (r *InfluxClient) Close() error {
	r.done = true
	close(r.closeCh)
	return nil
}