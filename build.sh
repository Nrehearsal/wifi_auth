#!/bin/bash

go build -o wifi_auth
mv wifi_auth resource
cd resource && ./wifi_auth
