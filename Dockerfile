FROM golang:1.21-alpine AS build

WORKDIR /cmd

COPY go.mod ./

# RUN ls -la && go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o chess cmd/main.go


FROM alpine

WORKDIR /app

RUN mkdir -p configs bin

COPY --from=build ./cmd/chess /app/bin/chess
COPY --from=build ./cmd/configs/main.yml /app/configs/main.yml
# COPY --from=build . ./configs/main.yml

# RUN ls -la && mkdir -p configs && sleep 10s

# COPY --from=build configs/main.yml configs/main.yml

WORKDIR /app/bin

# ENV GIN_MODE=release
EXPOSE 8001
# CMD ["/bin/ls", "-la"]
CMD ["./chess"]