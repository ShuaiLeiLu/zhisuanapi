<div align="center">

![new-api](/web/default/public/logo.png)

# New API

🍥 **Next-Generation LLM Gateway and AI Asset Management System**

<p align="center">
  <a href="./README.zh_CN.md">简体中文</a> |
  <a href="./README.zh_TW.md">繁體中文</a> |
  <strong>English</strong> |
  <a href="./README.fr.md">Français</a> |
  <a href="./README.ja.md">日本語</a>
</p>

<p align="center">
  <a href="https://raw.githubusercontent.com/Calcium-Ion/new-api/main/LICENSE">
    <img src="https://img.shields.io/github/license/Calcium-Ion/new-api?color=brightgreen" alt="license">
  </a><!--
  --><a href="https://github.com/Calcium-Ion/new-api/releases/latest">
    <img src="https://img.shields.io/github/v/release/Calcium-Ion/new-api?color=brightgreen&include_prereleases" alt="release">
  </a><!--
  --><a href="https://hub.docker.com/r/CalciumIon/new-api">
    <img src="https://img.shields.io/badge/docker-dockerHub-blue" alt="docker">
  </a><!--
  --><a href="https://goreportcard.com/report/github.com/Calcium-Ion/new-api">
    <img src="https://goreportcard.com/badge/github.com/Calcium-Ion/new-api" alt="GoReportCard">
  </a>
</p>

<p align="center">
  <a href="https://trendshift.io/repositories/20180" target="_blank">
    <img src="https://trendshift.io/api/badge/repositories/20180" alt="QuantumNous%2Fnew-api | Trendshift" style="width: 250px; height: 55px;" width="250" height="55"/>
  </a>
  <br>
  <a href="https://hellogithub.com/repository/QuantumNous/new-api" target="_blank">
    <img src="https://api.hellogithub.com/v1/widgets/recommend.svg?rid=539ac4217e69431684ad4a0bab768811&claim_uid=tbFPfKIDHpc4TzR" alt="Featured｜HelloGitHub" style="width: 250px; height: 54px;" width="250" height="54" />
  </a><!--
  --><a href="https://www.producthunt.com/products/new-api/launches/new-api?embed=true&utm_source=badge-featured&utm_medium=badge&utm_campaign=badge-new-api" target="_blank" rel="noopener noreferrer">
    <img src="https://api.producthunt.com/widgets/embed-image/v1/featured.svg?post_id=1047693&theme=light&t=1769577875005" alt="New API - All-in-one AI asset management gateway. | Product Hunt" style="width: 250px; height: 54px;" width="250" height="54" />
  </a>
</p>

<p align="center">
  <a href="#quick-start">Quick Start</a> •
  <a href="#features">Features</a> •
  <a href="#deployment">Deployment</a> •
  <a href="#documentation">Documentation</a> •
  <a href="#license">License</a>
</p>

</div>

## Project Description

New API is an AI API gateway for unified model access, user management, usage analytics, billing, rate limiting, and private deployment.

This repository is modified from the open-source New API project. Original project: <https://github.com/QuantumNous/new-api>.

> [!IMPORTANT]
> This project is intended only for lawful and authorized AI API gateway, organization-level authentication, multi-model management, usage analytics, cost accounting, and private deployment scenarios. Users must lawfully obtain upstream API keys, accounts, model services, and interface permissions, and must comply with upstream terms of service and applicable laws and regulations.

## Quick Start

### Docker Compose

```bash
git clone https://github.com/QuantumNous/new-api.git
cd new-api
nano docker-compose.yml
docker-compose up -d
```

Open `http://localhost:3000` after the service starts.

### Docker

```bash
docker pull calciumion/new-api:latest

docker run --name new-api -d --restart always \
  -p 3000:3000 \
  -e TZ=Asia/Shanghai \
  -v ./data:/data \
  calciumion/new-api:latest
```

For MySQL/PostgreSQL, configure `SQL_DSN`. For Redis cache or multi-machine deployment, configure `REDIS_CONN_STRING`, `SESSION_SECRET`, and `CRYPTO_SECRET`.

> [!WARNING]
> When operating this project as a public generative AI service or API resale service, users should first complete all required filing, licensing, content safety, real-name verification, log retention, tax, payment, and upstream authorization obligations.

## Documentation

| Resource | Link |
| --- | --- |
| Official Documentation | [docs.newapi.pro](https://docs.newapi.pro/en/docs) |
| Deployment Guide | [Installation Documentation](https://docs.newapi.pro/en/docs/installation) |
| Environment Variables | [Environment Variables](https://docs.newapi.pro/en/docs/installation/config-maintenance/environment-variables) |
| API Documentation | [API Documentation](https://docs.newapi.pro/en/docs/api) |
| FAQ | [FAQ](https://docs.newapi.pro/en/docs/support/faq) |
| DeepWiki | [Ask DeepWiki](https://deepwiki.com/QuantumNous/new-api) |

## Features

- Unified API gateway for OpenAI-compatible, Claude, Gemini, Azure, AWS Bedrock, and other providers
- Modern dashboard with user, token, channel, quota, billing, and usage-log management
- Multi-language frontend: Simplified Chinese, Traditional Chinese, English, French, Japanese
- OpenAI Responses, Realtime, Claude Messages, Gemini, embeddings, rerank, image, audio, and task interfaces
- Weighted channel routing, automatic retry, rate limiting, and cache support
- Organization-level accounting, top-up, quota allocation, and flexible billing policies
- OAuth/OIDC login, passkeys, 2FA, user groups, and permission controls
- Compatible with SQLite, MySQL, PostgreSQL, Redis, and in-memory cache

## Model Support

New API supports mainstream chat, reasoning, embedding, rerank, image, audio, video, and custom upstream models through unified gateway interfaces.

See the [API Documentation - Gateway Interface](https://docs.newapi.pro/en/docs/api) for the complete model and interface list.

## Deployment

> **Latest Docker image:** `calciumion/new-api:latest`

| Component | Requirement |
| --- | --- |
| Local database | SQLite, with Docker volume mounted to `/data` |
| Remote database | MySQL >= 5.7.8 or PostgreSQL >= 9.6 |
| Cache | Redis recommended, in-memory cache supported |
| Container engine | Docker / Docker Compose |

Common environment variables:

| Variable | Purpose |
| --- | --- |
| `SQL_DSN` | Remote database connection string |
| `REDIS_CONN_STRING` | Redis connection string |
| `SESSION_SECRET` | Required for stable sessions in multi-machine deployments |
| `CRYPTO_SECRET` | Required when encrypted data is shared through Redis |
| `STREAMING_TIMEOUT` | Streaming timeout in seconds |
| `MAX_REQUEST_BODY_MB` | Max request body size after decompression |
| `PYROSCOPE_APP_NAME` | Pyroscope application name, default `new-api` |
| `HOSTNAME` | Hostname tag for Pyroscope, default `new-api` |

More deployment methods:

- Docker Compose: recommended for most deployments
- Docker command: suitable for lightweight single-node deployments
- BaoTa Panel: search for **New-API** in the application store, or see [BT tutorial](./docs/BT.md)

## Trusted Partners

<details>
<summary>View partners and special thanks</summary>

<p align="center">
  <em>No particular order</em>
</p>

<p align="center">
  <a href="https://www.cherry-ai.com/" target="_blank">
    <img src="./docs/images/cherry-studio.png" alt="Cherry Studio" height="80" />
  </a><!--
  --><a href="https://github.com/iOfficeAI/AionUi/" target="_blank">
    <img src="./docs/images/aionui.png" alt="Aion UI" height="80" />
  </a><!--
  --><a href="https://bda.pku.edu.cn/" target="_blank">
    <img src="./docs/images/pku.png" alt="Peking University" height="80" />
  </a><!--
  --><a href="https://www.compshare.cn/?ytag=GPU_yy_gh_newapi" target="_blank">
    <img src="./docs/images/ucloud.png" alt="UCloud" height="80" />
  </a><!--
  --><a href="https://www.aliyun.com/" target="_blank">
    <img src="./docs/images/aliyun.png" alt="Alibaba Cloud" height="80" />
  </a><!--
  --><a href="https://io.net/" target="_blank">
    <img src="./docs/images/io-net.png" alt="IO.NET" height="80" />
  </a>
</p>

<p align="center">
  <a href="https://www.jetbrains.com/?from=new-api" target="_blank">
    <img src="https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.png" alt="JetBrains Logo" width="120" />
  </a>
</p>

<p align="center">
  <strong>Thanks to <a href="https://www.jetbrains.com/?from=new-api">JetBrains</a> for providing free open-source development license for this project</strong>
</p>

</details>

## Related Projects

| Project | Description |
| --- | --- |
| [One API](https://github.com/songquanpeng/one-api) | Original project base |
| [Midjourney-Proxy](https://github.com/novicezk/midjourney-proxy) | Midjourney interface support |
| [new-api-key-tool](https://github.com/Calcium-Ion/new-api-key-tool) | Key quota query tool |
| [new-api-horizon](https://github.com/Calcium-Ion/new-api-horizon) | New API high-performance optimized version |

## Help Support

- Documentation: [Official Documentation](https://docs.newapi.pro/en/docs)
- Community: [Communication Channels](https://docs.newapi.pro/en/docs/support/community-interaction)
- Issues: [Issue Feedback](https://github.com/Calcium-Ion/new-api/issues)
- Releases: [Latest Release](https://github.com/Calcium-Ion/new-api/releases)

Contributions are welcome: bug reports, feature proposals, documentation improvements, and code changes.

## License

This project is licensed under the [GNU Affero General Public License v3.0 (AGPLv3)](./LICENSE).

Additional terms under AGPLv3 Section 7 apply. Modified versions must preserve
the author attribution notice `Frontend design and development by New API
contributors.` in the appropriate legal notices and in any prominent about,
legal, footer, or attribution location presented by the user interface.

Modified versions that present a user interface must also preserve a visible
link to the original project: <https://github.com/QuantumNous/new-api>.

This is an open-source project developed based on [One API](https://github.com/songquanpeng/one-api) (MIT License).

If your organization's policies do not permit the use of AGPLv3-licensed software, or if you wish to avoid the open-source obligations of AGPLv3, please contact us at: [support@quantumnous.com](mailto:support@quantumnous.com)

## Star History

<div align="center">

[![Star History Chart](https://api.star-history.com/svg?repos=Calcium-Ion/new-api&type=Date)](https://star-history.com/#Calcium-Ion/new-api&Date)

</div>

<div align="center">

### 💖 Thank you for using New API

If this project is helpful to you, welcome to give us a ⭐️ Star！

**[Official Documentation](https://docs.newapi.pro/en/docs)** • **[Issue Feedback](https://github.com/Calcium-Ion/new-api/issues)** • **[Latest Release](https://github.com/Calcium-Ion/new-api/releases)**

<sub>Built with ❤️ by QuantumNous</sub>

</div>
