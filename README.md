## Email Worker

A simple worker using GO and RabbitMQ. It consumes from a queue a JSON with the following structure and sends an e-mail to the recipient.

```json
{
	"recipient": "some.recipient@somemail.com",
	"subject": "Subject",
	"body": "<h1>Hello World</h1>"
}
```

It sends using TLS and the GO native SMTP lib. There's a `config.yaml.example` with the environment variables used by the application.

### TODO

- Retry 
- Dead letter exchange