# Exchange rate notifier API
App for getting every day notification about NBU currency rate USD to UAH

## Stack 💻
[![My Skills](https://skillicons.dev/icons?i=golang,postgresql)](https://skillicons.dev)

## Description :cyclone:
App stores in PostgreSQL all emails. Then at 8:00 (crone job, can be set up in /configs/config.yml) service requests from [NBU API](https://bank.gov.ua/ua/open-data/api-dev) current currency rate and sends it to all emails using Gmail SMTP.

Used style guidelines: https://github.com/golang-standards/project-layout

## Example of notification :milky_way:
![image](https://github.com/DimaL-cloud/exchange-rate-notifier-api/assets/78265212/b5acec7a-cb79-4416-985e-ebeb0ed74523)


## Setting up :rocket:
1. Clone the repository:
```
git clone https://github.com/DimaL-cloud/exchange-rate-notifier-api.git
```
2. Configure .env file
3. [Optional] Set up crone job
4. Up services:
```
docker-compose up -d
```

## Alerts :warning:
1. Total failed rate notification jobs > 0
- if rate notification jobs was failed, users won't get notifications
2. Microservice is not alive
- if one of the services is not alive, users can face with difficulties
3. Total email send errors > 0
- if email was not sent, users won't get notifications