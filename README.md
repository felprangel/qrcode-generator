# qrcode-generator

Ferramenta CLI que gera QR codes em ASCII no terminal a partir de números, opcionalmente prefixados por uma URL base configurável.

## Pré-requisitos

- Toolchain do [Go](https://go.dev/dl/) 1.26 ou superior

## Configuração

1. Clone o repositório:

```sh
git clone https://github.com/felprangel/qrcode-generator
cd qrcode-generator
```

2. (Opcional) Defina um prefixo global aplicado a cada número antes de gerar o QR code:

```sh
qrcode-generator config set PREFIX=https://exemplo.com/
```

O comando grava o valor em `~/.config/qrcode-generator/config` (ou em `$XDG_CONFIG_HOME/qrcode-generator/config`, se definido), preservando comentários e outras linhas do arquivo.

| Variável | Padrão | Descrição                                          |
| -------- | ------ | -------------------------------------------------- |
| `PREFIX` | _(vazio)_ | String concatenada antes do número ao gerar o QR |

## Instalação

```sh
go install ./...
```

Isso coloca dois binários em `$(go env GOBIN)` (ou `$(go env GOPATH)/bin`): `qrcode-generator` e o alias curto `qr`. Ambos têm o mesmo comportamento. Garanta que esse diretório esteja no seu `$PATH`.

## Uso

```sh
qrcode-generator [-c|--clear] <número> [número...]
qrcode-generator config set KEY=VALUE

# alias equivalente
qr [-c|--clear] <número> [número...]
qr config set KEY=VALUE
```

**Flags:**

| Flag             | Descrição                                       |
| ---------------- | ----------------------------------------------- |
| `-c`, `--clear`  | Limpa o terminal antes de gerar os QR codes     |

**Exemplos:**

```sh
qrcode-generator 123
qr 1 2 3
qr -c 42
qrcode-generator --clear 1 2 3
```

A saída é um QR code renderizado em ASCII (half-blocks) impresso no stdout, um por número informado.

## Desinstalação

```sh
rm "$(go env GOPATH)/bin/qrcode-generator" "$(go env GOPATH)/bin/qr"
```
