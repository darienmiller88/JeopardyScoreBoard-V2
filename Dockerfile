FROM golang:1.25 AS build-stage

WORKDIR /app

# First, copy mod and sum files to directory,
COPY go.mod go.sum ./

# Afterwards, download the dependencies from the go.mod file
RUN go mod download

# Copy the entire directory into the build folder
COPY . .

RUN go build -o jeopardy-scoreboard .


# This is the run stage now, pulling from gcr
FROM gcr.io/distroless/base-debian12

# Create a new directory for the run time image
WORKDIR /app

# Copy the binary create during the build stage into the run time image
COPY --from=build-stage /app/jeopardy-scoreboard .

# Also copy the templates and static folders into the run time image
COPY --from=build-stage /app/templates ./templates
COPY --from=build-stage /app/static ./static

# expose port 8080 for the local machine
EXPOSE 8080

# Finally, run the go binary!
CMD [ "./jeopardy-scoreboard" ]