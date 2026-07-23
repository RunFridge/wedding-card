[English](../README.md) | **한국어**

# 웨딩 청첩장 웹사이트

[![Release](https://img.shields.io/github/v/release/RunFridge/wedding-card)](https://github.com/RunFridge/wedding-card/releases)
[![License](https://img.shields.io/github/license/RunFridge/wedding-card)](../LICENSE)
[![GitHub](https://img.shields.io/badge/GitHub-RunFridge%2Fwedding--card-181717?logo=github)](https://github.com/RunFridge/wedding-card)

스타듀 밸리 감성의 픽셀 아트로 꾸민, 직접 호스팅할 수 있는 모바일 청첩장 웹사이트입니다. 하객에게는 방명록, 미니게임, 사진 공유가 있는 인터랙티브 청첩장을, 신랑 신부에게는 코드를 건드리지 않고 모든 것을 관리할 수 있는 관리자 패널을 제공합니다.

도커 컨테이너 하나로 모든 것이 동작합니다. 클라우드 계정도, 월 구독료도, 필수 외부 서비스도 없습니다.

> **🌐 라이브 데모:** [demo-wedding.runfridge.dev](https://demo-wedding.runfridge.dev/) — 직접 호스팅하기 전에 먼저 체험해 보세요.
> ⚠️ 데모 서버이므로 언제든 중단되거나 데이터가 초기화될 수 있습니다.

## 하객이 만나는 것들

- **청첩장** — 신랑 신부와 양가 가족 소개, 예식 일시와 장소, 마음 전하실 곳(계좌 안내)
- **오시는 길** — 카카오맵 · 네이버 지도 · 티맵 · 구글 지도 링크와 지하철 · 버스 · 주차 · 전세버스 안내
- **방명록** — 하객이 직접 정한 비밀번호로 수정 · 삭제할 수 있는 축하 메시지. 신랑 신부만 볼 수 있는 비밀 메시지로 남길 수도 있습니다
- **카드 짝 맞추기 게임** — 웨딩 사진으로 즐기는 미니게임. 오락실 스타일 순위표, 수집형 업적, 명예의 전당까지
- **사진 공유** — 하객이 결혼식 당일 사진을 올릴 수 있습니다. 자르기 · 회전 · 필터 편집기가 내장되어 있고, 예식 직전부터 카운트다운과 함께 업로드가 열립니다
- **소소한 즐거움** — 낮/밤 테마, 닭과 병아리가 돌아다니는 픽셀 아트 배경 애니메이션, 색종이 축하 효과, 모바일 햅틱 피드백
- **심플 청첩장** — 화려한 화면이 부담스러운 하객을 위한 `/simple` 페이지. 큰 글씨로 꼭 필요한 정보만 담백하게 보여주며, 원하면 외부 청첩장 주소로 연결할 수도 있습니다

## 신랑 신부가 얻는 것들

`/-/admin/` 관리자 패널에서 모든 것을 직접 관리할 수 있습니다:

- **예식 정보** — 이름, 가족 소개, 일시, 예식장, 인사말, 계좌번호, 오시는 길을 모두 브라우저에서 수정합니다. 링크를 보낸 뒤에도 언제든지요
- **사진 관리** — 게임과 갤러리에 쓰일 사진을 고르고, 메인 사진을 지정하고, 하객이 올린 사진을 승인하거나 숨깁니다
- **방명록 관리** — 메시지를 숨기거나 삭제하고, 비밀 메시지를 읽어봅니다
- **자동 검열 (선택)** — OpenAI 검열 서비스가 하객 사진과 메시지를 자동으로 심사합니다. 언제든 직접 번복할 수 있습니다
- **통계** — 일별 방문자 차트와 실시간 서버 로그
- **순위표 관리** — 장난 기록을 지우거나 게임 순위를 초기화합니다

## 시작하기

[Docker](https://docs.docker.com/get-docker/)가 설치된 서버가 필요합니다. 작은 VPS나 홈 서버면 충분합니다.

**1. 실행**

```bash
docker run -d --name wedding-server \
  -p 8080:8080 \
  -v wedding-data:/data \
  -e TZ=Asia/Seoul \
  --restart unless-stopped \
  ghcr.io/runfridge/wedding-card:latest
```

예식이 한국이 아니라면 `Asia/Seoul`을 해당 시간대로 바꾸세요. Docker Compose를 선호한다면 이 저장소의 [`docker-compose.yml`](../docker-compose.yml)을 내려받아 `docker compose up -d`를 실행하면 됩니다.

**2. 관리자 비밀번호 확인**

첫 실행 시 비밀번호가 자동으로 생성되어 서버 로그에 출력됩니다:

```bash
docker logs wedding-server
```

`/data/admin_password.txt` 파일로도 저장됩니다 (직접 비밀번호를 설정하면 삭제됩니다).

**3. 초기 설정**

`http://localhost:8080/-/admin/`에 접속해 로그인하세요. 설정 마법사가 새 비밀번호 설정, 사진 저장소(기본값은 서버에 저장), 선택 사항인 콘텐츠 검열을 차례로 안내합니다.

**4. 나만의 청첩장 만들기**

**설정**에서 예식 정보를 채우고, 사진을 올리고, 링크를 하객들에게 공유하세요. 🎉

> **팁:** 도메인과 HTTPS로 서비스하려면 8080 포트 앞에 리버스 프록시(Caddy, Nginx, Traefik 등)를 두세요.

## 데모 체험하기

라이브 데모가 [demo-wedding.runfridge.dev](https://demo-wedding.runfridge.dev/)에서 운영 중이며, 관리자 패널은 [demo-wedding.runfridge.dev/-/admin/](https://demo-wedding.runfridge.dev/-/admin/)에서 비밀번호 `demo_1234!`로 로그인할 수 있습니다. 방명록을 남기고, 게임을 해보고, 사진을 승인하고, 예식 정보를 바꿔보세요. 단, 남기고 싶은 건 남기지 마세요 — **모든 데모 데이터는 매주 초기화됩니다**.

직접 데모 인스턴스를 운영하고 싶다면 개발 가이드의 [Demo Mode](DEVELOPMENT.md#demo-mode)(영문)를 참고하세요.

## 데이터와 백업

모든 데이터는 `wedding-data` 도커 볼륨 안에 있습니다. 데이터베이스와, 기본 로컬 저장소를 쓴다면 업로드된 사진까지 전부요. 이 볼륨 하나만 백업하면 청첩장 사이트 전체가 백업됩니다.

## 콘텐츠 검열 (선택)

하객 사진을 일일이 승인하기 번거롭다면 **관리자 → 시스템**(또는 초기 설정 마법사)에서 OpenAI API 키로 검열을 켜세요. 문제없는 사진은 자동 승인되고, 부적절한 방명록 메시지는 자동으로 숨겨지며, 모든 결정은 관리자 패널에서 번복할 수 있습니다.

## 오프라인 / 폐쇄망 설치

외부 레지스트리에 접근할 수 없는 서버를 위해 각 [GitHub 릴리스](https://github.com/RunFridge/wedding-card/releases)에 이미지 압축 파일이 첨부됩니다:

```bash
docker load < wedding-server-<version>-amd64.tar.gz   # 또는 -arm64
```

## 언어

청첩장과 관리자 패널 모두 한국어와 영어를 지원합니다. 언어는 방문자의 브라우저 설정을 따르며 언제든 바꿀 수 있고, 직접 입력한 내용(인사말, 안내 문구, 오시는 길 등)은 작성한 그대로 표시됩니다.

## 개발자를 위해

소스 빌드, API 탐색, 언어 추가, 기여에 관심이 있다면 [개발 가이드](DEVELOPMENT.md)(영문)를 확인하세요.

## 라이선스

[MIT](../LICENSE)
