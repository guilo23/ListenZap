# ListenZap

Pipeline de monitoramento de mensagens do WhatsApp em tempo real, com inserção de contatos na agenda do Android via API externa.

## Arquitetura
**Componentes:**
- **AndroidSide** — binário nativo em Go que roda dentro do emulador Android com root
- **Redis** — intermediário entre o agente e a API
- **Api** — aplicação .NET que consome mensagens e expõe endpoint de contatos

## Pré-requisitos

- Emulador Android com root (GenyMotion + Android 10 + OpenGApps)
- WhatsApp instalado e ativado no emulador
- Redis 5.0+
- Go 1.21+
- .NET 8+
- ADB instalado

## Como rodar

### 1. Redis

```bash
docker run -d -p 6379:6379 redis
```

### 2. Agente Go (AndroidSide)

Compilar:
```bash
cd AndroidSide
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o agent .
```

No Windows (PowerShell):
```powershell
$env:GOOS="linux"; $env:GOARCH="amd64"; $env:CGO_ENABLED="0"; go build -o agent .
```

Enviar para o emulador e executar:
```bash
adb push agent /data/local/tmp/agent
adb shell chmod +x /data/local/tmp/agent
adb shell /data/local/tmp/agent
```

### 3. API .NET

```bash
cd Api
dotnet run
```

A API sobe em `http://localhost:5205`

## Endpoints

### POST /contacts

Insere um contato na agenda do Android.

**PowerShell:**
```powershell
Invoke-WebRequest -Uri "http://localhost:5205/api/Contacts"
-Method POST -ContentType "application/json"
-Body '{"name":"Joao Silva","number":"5511999999999"}'
```

**Body:**
```json
{
  "name": "Nome do Contato",
  "number": "5511999999999"
}
```

**Response:** `202 Accepted`
