FROM golang:1.19-bullseye AS build

WORKDIR /src

COPY ./go.sum ./go.mod ./

RUN go mod download

COPY . .

RUN mkdir ./output
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -a -installsuffix cgo -ldflags '-s -w -extldflags "-static"' -o ./output ./cmd/*

FROM gcr.io/distroless/static as final

COPY --from=build --chown=nonroot:nonroot /src/output /

USER nonroot:nonroot

ENTRYPOINT ["/roofwaterd"]
