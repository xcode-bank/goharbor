ALTER TABLE project ADD COLUMN IF NOT EXISTS registry_id int;
ALTER TABLE IF EXISTS cve_whitelist RENAME TO cve_allowlist;
UPDATE role SET name='maintainer' WHERE name='master';
UPDATE project_metadata SET name='reuse_sys_cve_allowlist' WHERE name='reuse_sys_cve_whitelist';

CREATE TABLE IF NOT EXISTS execution (
    id SERIAL NOT NULL,
    vendor_type varchar(16) NOT NULL,
    vendor_id int,
    status varchar(16),
    status_message text,
    trigger varchar(16) NOT NULL,
    extra_attrs JSON,
    start_time timestamp DEFAULT CURRENT_TIMESTAMP,
    end_time timestamp,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS task (
    id SERIAL PRIMARY KEY NOT NULL,
    execution_id int NOT NULL,
    job_id varchar(64),
    status varchar(16) NOT NULL,
    status_code int NOT NULL,
    status_revision int,
    status_message text,
    run_count int,
    extra_attrs JSON,
    creation_time timestamp DEFAULT CURRENT_TIMESTAMP,
    start_time timestamp,
    update_time timestamp,
    end_time timestamp,
    FOREIGN KEY (execution_id) REFERENCES execution(id)
);

ALTER TABLE blob ADD COLUMN IF NOT EXISTS update_time timestamp default CURRENT_TIMESTAMP;
ALTER TABLE blob ADD COLUMN IF NOT EXISTS status varchar(255);
ALTER TABLE blob ADD COLUMN IF NOT EXISTS version BIGINT default 0;
CREATE INDEX IF NOT EXISTS idx_status ON blob (status);
CREATE INDEX IF NOT EXISTS idx_version ON blob (version);

CREATE TABLE p2p_preheat_instance (
  id          SERIAL PRIMARY KEY NOT NULL,
  name        varchar(255) NOT NULL,
  description varchar(255),
  vendor	  varchar(255) NOT NULL,
  endpoint    varchar(255) NOT NULL,
  auth_mode   varchar(255),
  auth_data   text,
  enabled     boolean,
  is_default  boolean,
  insecure    boolean,
  setup_timestamp int,
  UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS p2p_preheat_policy (
    id SERIAL PRIMARY KEY NOT NULL,
    name varchar(255) NOT NULL,
    description varchar(1024),
    project_id int NOT NULL,
    provider_id int NOT NULL,
    filters varchar(1024),
    trigger varchar(255),
    enabled boolean,
    creation_time timestamp,
    update_time timestamp,
    UNIQUE (name, project_id)
);

ALTER TABLE schedule ADD COLUMN IF NOT EXISTS vendor_type varchar(16);
ALTER TABLE schedule ADD COLUMN IF NOT EXISTS vendor_id int;
ALTER TABLE schedule ADD COLUMN IF NOT EXISTS cron varchar(64);
ALTER TABLE schedule ADD COLUMN IF NOT EXISTS callback_func_name varchar(128);
ALTER TABLE schedule ADD COLUMN IF NOT EXISTS callback_func_param text;

/*abstract the cron, callback function parameters from table retention_policy*/
UPDATE schedule
SET vendor_type= 'RETENTION', vendor_id=retention.id, cron = retention.cron,
    callback_func_name = 'RETENTION', callback_func_param=concat('{"PolicyID":', retention.id, ',"Trigger":"Schedule"}')
FROM (
    SELECT id, data::json->'trigger'->'references'->>'job_id' AS schedule_id,
        data::json->'trigger'->'settings'->>'cron' AS cron
        FROM retention_policy
    ) AS retention
WHERE schedule.id=retention.schedule_id::int;

/*create new execution and task record for each schedule*/
DO $$
DECLARE
    sched RECORD;
    exec_id integer;
    status_code integer;
BEGIN
    FOR sched IN SELECT * FROM schedule
    LOOP
      INSERT INTO execution (vendor_type, vendor_id, trigger) VALUES ('SCHEDULER', sched.id, 'MANUAL') RETURNING id INTO exec_id;
      IF sched.status = 'Pending' THEN
        status_code = 0;
      ELSIF sched.status = 'Scheduled' THEN
        status_code = 1;
      ELSIF sched.status = 'Running' THEN
        status_code = 2;
      ELSIF sched.status = 'Stopped' OR sched.status = 'Error' OR sched.status = 'Success' THEN
        status_code = 3;
      ELSE
        status_code = 0;
      END IF;
      INSERT INTO task (execution_id, job_id, status, status_code, status_revision, run_count) VALUES (exec_id, sched.job_id, sched.status, status_code, 0, 0);
    END LOOP;
END $$;

ALTER TABLE schedule DROP COLUMN IF EXISTS job_id;
ALTER TABLE schedule DROP COLUMN IF EXISTS status;

UPDATE registry SET type = 'quay' WHERE type = 'quay-io';

CREATE TABLE IF NOT EXISTS data_migrations (
    version int
);
INSERT INTO data_migrations (version) VALUES (
    CASE
        /*if the "extra_attrs" isn't null, it means that the deployment upgrades from v2.0*/
        WHEN (SELECT Count(*) FROM artifact WHERE extra_attrs!='')>0 THEN 30
        ELSE 0
    END
);
ALTER TABLE schema_migrations DROP COLUMN IF EXISTS data_version;

ALTER TABLE artifact ADD COLUMN IF NOT EXISTS icon varchar(255);

UPDATE artifact
SET icon=(
CASE
    WHEN type='IMAGE' THEN
        'sha256:0048162a053eef4d4ce3fe7518615bef084403614f8bca43b40ae2e762e11e06'
    WHEN type='CHART' THEN
        'sha256:61cf3a178ff0f75bf08a25d96b75cf7355dc197749a9f128ed3ef34b0df05518'
    WHEN type='CNAB' THEN
        'sha256:089bdda265c14d8686111402c8ad629e8177a1ceb7dcd0f7f39b6480f623b3bd'
    ELSE
        'sha256:da834479c923584f4cbcdecc0dac61f32bef1d51e8aae598cf16bd154efab49f'
END);
