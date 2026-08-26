# qrcode-generator

Ferramenta CLI que gera QR codes em ASCII no terminal a partir de qualquer texto, opcionalmente prefixado por um _preset_ configurável.

## Pré-requisitos

- Toolchain do [Go](https://go.dev/dl/) 1.26 ou superior

## Configuração

1. Clone o repositório:

```sh
git clone https://github.com/felprangel/qrcode-generator
cd qrcode-generator
```

2. (Opcional) Defina _presets_ — prefixos nomeados concatenados antes do texto antes de gerar o QR code:

```sh
qrcode-generator config set zap=https://wa.me/
qrcode-generator config set site=https://exemplo.com/
qrcode-generator config default zap   # usa "zap" quando -p é omitido
```

O comando grava os valores em `~/.config/qrcode-generator/config` (ou em `$XDG_CONFIG_HOME/qrcode-generator/config`, se definido), preservando comentários e outras linhas do arquivo. Nomes de preset são case-insensitive.

Selecione um preset na hora de gerar com `-p NOME`. Sem `-p`, é usado o preset default (ou o preset `prefix`, mantido por compatibilidade com versões antigas que usavam `PREFIX=`).

| Comando                        | Descrição                                        |
| ------------------------------ | ------------------------------------------------ |
| `config set NOME=PREFIXO`      | Cria ou atualiza um preset                       |
| `config default NOME`          | Define qual preset é usado quando `-p` é omitido |

## Instalação

```sh
go install ./...
```

Isso coloca dois binários em `$(go env GOBIN)` (ou `$(go env GOPATH)/bin`): `qrcode-generator` e o alias curto `qr`. Ambos têm o mesmo comportamento. Garanta que esse diretório esteja no seu `$PATH`.

## Uso

```sh
qrcode-generator [-c|--clear] [-p|--preset NOME] <texto> [texto...]
qrcode-generator config set NOME=PREFIXO
qrcode-generator config default NOME

# alias equivalente
qr [-c|--clear] [-p|--preset NOME] <texto> [texto...]
```

**Flags:**

| Flag                  | Descrição                                       |
| --------------------- | ----------------------------------------------- |
| `-c`, `--clear`       | Limpa o terminal antes de gerar os QR codes     |
| `-p`, `--preset NOME` | Usa o preset NOME como prefixo                   |
| `--no-preset`         | Gera sem prefixo, ignorando o preset default     |

**Exemplos:**

```sh
qr https://exemplo.com
qr -p zap 5511999999999
qr texto\ qualquer
qr -c 42
qrcode-generator --clear 1 2 3
```

A saída é um QR code renderizado em ASCII (half-blocks) impresso no stdout, um por texto informado.

## Desinstalação

```sh
rm "$(go env GOPATH)/bin/qrcode-generator" "$(go env GOPATH)/bin/qr"
```
