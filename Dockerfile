FROM golang:1.13-alpine

WORKDIR /go/src/github.com/ravil23/usebot
COPY ./data/gia11/fipi/parsed /data
COPY ./telegrambot ./telegrambot

RUN cd telegrambot \
    && go get -v -d ./... \
    && go install -v ./...

ENTRYPOINT /go/bin/telegrambot \
    --russian /data/tasks_subject_russian.json \
    --math-advanced /data/tasks_subject_math_advanced.json \
    --math-basic /data/tasks_subject_math_basic.json \
    --physics /data/tasks_subject_physics.json \
    --chemistry /data/tasks_subject_chemistry.json \
    --it /data/tasks_subject_it.json \
    --biology /data/tasks_subject_biology.json \
    --history /data/tasks_subject_history.json \
    --geography /data/tasks_subject_geography.json \
    --english /data/tasks_subject_english.json \
    --german /data/tasks_subject_german.json \
    --french /data/tasks_subject_french.json \
    --social /data/tasks_subject_social.json \
    --spanish /data/tasks_subject_spanish.json \
    --literature /data/tasks_subject_literature.json
