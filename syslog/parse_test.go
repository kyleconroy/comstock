package syslog

import (
	"testing"
)

var logs = []byte(`64 <190>1 2015-11-05T08:37:14.373088+00:00 host app web.1 - 世界
74 <190>1 2015-11-05T08:37:14.373092+00:00 host app web.1 - Many new lines: 
58 <190>1 2015-11-05T08:37:14.373093+00:00 host app web.1 - 
58 <190>1 2015-11-05T08:37:14.373094+00:00 host app web.1 - 
58 <190>1 2015-11-05T08:37:14.373095+00:00 host app web.1 - 
58 <190>1 2015-11-05T08:37:14.373095+00:00 host app web.1 - 
458 <190>1 2015-11-05T08:37:14.373101+00:00 host app web.1 - foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo 
`)

var demo = []byte(`83 <40>1 2012-11-30T06:45:29+00:00 host app web.3 - State changed from starting to up
119 <40>1 2012-11-30T06:45:26+00:00 host app web.3 - Starting process with command 'bundle exec rackup config.ru -p 24405'
`)

func TestParseMessage(t *testing.T) {
	messages, err := ParseFrame(logs, 7)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(messages)
}
