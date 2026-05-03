# 🌟 Protos (GRPC Contracts)

Централизованное хранилище gRPC контрактов (`.proto` файлы).
Этот репозиторий/папка предназначен для использования в качестве зависимостей в микросервисах (Auth Service, Shop Service и др.).

---

## 📂 Структура

```text
protos/
├── Taskfile.yml          # Скрипт автоматической генерации
├── gen/
│   ├── go/               # Сгенерированный код для Go
│   │   └── sso/
│   │       ├── sso.pb.go
│   │       └── sso_grpc.pb.go
│   └── python/           # Сгенерированный код для Python
│       ├── __init__.py
│       └── sso/
│           ├── __init__.py
│           ├── sso_pb2.py
│           └── sso_pb2_grpc.py
└── proto/
    └── sso/
        └── sso.proto    # Контракты авторизации
```

---

## 🛠 Требования
| Язык | Инструменты | Команда установки |
|------|-------------|------------------|
| **Общие** | • **[protoc](https://github.com/protocolbuffers/protobuf/releases)** — компилятор протоколов.<br>• **[Task](https://taskfile.dev/installation/)** — менеджер задач (аналог Make). | Скачать по ссылкам <-- |
| **Go** | • `protoc-gen-go` — генератор Go-кода.<br>• `protoc-gen-go-grpc` — генератор gRPC-заглушек для Go. | `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`<br>`go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest` |
| **Python** | • `grpcio` — рантайм gRPC для Python.<br>• `grpcio-tools` — утилиты для генерации кода. | `pip install grpcio grpcio-tools` |

---

## ⚡Генерация кода
Для генерации Go-структур и gRPC-клиентов/серверов выполните в корневой папке protos:
```
task generate_go
# или
task gen_go
```
Результат появится в папке gen/go/proto/sso/.

Только Python
```
task generate_py
# или
task gen_py
```
Результат появится в папке gen/python/proto/sso/

---

## 🔌 Подключение в другие сервисы
▶️ Go-сервисы
- Добавьте репозиторий как подмодуль:
```
go get github.com/tambrama/protos
```
- Импортируйте сгенерированные пакеты в коде:

```
import pb "github.com/tambrama/protos/gen/go/proto/sso"
```
▶️ Python-сервисы (Flask / FastAPI)

1. Скопируйте или подключите сгенерированный код:
Вариант А (подмодуль):

```
git submodule add https://github.com/tambrama/protos.git
# Или скопируйте папку gen/python/ в ваш проект
```
Вариант Б (локальная генерация):
- Скопируйте gen/python/ в ваш Python-проект

2. Добавьте зависимости в requirements.txt:

```
grpcio>=1.60.0
grpcio-tools>=1.60.0
```

3. Импортируйте сгенерированные модули:

```
from gen.python.auth import auth_pb2, auth_pb2_grpc
```
Пример использования клиента:

```
# auth_client.py
import grpc
import os
from gen.python.auth import auth_pb2, auth_pb2_grpc

class AuthServiceClient:
    def __init__(self):
        host = os.getenv("AUTH_SERVICE_HOST", "localhost")
        port = os.getenv("AUTH_SERVICE_PORT", "50051")
        self.channel = grpc.insecure_channel(f"{host}:{port}")
        self.stub = auth_pb2_grpc.AuthStub(self.channel)
    
    def validate_token(self, token: str) -> dict:
        try:
            request = auth_pb2.ValidateRequest(token=token)
            metadata = [("authorization", f"Bearer {token}")]
            response = self.stub.Validate(request, metadata=metadata, timeout=5.0)
            return {"valid": True, "user_id": response.user_id, "app_id": response.app_id}
        except grpc.RpcError as e:
            return {"valid": False, "error": str(e.details())}
        finally:
            self.channel.close()
```
---

## 🔄 Обновление контрактов

1. Внесите изменения в proto/sso/auth.proto
2. Запустите генерацию: task generate
3. Закоммитьте изменения в proto/ и gen/
4. Обновите зависимость в микросервисах

```
# Go
go get -u github.com/tambrama/protos@latest

# Python
# Скопируйте обновленную папку gen/python/ в ваши сервисы
```

---