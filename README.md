# aws-slack-billing
Notify your AWS account monthly cost to your Slack channel everyday

## Installation
https://ap-northeast-1.console.aws.amazon.com/lambda/home?region=ap-northeast-1#/create/app?applicationId=arn:aws:serverlessrepo:ap-northeast-1:495476032358:applications/aws-notify-billing
Deploy to your AWS environment

- `Channel` is the slack channel to be notified at
- `SlackWebhookUrl` is incoming-webhook that associated with the slack thannel to be notified at. You must generate this url before deploying this application
- `TZ` is TimeZone this application works at
