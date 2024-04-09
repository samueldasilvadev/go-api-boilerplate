FROM golang:1.22-alpine as build
WORKDIR /app
COPY . .

RUN go mod download
RUN go build -tags lambda.norpc cmd/aws_lambda/main.go

FROM public.ecr.aws/lambda/provided:al2023

COPY --from=build /app/main ./main

ENV ENVIRONMENT=development
ENV VERSION=0.0.1
ENV APP=go-skeleton
ENV DB_USER=root
ENV DB_PASS=root
ENV DB_URL=192.168.101.9
ENV DB_PORT=3306
ENV DB_DATABASE=skeleton
ENV DB_DRIVER=mysql

ENV DB_TEST_DATABASE=test

ENV JWT_SECRET=secret
ENV JWT_EXPIRATION=10

ENV ACCESS_CONTEXT=teste
ENV ACCESS_SECRET=yweurigyquwruieygqw
ENV ACCESS_TOKEN=yweurigyquwruieygqwasda

ENV HTTP_PORT=9000


RUN touch .env

ENTRYPOINT [ "./main" ]
