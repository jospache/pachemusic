# API Golang Music Download YouTube 🎵

Aplicação web desenvolvida em **Go (Golang)** com frontend servido pelo próprio backend para busca e conversão autorizada de vídeos do YouTube em **MP3** ou **MP4**, utilizando [yt-dlp](https://github.com/yt-dlp/yt-dlp).

---

## 🚀 Funcionalidades

- Buscar vídeos do YouTube por URL ou nome.
- Baixar vídeos em formato **MP3 (áudio)** ou **MP4 (vídeo)**.
- Utilizar cookies para downloads de vídeos privados ou restritos.
- Servir o site e os downloads via **API REST**.
- Configuração simplificada com **Docker**.

---

## 📂 Estrutura do Projeto

```
.
├── cmd/            # Arquivo principal para inicialização da API
├── common/         # Utilitários e funções comuns
├── handlers/       # Handlers responsáveis pelas rotas da API
├── services/       # Serviços que implementam a lógica de download e conversão
├── utils/          # Funções auxiliares
├── Dockerfile      # Configuração para execução via Docker
├── cookies.txt     # Opcional; forneça em runtime, nunca publique cookies reais
├── go.mod          # Dependências Go
├── go.sum          # Hash das dependências Go
├── LICENSE         # Licença GPL-3.0
└── README.md       # Documentação
```

---

## ⚙️ Requisitos

Antes de executar a aplicação, instale as seguintes dependências:

- [Go 1.21+](https://go.dev/)
- [yt-dlp](https://github.com/yt-dlp/yt-dlp)
- [FFmpeg](https://ffmpeg.org/) (necessário para conversão de áudio/vídeo)

Verifique se estão disponíveis no **PATH**:

```bash
yt-dlp --version
ffmpeg -version
```

---

## ▶️ Como Executar

### 1. Clonar o Repositório

```bash
git clone https://github.com/marcosoleniuk/api-golang-music-download-youtube.git
cd api-golang-music-download-youtube
```

### 2. Instalar Dependências

```bash
go mod tidy
```

### 3. Rodar Localmente

```bash
go run cmd/main.go
```

A API estará disponível em:

```
http://localhost:8080
```

---

## 🐳 Executando com Docker

### 1. Build da Imagem

```bash
docker build -t youtube-music-api .
```

### 2. Rodar o Container

```bash
docker run --rm -p 8080:8080 -e API_KEY_YOUTUBE=sua_chave youtube-music-api
```

---

## 📡 Endpoints da API

### 🔹 `GET /search?q=TERM`

Busca vídeos no YouTube e retorna os resultados em JSON.

### 🔹 `GET /convert?youtubelink=YOUTUBE_URL&format=mp3|mp4`

Converte o vídeo e retorna o caminho do arquivo gerado.

### 🔹 `GET /download/:filename`

Baixa um arquivo convertido.

---

## 🔐 Uso de Cookies

Para vídeos restritos, monte um `cookies.txt` em runtime. O downloader funciona sem cookies para vídeos públicos.

```bash
docker run --rm -p 8080:8080 -e API_KEY_YOUTUBE=sua_chave \
	-v "${PWD}/cookies.txt:/app/cookies.txt:ro" youtube-music-api
```

Em plataformas como o Render, não publique o arquivo no GitHub. Exporte os cookies em formato Netscape,
converta o conteúdo para Base64 e adicione o resultado como variável secreta `YOUTUBE_COOKIES_B64`.
O backend cria o arquivo apenas dentro do container durante a conversão.

---

## 📜 Licença

Este projeto está sob a licença **GPL-3.0**.  
Você pode usar, modificar e distribuir, desde que mantenha os créditos e preserve a mesma licença.

---

## 🤝 Contribuições

Contribuições são bem-vindas!  
Sinta-se à vontade para abrir **issues** ou enviar **pull requests**.

---

## 🌟 Autor

Desenvolvido por [**Marcos Oleniuk**](https://github.com/marcosoleniuk) 🚀
