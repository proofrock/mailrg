# MailRG

This app implements a server that opens a POST endpoint that, when invoked with a JSON that describes an email message, sends the message.

The endpoint is `serve/`, and the JSON is:

```bash
{
    "token": "CorrectHorseBatteryStaple",
    "from": "john.doe@gmail.com",
    "to": ["me@johndoe.it"],
    "cc": ["you@johndoe.it"],
    "bcc": ["secret@johndoe.it"],
    "subject": "Hey!",
    "body": "Fun! <b>Cool!</b>",
    "html": true,
    "attachments": ["file1.txt", "file2.pdf"]
}
```

Note that the attachments are relative to the `DATA_DIR` specified, and `cc`, `bcc`, `to` and `attachments` nodes can be omitted.

It can be installed via Docker:

```bash
docker run -p 2163:2163 -v $(pwd):/data \
  -e SMTP_SERVER=smtp.gmail.com \
  -e SMTP_USER="john.doe@gmail.com" \
  -e SMTP_PASS="xyz" \
  -e MAILRG_TOKEN="CorrectHorseBatteryStaple" \
  mailrg
```