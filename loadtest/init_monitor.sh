#!/bin/zsh

#kubectl -n devspace-pro set env deployment/loft LOFTDEBUG=true

kubectl -n devspace-pro port-forward $(kubectl -n devspace-pro get pods -l app=loft -o jsonpath="{.items[0].metadata.name}") 8080:8080
