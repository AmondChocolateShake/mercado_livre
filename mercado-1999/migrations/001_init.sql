-- TODO: items 테이블 설계 (internal/item/model.go의 Item 구조체 참고)
--   - 가격은 어떤 타입으로 저장할까? (FLOAT vs DECIMAL vs BIGINT)
--   - condition은 MySQL 예약어다. 컬럼 이름을 어떻게 할까?
--   - seller_id는 users.id를 참조해야 한다. (외래 키)
--   - 자주 조회되는 컬럼에 인덱스가 필요할까?
--   - 공개 ID가 MLA100000000처럼 보이게 하려면 AUTO_INCREMENT 시작값은?
--
-- TODO: orders 테이블 설계 (internal/order/model.go의 Order 구조체 참고)
--   - item_id, buyer_id, seller_id 외래 키
--   - 주문 당시 가격(unit_price)을 왜 items.price와 별도로 저장해야 할까?
--
-- 주의: 마이그레이션 실행기는 ; 기준으로 문장을 나누므로 주석 안에 세미콜론을 쓰지 마세요





CREATE TABLE IF NOT EXISTS users (
    id         BIGINT       NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id    VARCHAR(50)  NOT NULL,
    user_name   VARCHAR(50)  NOT NULL,
    nickname    VARCHAR(50)  NOT NULL,
    email      VARCHAR(255) NOT NULL,
    phone      VARCHAR(20)  NULL,

    created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME     NULL DEFAULT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uk_users_email (email),
    UNIQUE KEY uk_users_user_id (user_id),
    UNIQUE KEY uk_users_nickname (nickname)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;





CREATE TABLE IF NOT EXISTS items (
    id                  BIGINT          NOT NULL AUTO_INCREMENT PRIMARY KEY,
    site_id             varchar(20)     NOT NULL,
    seller_id           BIGINT          NOT NULL,
    title               VARCHAR(255)    NOT NULL,
    price               DECIMAL(15, 2)  NOT NULL,
    currency_id         VARCHAR(3)      NOT NULL,
    available_quantity  INT             NOT NULL,
    sold_quantity       INT             NOT NULL,
    buying_mode         VARCHAR(20)     NOT NULL,
    item_condition           VARCHAR(20)     NOT NULL,
    status              VARCHAR(20)     NOT NULL,
    created_at          DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at          DATETIME        NULL DEFAULT NULL,
    CONSTRAINT 
        fk_items_seller_id FOREIGN KEY (seller_id) REFERENCES users(id) 
        ON DELETE CASCADE
) AUTO_INCREMENT=100000000 ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;




CREATE TABLE IF NOT EXISTS orders (
    id                  BIGINT          NOT NULL AUTO_INCREMENT PRIMARY KEY,
    item_id             BIGINT          NOT NULL,
    buyer_id            BIGINT          NOT NULL,
    seller_id           BIGINT          NOT NULL,
    quantity            INT             NOT NULL,
    unit_price          DECIMAL(15, 2)  NOT NULL,
    total_amount        DECIMAL(15, 2)  NOT NULL,
    currency_id         VARCHAR(3)      NOT NULL,
    status              VARCHAR(20)     NOT NULL,
    created_at          DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at          DATETIME        NULL DEFAULT NULL,
    CONSTRAINT 
        fk_orders_item_id FOREIGN KEY (item_id) REFERENCES items(id),
    CONSTRAINT 
        fk_orders_buyer_id FOREIGN KEY (buyer_id) REFERENCES users(id),
    CONSTRAINT 
        fk_orders_seller_id FOREIGN KEY (seller_id) REFERENCES users(id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;