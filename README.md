# Old Access Key Pair Finder

keyfinder는 AWS 계정 내 모든 IAM User들의 모든 Access Key Pair들에 대해 생성된지 N 시간 보다 오래된 Access Key Pair가 있는지 찾아보고,

N시간 보다 오래된 Access Key Pair가 있을 경우 Slack webhook으로 해당 IAM User의 이름, 해당 Access Key ID, 해당 Access Key Pair가 생성된 날짜를

알림 메세지로 보내주는 애플리케이션입니다. 

keyfinder는 HTTP web server 기반으로 작동합니다. query string으로 hours와 url 이라는 매개변수를 받아 역할을 수행합니다.

아래 예시처럼 keyfinder에게 요청하면 생성된지 8 시간보다 오래된 Access Key Pair를 찾고 https://dummyslackwebhook.com/asdf 이라는 주소로 검출 결과를 http post 요청으로 보냅니다.

- 예시 : 브라우저에서 localhost:8080/?url=https://dummyslackwebhook.com/asdf&hours=8 입력

## Limitations

### Access Key Pair는 컨테이너에 환경변수로 직접 넣어주어야 합니다.

Docker로 실행할 경우 `docker run` 명령어에 `-e` 옵션을 넣고 넣어줄 수 있습니다.

Kubernetes로 배포할 경우 manifest 내에 넣어주어야 합니다.

### 메세지를 보낼 Slack 채널 지정하는 법

별도로 지정해주지 않을 경우 #example 채널에 요청을 보내게 넣어놓았습니다.(`main.go` 내 `func rootHandler`내에 정의)

다른 채널로 override 하고 싶을 경우가 있을 수 있으므로 그럴 경우 query string parameter로 받게 하였습니다.

!!주의 : channel 명이 예를 들어 #example 인 경우 example 이라고만 입력해야 해당 채널로 보냅니다.

- #examplechannel로 보내는 예시 : localhost:8080/?channel=examplechannel&url=https://dummyslackwebhook.com/asdf&hours=8

## Run locally

### Prerequisites

You may need AWS credential with proper IAM Permissions to run keyfinder locally.

```bash
git clone https://github.com/augustkang/keyfinder

cd keyfinder

go build -o main ./

./main

# Or just run
go run main.go
```

## Run as Docker container
```

# Pull from Docker hub
docker pull donghyunkang/keyfinder:latest

# Build image locally
docker build -t keyfinder:latest ./

# Docker run (with pulled image)
docker run -e AWS_ACCESS_KEY_ID="YOUR-ACCESS-KEY" -e AWS_SECRET_ACCESS_KEY="YOUR-SECRET-KEY" -e AWS_REGION=ap-northeast-2 -d -p 8080:8080 --name keyfinder donghyunkang/keyfinder:latest

```

## Run on Kubernetes

### Prerequisites

Edit `keyfinder.yaml` to use AWS Access Key inside Container's `env` property.

`vim keyfinder.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: keyfinder
spec:
  replicas: 1
  selector:
    matchLabels:
      app: keyfinder
  template:
    metadata:
      labels:
        app: keyfinder
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
```

### How to deploy and expose keyfinder application as Kubernetes Service
```
# create deployment
kubectl apply -f keyfinder.yaml

# expose deployment as service
kubectl expose deployment keyfinder --type=LoadBalancer --port=8080

# If using Minikube, Run below command to access service
minikube service keyfinder
```
## Usage

### Browser

Enter localhost:portnumber/?url=https://SLACK-WEBHOOK-URL&hours=N from Browser
