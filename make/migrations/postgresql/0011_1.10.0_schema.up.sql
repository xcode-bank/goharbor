/*Table for keeping the plug scanner registration*/
CREATE TABLE scanner_registration
(
    id SERIAL PRIMARY KEY NOT NULL,
    uuid VARCHAR(64) UNIQUE NOT NULL,
    url VARCHAR(256) UNIQUE NOT NULL,
    name VARCHAR(128) UNIQUE NOT NULL,
    description VARCHAR(1024) NULL,
    auth VARCHAR(16) NOT NULL,
    access_cred VARCHAR(512) NULL,
    adapter VARCHAR(128) NOT NULL,
    vendor VARCHAR(128) NOT NULL,
    version VARCHAR(32) NOT NULL,
    disabled BOOLEAN NOT NULL DEFAULT FALSE,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    skip_cert_verify BOOLEAN NOT NULL DEFAULT FALSE,
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

/*Table for keeping the scanner report. The report details are stored as JSONB*/
CREATE TABLE scanner_report
(
    id SERIAL PRIMARY KEY NOT NULL,
    digest VARCHAR(256) NOT NULL,
    registration_id VARCHAR(64) NOT NULL,
    job_id VARCHAR(32),
    status VARCHAR(16) NOT NULL,
    status_code INTEGER DEFAULT 0,
    report JSON,
    start_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    end_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(digest, registration_id)
)
