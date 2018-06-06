## Golang Brasil Translation Task Force
[![Build Status](https://travis-ci.org/golangbr/golangbr.org.svg?branch=master)](https://travis-ci.org/golangbr/golangbr.org)

Estrutura do Repositório:
```
|+godoc/    godoc q roda no GAE
|+newgo/    deps p/ godoc rodar lá
|~pt_BR/    dir dos arquivos pro godoc subir
| |-doc/    dir com os arquivos traduzidos <==
| |-lib/    dir de templates e assets p/ godoc
|~app.yaml  config do GAE
```

### Subindo localmente

```
$ wget https://storage.googleapis.com/appengine-sdks/featured/go_appengine_sdk_linux_amd64-1.9.40.zip; unzip go_appengine_sdk_linux_amd64-1.9.40.zip
$ git clone https://github.com/golangbr/golangbr.org.git
$ go_appengine/dev_appserver.py golangbr.org
$ chromium-browser localhost:8080
```
Para outras versões etc..
+ [https://developers.google.com/appengine/downloads?hl=pt-br#Google_App_Engine_SDK_for_Go](https://developers.google.com/appengine/downloads?hl=pt-br#Google_App_Engine_SDK_for_Go)


### Contribuindo

+ Abra uma issue informando no que deseja colaborar, ou está trabalhando.
+ Inclua seu nome/email no arquivo CONTRIBUIDORES.
+ Dê o pull/request ;-)
