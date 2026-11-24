# Mercado Libre: 1999 → 2026 — Go

남미 최대 이커머스 **Mercado Libre**의 핵심 마켓플레이스를 **두 시대**로 나누어 Go로 구현한 프로젝트입니다.

- **[`mercado-1999/`](./mercado-1999)** — 창업 초기. 아르헨티나에서 **경매 사이트**로 시작해 작은 팀이 빠르게 시장을 검증하던 시기의 **단순한 모듈러 모놀리스**
- **[`mercado-2026/`](./mercado-2026)** — 현재. 대규모 트래픽과 조직 규모에 맞춘 **분산 마이크로서비스 아키텍처**

같은 도메인(상품, 입찰, 주문, 질문, 평판)을 두 시대의 제약 조건으로 각각 설계하여,
**"서비스가 성장하면 아키텍처는 무엇을, 왜 바꿔야 하는가"**를 코드로 보여주는 것이 목표입니다.

> 1999년 버전은 당시 Mercado Libre의 실제 소스가 아니라, 초기 서비스 모델과 당시 스타트업의 제약(작은 팀, 적은 트래픽, 빠른 출시)을 가정해 재해석한 구현입니다.
> 도메인 명칭, 필드 값, ID 형식, URL 경로는 [Mercado Libre 공개 API](https://developers.mercadolibre.com)의 용어를 따릅니다.

---

## 한눈에 비교

| 구분 | 1999 — 초기 | 2026 — 현재 |
|---|---|---|
| **목표** | 빠른 출시, 시장 검증 | 확장성, 장애 격리, 팀 독립 배포 |
| **판매 방식** | 경매(`auction`) 중심 | 고정가(`buy_it_now`) 중심 + 경매 |
| **서비스 국가 (site)** | 아르헨티나(`MLA`) 하나 | `MLA`, `MLB`(브라질), `MLM`(멕시코) 등 다중 사이트 |
| **아키텍처** | 모듈러 모놀리스 (단일 바이너리) | 도메인별 마이크로서비스 |
| **서비스 간 통신** | 함수 호출 (같은 프로세스) | REST/gRPC (동기) + 이벤트 (비동기) |
| **데이터베이스** | 단일 MySQL | 서비스별 PostgreSQL 분리 (Database per Service) |
| **검색** | DB `LIKE` 쿼리 | 별도 `search` 서비스로 분리 |
| **캐시** | 없음 / 인메모리 | Redis |
| **메시지 큐** | 없음 (동기 처리) | Kafka |
| **동시성 제어 (입찰)** | DB 트랜잭션 + 행 잠금 | 분산 락 / 낙관적 잠금 + 이벤트 |
| **테스트** | `go test` 단위 테스트 | 단위 + 통합(testcontainers) + 부하 테스트(k6) |
| **배포** | 단일 서버 | Docker Compose (로컬), 컨테이너 기반 |
| **관측성** | 표준 로그 | 구조화 로그, Prometheus 지표, 분산 트레이싱 |

---

## 도메인 용어

Mercado Libre API에서 쓰는 이름을 그대로 패키지·서비스 이름으로 사용합니다.

| 도메인 | 의미 | API 리소스 |
|---|---|---|
| **site** | 국가 단위 사이트 (`MLA` = 아르헨티나) | `/sites/{site_id}` |
| **category** | 상품 분류 트리 | `/categories/{category_id}` |
| **user** | 회원. 판매자 평판(`seller_reputation`) 포함 | `/users/{user_id}` |
| **item** | 판매자가 올린 판매글. 가격, 상태, 판매 방식을 가짐 | `/items/{item_id}` |
| **bid** | 경매 입찰. "이 가격에 사겠다"는 제시 | `/items/{item_id}/bids` (직접 설계) |
| **order** | 낙찰 또는 즉시구매로 체결된 거래 | `/orders/{order_id}` |
| **question** | 구매 예정자의 공개 질문과 판매자 답변 | `/questions`, `/answers` |
| **feedback** | 거래(order) 종료 후 구매자·판매자 상호 평가 | `/orders/{order_id}/feedback` |

### 상품(product)이 아니라 판매글(item)인 이유
같은 "아이폰 15"라도 판매자 100명이 올리면 item은 100개입니다. 가격, 신품/중고 여부, 경매 여부, 질문, 입찰은 모두 물건 자체가 아니라 **판매글(item)에 붙는 정보**입니다.

### 경매(bid) 흐름
```text
판매자 A: 중고 아이폰 등록 (buying_mode: "auction", 시작가 500,000, 마감 3일 후)

구매자 B: 520,000 입찰   ← bid
구매자 C: 550,000 입찰   ← bid
구매자 B: 580,000 입찰   ← bid

마감 → 최고 입찰자 B 낙찰 → order 생성 → 거래 후 feedback 작성 → seller_reputation 반영
```

`bid` 도메인의 책임:
- **입찰 검증:** 현재 최고가보다 높은지, 판매자 본인 입찰이 아닌지, 경매가 진행 중인지
- **최고가 갱신:** item의 현재 가격 갱신
- **동시성 제어:** 같은 시점의 동시 입찰 중 하나만 성공해야 함 (이 프로젝트의 핵심 기술 포인트)
- **마감 처리:** 마감 시 최고 입찰자를 낙찰자로 확정하고 `order` 생성

### ID 형식
```text
item:     MLA123456789   (사이트 코드 + 숫자)
category: MLA1051
user:     123456789      (숫자)
```

### 주요 필드 값
```go
// item
Condition:  "new" | "used"
BuyingMode: "auction" | "buy_it_now"
Status:     "active" | "paused" | "closed"
CurrencyID: "ARS"

// question
Status: "UNANSWERED" | "ANSWERED"

// feedback
Rating: "positive" | "neutral" | "negative"
```

### API 경로
```text
GET  /sites/MLA/search?q=iphone        # 검색
GET  /sites/MLA/categories             # 카테고리 목록
POST /items                            # 상품 등록
GET  /items/{item_id}                  # 상품 조회
POST /items/{item_id}/bids             # 입찰
POST /questions                        # 질문 등록
POST /answers                          # 답변 등록
GET  /orders/{order_id}                # 주문 조회
POST /orders/{order_id}/feedback       # 거래 후기
GET  /users/{user_id}                  # 회원 조회 (seller_reputation 포함)
```

### 범위 밖
- **결제(Mercado Pago):** 2003년 출시. 1999 버전에서는 구매자·판매자가 직접 거래한다고 가정합니다.
- **배송(Mercado Envios):** 두 버전 모두 다루지 않습니다.

---

## 1999 — 초기 마켓플레이스

### 설계 철학
- **경매 중심:** 1999년 Mercado Libre는 경매 사이트로 시작했습니다. `bid` 도메인과 경매 마감 처리를 핵심으로 구현합니다.
- **핵심 마켓플레이스 우선:** 결제 게이트웨이나 자체 물류 없이도 성립하는 **매칭·거래 의도·신뢰(평판)**를 우선 구현합니다.
- **거래 신뢰도(Reputation):** 온라인 거래 초기의 낮은 신뢰를 보완하기 위해 `feedback`과 `seller_reputation` 집계를 핵심 도메인으로 다룹니다.
- **실용적인 모놀리스:** 변경 속도를 높이기 위해 하나의 바이너리, 하나의 DB로 구성하되 도메인 경계는 패키지로 명확히 나눕니다.

### 기술 스택
- **Language:** Go `1.2x`
- **HTTP / Routing:** Gin
- **Database:** MySQL
- **Cache / Queue:** 사용하지 않음 (필요 시 인메모리)
- **Test:** `go test`

### 디렉터리 구조
```text
mercado-1999/
├── cmd/
│   └── api/
│       └── main.go           # 엔트리포인트, 의존성 조립, 서버 실행
├── internal/
│   ├── config/               # 환경변수 로딩 및 설정 구조체
│   ├── database/             # DB 연결 및 마이그레이션
│   ├── site/                 # 사이트 (MLA 하나만)
│   ├── category/             # 카테고리 트리
│   ├── user/                 # 회원, seller_reputation 계산
│   ├── item/                 # 판매글 등록/조회/검색/상태 관리
│   ├── bid/                  # 경매 입찰, 동시성 제어, 마감 처리
│   ├── order/                # 낙찰/즉시구매 거래 체결
│   ├── question/             # 상품 질문과 답변
│   └── feedback/             # 거래 후기
├── migrations/               # SQL 스키마 마이그레이션 파일
├── .env.example
├── go.mod
└── go.sum
```

각 도메인 패키지는 `handler.go`(HTTP) → `service.go`(비즈니스 규칙) → `repository.go`(DB) → `model.go`(구조체) 순서로 책임을 나눕니다.

---

## 2026 — 현대 마켓플레이스

### 설계 철학
- **서비스 단위 독립성:** 도메인별로 서비스와 DB를 분리해 팀이 독립적으로 개발·배포할 수 있게 합니다.
- **비동기 이벤트 기반:** 낙찰 → 주문 → 평판 갱신 같은 후속 처리를 이벤트로 연결해 결합도를 낮추고 장애 전파를 막습니다. Kafka 토픽 이름은 Mercado Libre 알림 토픽(`items`, `questions`, `orders_v2`)을 따릅니다.
- **읽기 성능 최적화:** 트래픽이 가장 큰 검색을 `search` 서비스로 분리하고, 상품 조회는 Redis 캐시로 DB 부하를 줄입니다.
- **판매 방식 변화 반영:** 시장이 고정가(`buy_it_now`) 중심으로 바뀐 흐름을 반영하되, 인기 경매 마감 직전의 입찰 폭주는 분산 환경에서 처리합니다.
- **다중 사이트:** `MLA` 외에 `MLB`, `MLM` 등 여러 국가 사이트를 지원합니다.
- **측정 가능한 개선:** 1999 버전과 같은 시나리오로 부하 테스트를 돌려 성능 차이를 수치로 기록합니다.

### 기술 스택
- **Language:** Go `1.2x`
- **HTTP / RPC:** Gin (외부 API), gRPC (서비스 간 통신)
- **Database:** PostgreSQL (서비스별 분리)
- **Cache:** Redis
- **Message Queue:** Kafka
- **Test:** `go test` + testify, testcontainers-go, k6
- **Infra / Observability:** Docker Compose, Prometheus, OpenTelemetry

### 디렉터리 구조
```text
mercado-2026/
├── services/
│   ├── users/                # 회원, seller_reputation (feedback 이벤트 구독)
│   ├── items/                # 판매글 등록/조회 (Redis 캐시)
│   ├── search/               # 검색 (items 이벤트로 인덱스 갱신)
│   ├── categories/           # 카테고리 트리
│   ├── bids/                 # 경매 입찰, 분산 동시성 제어
│   ├── orders/               # 거래 체결
│   ├── questions/            # 상품 질문과 답변
│   ├── feedback/             # 거래 후기 (orders 이벤트 구독)
│   └── notifications/        # 이벤트 → 판매자/구매자 알림
├── pkg/                      # 공통 라이브러리 (로깅, 설정, 이벤트 스키마)
├── proto/                    # gRPC 정의
├── deploy/
│   └── docker-compose.yml    # PostgreSQL, Redis, Kafka, 전체 서비스 실행
├── loadtest/                 # k6 부하 테스트 시나리오
└── go.work                   # 서비스별 go module을 묶는 워크스페이스
```

---

## 진화 과정에서 바뀐 것과 이유

| 문제 상황 (가정) | 1999의 한계 | 2026의 해결 |
|---|---|---|
| 상품 조회 트래픽 급증 | 모든 요청이 DB로 직행 | Redis 캐시로 읽기 부하 분산 |
| 검색 트래픽 증가 | `LIKE` 쿼리가 메인 DB를 점유 | `search` 서비스 분리, 이벤트로 인덱스 갱신 |
| 인기 경매 마감 직전 입찰 집중 | 행 잠금 경합으로 응답 지연 | 잠금 범위 축소 + 이벤트 기반 후속 처리 |
| 평판 갱신 로직 장애 | 거래 체결 요청까지 함께 실패 | 이벤트 비동기 처리로 장애 격리 |
| 다른 국가로 서비스 확장 | `MLA` 하나만 가정한 구조 | `site_id` 기반 다중 사이트 |
| 서비스별 데이터 요구 차이 | 단일 MySQL에 모든 도메인 | 서비스별 PostgreSQL로 분리 |
| 팀과 기능 증가 | 하나의 코드베이스에서 배포 충돌 | 서비스별 독립 배포 |

> 각 항목의 측정 결과(응답 시간, 처리량)는 부하 테스트 후 이 섹션에 추가합니다.

---

## 실행 방법

```bash
# 1999 — 단일 서버
cd mercado-1999
cp .env.example .env
go run ./cmd/api

# 2026 — 전체 서비스
cd mercado-2026
docker compose -f deploy/docker-compose.yml up
```
