import smtplib
# Here are the email package modules we'll need.
from email.message import EmailMessage

# Create the container email message.
msg = EmailMessage()
msg['Subject'] = 'Our family reunion'
# me == the sender's email address
# family = the list of all recipients' email addresses
msg['From'] = 'you@restsend.com'
msg['To'] = ', '.join(['oss@restsend.com'])
msg.preamble = 'You will not see this in a MIME-aware mail reader.\n'

img_data = open('./ui/src/assets/logo.png', 'rb').read()
msg.add_attachment(img_data, maintype='image',
                   subtype='png', filename='logo.png')

conn = smtplib.SMTP('localhost', port=9025)
# Mock login
conn.login('you@restsend.com', 'smtppreview')

r = conn.send_message(msg)

print('send ok, result:', r)
