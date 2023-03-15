# go-musthave-devops-tpl

Шаблон репозитория для практического трека «Go в DevOps».

# Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` - адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

# Обновление шаблона

Чтобы получать обновления автотестов и других частей шаблона, выполните следующую команду:

```
git remote add -m main template https://github.com/yandex-praktikum/go-musthave-devops-tpl.git
```

Для обновления кода автотестов выполните команду:

(Для Unix систем)

```
git fetch template && git checkout template/main .github
```

(Для Windows PowerShell)

```
(git fetch template) -and (git checkout template/main .github)
```

Затем добавьте полученные изменения в свой репозиторий.



// https://golangci-lint.run/usage/install/
//
//C:/Golang/devopstest.exe
//-test.v -test.run=^TestIteration2[b]*$
//-binary-path=C:\Golang\go-advanced-devops-tpl\cmd\server\server.exe
//-source-path=C:\Golang\go-advanced-devops-tpl\cmd\server
//
//go build -race ./
//
//c:\apps\devopstest-windows-amd64.exe
//--test.v
//--test.run=^TestIteration2$
//--binary-path=./cmd/server/server.exe
//--agent-binary-path=./cmd/agent/agent.exe
