FROM golang:1.13-alpine

WORKDIR /go/src/github.com/ravil23/usebot
COPY ./data/gia11/fipi/parsed /data
COPY ./telegrambot ./telegrambot

RUN cd telegrambot \
    && go get -v -d ./... \
    && go install -v ./...

ENV SUBJECT_RUSSIAN="/data/tasks_subject_russian.json"
ENV SUBJECT_MATH_ADVANCED="/data/tasks_subject_math_advanced.json"
ENV SUBJECT_MATH_BASIC="/data/tasks_subject_math_basic.json"
ENV SUBJECT_PHYSICS="/data/tasks_subject_physics.json"
ENV SUBJECT_CHEMISTRY="/data/tasks_subject_chemistry.json"
ENV SUBJECT_IT="/data/tasks_subject_it.json"
ENV SUBJECT_BIOLOGY="/data/tasks_subject_biology.json"
ENV SUBJECT_HISTORY="/data/tasks_subject_history.json"
ENV SUBJECT_GEOGRAPHY="/data/tasks_subject_geography.json"
ENV SUBJECT_ENGLISH="/data/tasks_subject_english.json"
ENV SUBJECT_GERMAN="/data/tasks_subject_german.json"
ENV SUBJECT_FRENCH="/data/tasks_subject_french.json"
ENV SUBJECT_SOCIAL="/data/tasks_subject_social.json"
ENV SUBJECT_SPANISH="/data/tasks_subject_spanish.json"
ENV SUBJECT_LITERATURE="/data/tasks_subject_literature.json"

ENTRYPOINT /go/bin/telegrambot