package main

import (
	"bytes"
	"io"
	"testing"

	"github.com/DusanKasan/parsemail"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	data := []byte(`Subject: Our family reunion
From: you@restsend.com
To: oss@restsend.com
Content-Type: multipart/mixed; boundary="===============8642179929463994206=="

You will not see this in a MIME-aware mail reader.

--===============8642179929463994206==
Content-Type: image/png
Content-Transfer-Encoding: base64
Content-Disposition: attachment; filename="logo.png"
MIME-Version: 1.0

aW1nX2RhdGE=

--===============8642179929463994206==--`)
	r := bytes.NewReader(data)
	msg, err := parsemail.Parse(r)
	assert.Nil(t, err)
	assert.Equal(t, "Our family reunion", msg.Subject)
	assert.Equal(t, 1, len(msg.Attachments))
	assert.Equal(t, "logo.png", msg.Attachments[0].Filename)
	data, err = io.ReadAll(msg.Attachments[0].Data)
	assert.Nil(t, err)
	assert.Equal(t, []byte("img_data"), data)
}
