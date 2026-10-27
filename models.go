// Package contracts 定义跨端使用的 v1 探测协议，不依赖具体服务或界面。
package contracts

const Protocol = 1
const Version = "0.1.1"

type Check struct {
	Type               string            `json:"type"`
	Target             string            `json:"target"`
	Port               int               `json:"port"`
	Timeout            int               `json:"timeout"`
	Method             string            `json:"method"`
	Headers            map[string]string `json:"headers,omitempty"`
	ExpectedStatus     int               `json:"expected_status"`
	Contains           string            `json:"contains,omitempty"`
	FollowRedirects    bool              `json:"follow_redirects"`
	InsecureTLS        bool              `json:"insecure_tls"`
	CertificateMinDays int               `json:"certificate_min_days"`
	DNSRecord          string            `json:"dns_record"`
	Resolver           string            `json:"resolver,omitempty"`
	ExpectedDNS        string            `json:"expected_dns,omitempty"`
	Count              int               `json:"count"`
	Size               int               `json:"size"`
	MaxLoss            float64           `json:"max_loss"`
	Family             string            `json:"family"`
	Retries            int               `json:"retries"`
}

type Monitor struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Group            string   `json:"group"`
	Tags             []string `json:"tags"`
	Check            Check    `json:"check"`
	NodeIDs          []string `json:"node_ids"`
	ChannelIDs       []string `json:"channel_ids"`
	Interval         int      `json:"interval"`
	FailThreshold    int      `json:"fail_threshold"`
	RecoverThreshold int      `json:"recover_threshold"`
	MinFailedNodes   int      `json:"min_failed_nodes"`
	Enabled          bool     `json:"enabled"`
	Archived         bool     `json:"archived"`
	Version          int      `json:"version"`
	UpdatedAt        int64    `json:"updated_at"`
	Public           bool     `json:"public"`
}

type Policy struct {
	AllowCIDRs []string `json:"allow_cidrs"`
}

type Assignment struct {
	ID      string  `json:"id"`
	Monitor Monitor `json:"monitor"`
	NodeID  string  `json:"node_id"`
	Start   int64   `json:"start"`
	Until   int64   `json:"until"`
	Policy  Policy  `json:"policy"`
}

type Job struct {
	ID        string `json:"id"`
	NodeID    string `json:"node_id"`
	Check     Check  `json:"check"`
	Policy    Policy `json:"policy"`
	Status    string `json:"status"`
	CreatedAt int64  `json:"created_at"`
	Deadline  int64  `json:"deadline"`
}

type Result struct {
	ID           string         `json:"id"`
	AssignmentID string         `json:"assignment_id,omitempty"`
	MonitorID    string         `json:"monitor_id,omitempty"`
	JobID        string         `json:"job_id,omitempty"`
	NodeID       string         `json:"node_id"`
	Version      int            `json:"version"`
	ScheduledAt  int64          `json:"scheduled_at"`
	StartedAt    int64          `json:"started_at"`
	ReceivedAt   int64          `json:"received_at"`
	Status       string         `json:"status"`
	Stage        string         `json:"stage"`
	Error        string         `json:"error,omitempty"`
	Latency      *float64       `json:"latency_ms"`
	Data         map[string]any `json:"data"`
	Late         bool           `json:"late"`
	Attempts     int            `json:"attempts"`
	FirstError   string         `json:"first_error,omitempty"`
}

type Poll struct {
	Protocol     int      `json:"protocol"`
	Version      string   `json:"version"`
	Platform     string   `json:"platform"`
	Capabilities []string `json:"capabilities"`
	Queue        int      `json:"queue"`
	Dropped      int      `json:"dropped"`
	Time         int64    `json:"time"`
}

type Work struct {
	Assignments []Assignment `json:"assignments"`
	Jobs        []Job        `json:"jobs"`
	Cancelled   []string     `json:"cancelled"`
	ServerTime  int64        `json:"server_time"`
}
