import smtplib
# Here are the email package modules we'll need.
from email.message import EmailMessage
from email.utils import make_msgid
# Create the container email message.
msg = EmailMessage()
msg['Subject'] = 'Hello SMTP Preview'
# me == the sender's email address
# family = the list of all recipients' email addresses
msg['From'] = 'you@restsend.com'
msg['To'] = ', '.join(['oss@restsend.com'])
msg.preamble = 'SMTP Preview is a tool for testing and inspecting email sending\n'

msg.set_content("""\
Salut!

Cela ressemble à un excellent recipie[1] déjeuner.

[1] http://www.yummly.com/recipe/Roasted-Asparagus-Epicurious-203718

--Pepé
""")

asparagus_cid = make_msgid()
msg.add_alternative("""\
<html>
  <head></head>
  <body>
    <p>Salut!</p>
    <p>Cela ressemble à un excellent
        <a href="http://www.yummly.com/recipe/Roasted-Asparagus-Epicurious-203718">
            recipie
        </a> déjeuner.
    </p>
    <img src="cid:{asparagus_cid}" />
  </body>
</html>
""".format(asparagus_cid=asparagus_cid[1:-1]), subtype='html')

with open("./ui/src/assets/logo.png", 'rb') as img:
    msg.get_payload()[1].add_related(img.read(), 'image', 'png',
                                     cid=asparagus_cid)

#img_data = open('./ui/src/assets/logo.png', 'rb').read()
# msg.add_attachment(img_data, maintype='image',
#                   subtype='png', filename='logo.png')

conn = smtplib.SMTP('localhost', port=9025)
# Mock login
conn.login('you@restsend.com', 'smtppreview')

r = conn.send_message(msg)

print('send ok, result:', r)
