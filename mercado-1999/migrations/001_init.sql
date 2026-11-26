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