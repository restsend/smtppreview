SMTP Preview - SMTP Test tools
======

SMTP Preview is a tool for testing and inspecting email sending. It allows you to preview the layout and notification appearance of emails on different clients, as well as detecting potential SPAM warnings. It can be integrated into the development process for SMTP sending tests.

SMTP Preview  is maintained by [restsend.com](https://restsend.com/?from=smtppreview_github).

![Screenshot](ui/src/assets/screenshot1.jpeg)

## Quick start
>
> SMTP Preview is written in Golang, and running smtppreview requires the Golang development environment. If you do not have the Golang development environment, you can refer to the official documentation of [golang](https://golang.org/?from=smtppreview).

- Build `smtppreivew`

```shell
git clone https://github.com/restsend/smtppreview.git
cd smtppreview
go build 
```

- Run `smtppreview`, the default serve port is `9025`

```shell
./smtppreview
```

- Visit web console, `http://localhost:8000`

- Send test mail via python code:

```python
import smtplib
from email.mime.text import MIMEText
from email.header import Header

message = MIMEText('Hello SMTP Preview', 'plain', 'utf-8')
message['From'] = Header("unittest@restsend.com", 'utf-8')
message['To'] = Header("oss@restsend.com", 'utf-8')

subject = 'Hello SMTP Preview'
message['Subject'] = Header(subject, 'utf-8')

conn = smtplib.SMTP('localhost', port=9025)
# Mock login
conn.login('you@restsend.com', 'smtppreview')

r = conn.sendmail('you@restsend.com',
                  ['oss@restsend.com', 'hello@restsend.com'],
                  message.as_string())

print('send ok, result:', r)
```

## How SMTP Preview work?
