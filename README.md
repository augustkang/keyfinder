# keyfinder

keyfinder는 AWS 계정 내 IAM User들에 대해 주어진 시간보다 오래된 Access Key Pair를 찾아 Slack webhook으로 알림을 보내주는 애플리케이션입니다.

아래와 같이 동작합니다.

- AWS 계정 내 모든 IAM User들의 모든 Access Key Pair들 중에서
- 생성 후 N 시간 보다 오래된 Access Key Pair가 있는지 찾는다.
- 생성된 지 N시간 보다 오래된 Access Key Pair가 있을 경우, Slack webhook으로 해당 IAM User의 이름, 해당 Access Key ID, 해당 Access Key Pair가 생성된 날짜를 알림 메세지로 보낸다

keyfinder는 HTTP web server 기반으로 작동합니다. query string으로 hours와 url 이라는 매개변수를 받아 역할을 수행합니다. **로컬에서 구동하면 8080 포트를 이용합니다.**

아래 예시처럼 keyfinder에게 요청하면 생성된지 8 시간보다 오래된 Access Key Pair를 찾고 https://dummyslackwebhook.com/asdf 이라는 주소로 검출 결과를 http post 요청으로 보냅니다.

- 예시 : 브라우저에서 localhost:8080/?url=https://dummyslackwebhook.com/asdf&hours=8 입력

### Query string으로 제공할 수 있는 Parameter의 종류
- url(필수) : Slack webhook URL
- hours(필수) : 생성된 후 몇 시간 지난 Access Key Pair 들을 찾을건지
- channel(선택) : 해당 Slack workspace의 어느 Channel로 보낼 것인지. 기본값 : #example
- - !!주의 : channel=name 형태로 넘겨줘야 합니다. channel=#name 형식 안됨!!

### 메세지를 보낼 Slack 채널 지정하는 법

채널을 별도로 지정해주지 않을 경우 #example 채널에 요청을 보내게 코드에 넣어놓았습니다.(`main.go` 내 `func rootHandler`내에 정의)

다른 채널로 override 하고 싶을 경우가 있을 수 있으므로 그럴 경우 query string parameter로 받게 하였습니다.

!!주의 : channel 명이 예를 들어 '#examplechannel' 인 경우 'examplechannel' 이라고만 입력해야 해당 채널로 보냅니다.

## Prerequisites

### IAM Permissions

keyfinder가 작동하려면 IAM action들을 허용해줘야 합니다.

- iam:ListUsers
- iam:ListAccessKeys

## 실행하려면?

직접 실행하려면 AWS Access Key Pair를 `aws configure` 명령어로 준비해두어야 합니다.

또는 컨테이너로 실행한다면 Access Key Pair를 컨테이너에 환경변수로 넣어주어야 합니다.

Docker로 실행할 경우 `docker run` 명령어에 `-e` 옵션을 넣고 넣어줄 수 있습니다.

Kubernetes로 배포할 경우 manifest 내에 넣어주어야 합니다.

## Run locally

If run on local, keyfinder will use port number **8080**

```bash
[august@dummy-pc ~]$ git clone https://github.com/augustkang/keyfinder

[august@dummy-pc ~]$ cd keyfinder

[august@dummy-pc ~]$ go build -o main ./

[august@dummy-pc ~]$ ./main

# Or just run
[august@dummy-pc ~]$ go run main.go
```

**Then, open your browser and enter below address (localhost:8080) as belows**

'localhost:8080/?hours=24&url=https://dummyslackwebhook.com/asdf&channel=examplechannel'

## Run as Docker container

keyfinder needs IAM permission to run properly.

Please pass aws credential as environment variable to keyfinder container.

1. docker-compose

edit docker-compose.yaml to use your own AWS Credential.

```yaml
version: '3'

services:
  keyfinder:
    image: donghyunkang/keyfinder:latest
    environment:
    - AWS_ACCESS_KEY_ID= #PUT ACCESS KEY ID
    - AWS_SECRET_ACCESS_KEY= #PUT SECRET ACCESS KEY
    - AWS_REGION=ap-northeast-2
    ports:
    - "8080:8080"
```

Then run docker-compose

```bash
[august@dummy-pc ~]$ docker-compose up -d
```

2. docker run

```bash
# Docker run
[august@dummy-pc ~]$ docker run -e AWS_ACCESS_KEY_ID="YOUR-ACCESS-KEY" -e AWS_SECRET_ACCESS_KEY="YOUR-SECRET-KEY" -e AWS_REGION=ap-northeast-2 -d -p 8080:8080 --name keyfinder donghyunkang/keyfinder:latest
```

**Then, access 'localhost:8080/?hours=24&url=https://dummyslackwebhook.com/asdf&channel=examplechannel'**


### To build or pull image from Docker hub,

```bash
# Build image locally
[august@dummy-pc ~]$ docker build -t keyfinder:latest ./

# Pull from Docker hub
[august@dummy-pc ~]$ docker pull donghyunkang/keyfinder:latest
```

## Run on Kubernetes

### Prerequisites

Edit `keyfinder.yaml` to use AWS Access Key inside Container's `env` property.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: keyfinder
  ...
  ...
    spec:
      containers:
      - name: keyfinder
        image: donghyunkang/keyfinder:latest
        env:
        - name: AWS_ACCESS_KEY_ID
          value: "YOUR-ACCESS-KEY-ID"  # Put AWS Access Key ID Here
        - name: AWS_SECRET_ACCESS_KEY 
          value: "YOUR-SECRET-ACCESS-KEY" # Put AWS Secret Access Key Here
        - name: AWS_REGION
          value: "ap-northeast-2"
...
...
```

### Deploy keyfinder on Kubernetes

Please refer `keyfinder.yaml`. This manifest contains below Kubernetes objects

- Deployment : Deploys keyfinder pods.
- Service : Exposes keyfinder deployment NodePort type.

I used `nodePort` 32000. You can edit this port number as per your requirements.

```bash
[august@dummy-pc ~]$ kubectl apply -f keyfinder.yaml
deployment.apps/keyfinder created
service/keyfinder-service created

```

Then send reuqest to nodeIP:32000/?hours={hours}&url={WEBHOOK-URL}&channel={SLACK-CHANNEL-NAME} without '#' (For #mychannel channel? type mychannel only.)

Let's say hours (N) is 24, webhook url is http://abcdwebhook.com/asdf.
And slack channel name is #mychannel, for example.
Your Worker node's IP is 1.1.1.1, and nodePort set as 32000.

Then full URI (pass to keyfinder) would

- 1.1.1.1:32000/?hours=24&url=http://abcdwebhook.com/asdf&channel=mychannel

## Minikube

You may need to edit keyfinder-service.

After create deployment(or pod) and create service, Run below command to access service.
```bash
[august@dummy-pc ~]$ minikube service keyfinder
```

### Request to keyfinder

nodeIP:nodePort/?url=URL&hours=N

add channel parameter as query string to override Slack channel.
(Default : example)

## TODO
- Inject AWS Credential securely(Or retrieve from other source)
