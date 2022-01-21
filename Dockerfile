# This Dockerfile builds a lightweight distribution image.
# It can be used to run atwhy without installing it in your system.
FROM golang:1.17-bullseye AS builder

RUN mkdir /app
COPY . /app
WORKDIR /app

RUN make dep
RUN make

# Final stage only containing the binary
FROM gcr.io/distroless/base-debian11 AS final

LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.name="atwhy"
LABEL org.label-schema.description="A simple CLI for tracking your working time."
LABEL org.label-schema.url="https://github.com/Tiffinger-Thiel-GmbH/atwhy"
LABEL org.label-schema.vcs-url="https://github.com/Tiffinger-Thiel-GmbH/atwhy"
LABEL org.opencontainers.image.source = "https://github.com/Tiffinger-Thiel-GmbH/atwhy"

COPY --from=builder /app/atwhy /bin

WORKDIR /project

CMD ["atwhy"]