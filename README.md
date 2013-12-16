##Golang-BR Translation Task Force

Estrutura do Repositório:
```
|+godoc/    godoc q roda no GAE
|+newgo/    deps p/ godoc rodar lá
|~pt_BR/    dir dos arquivos pro godoc subir
| |-doc/    **dir com os arquivos traduzidos**
| |-lib/    dir de templates e assets p/ godoc
|~app.yaml  config do GAE
```

### Subindo localmente

```
$ wget http://googleappengine.googlecode.com/files/go_appengine_sdk_linux_amd64-1.8.8.zip ; unzip go_appengine_sdk_linux_amd64-1.8.8.zip
$ git clone https://github.com/golang-br/golang-br.appspot.com.git
$ go_appengine/dev_appserver.py golang-br.appspot.com.git/
$ chromium-browser localhost:8080 
```
Para outras versões etc.. 
[https://developers.google.com/appengine/downloads?hl=pt-br#Google_App_Engine_SDK_for_Go](https://developers.google.com/appengine/downloads?hl=pt-br#Google_App_Engine_SDK_for_Go)
