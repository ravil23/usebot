FROM golang:1.13-alpine

WORKDIR /go/src/github.com/ravil23/usebot/telegrambot
COPY . .

RUN go get -v -d ./... \
    && go install -v ./...

ENTRYPOINT /go/bin/telegrambot \
    --russian /data/gia11/fipi/parsed/tasks_subject_russian.json \
    --math-advanced /data/gia11/fipi/parsed/tasks_subject_math_advanced.json \
    --math-basic /data/gia11/fipi/parsed/tasks_subject_math_basic.json \
    --physics /data/gia11/fipi/parsed/tasks_subject_physics.json \
    --chemistry /data/gia11/fipi/parsed/tasks_subject_chemistry.json \
    --it /data/gia11/fipi/parsed/tasks_subject_it.json \
    --biology /data/gia11/fipi/parsed/tasks_subject_biology.json \
    --history /data/gia11/fipi/parsed/tasks_subject_history.json \
    --geography /data/gia11/fipi/parsed/tasks_subject_geography.json \
    --english /data/gia11/fipi/parsed/tasks_subject_english.json \
    --german /data/gia11/fipi/parsed/tasks_subject_german.json \
    --french /data/gia11/fipi/parsed/tasks_subject_french.json \
    --social /data/gia11/fipi/parsed/tasks_subject_social.json \
    --spanish /data/gia11/fipi/parsed/tasks_subject_spanish.json \
    --literature /data/gia11/fipi/parsed/tasks_subject_literature.json
