/*add robot account table*/
CREATE TABLE robot (
 id SERIAL PRIMARY KEY NOT NULL,
 name varchar(255),
 description varchar(1024),
 project_id int,
 expiresat bigint,
 disabled boolean DEFAULT false NOT NULL,
 creation_time timestamp default CURRENT_TIMESTAMP,
 update_time timestamp default CURRENT_TIMESTAMP,
 CONSTRAINT unique_robot UNIQUE (name, project_id)
);

CREATE TRIGGER robot_update_time_at_modtime BEFORE UPDATE ON robot FOR EACH ROW EXECUTE PROCEDURE update_update_time_at_column();

/*add master role*/
INSERT INTO role (role_code, name) VALUES ('DRWS', 'master');

/*delete replication jobs whose policy has been marked as "deleted"*/
DELETE FROM replication_job AS j
USING replication_policy AS p
WHERE j.policy_id = p.id AND p.deleted = TRUE;

/*delete replication policy which has been marked as "deleted"*/
DELETE FROM replication_policy AS p
WHERE p.deleted = TRUE;

CREATE TABLE registry (
 id SERIAL PRIMARY KEY NOT NULL,
 name varchar(256),
 url varchar(256),
 credential_type varchar(16),
 access_key varchar(128),
 access_secret varchar(1024),
 type varchar(32),
 insecure boolean,
 description varchar(1024),
 health varchar(16),
 creation_time timestamp default CURRENT_TIMESTAMP,
 update_time timestamp default CURRENT_TIMESTAMP,
 CONSTRAINT unique_registry_name UNIQUE (name)
);

CREATE TABLE "replication_policy_ng" (
  "id" SERIAL PRIMARY KEY NOT NULL,
  "name" varchar(256),
  "description" text,
  "creator" varchar(256),
  "src_registry_id" int4,
  "src_namespaces" varchar(256),
  "dest_registry_id" int4,
  "dest_namespace" varchar(256),
  "override" bool NOT NULL DEFAULT false,
  "enabled" bool NOT NULL DEFAULT true,
  "cron_str" varchar(256),
  "filters" varchar(1024),
  "replicate_deletion" bool NOT NULL DEFAULT false,
  "start_time" timestamp(6),
  "deleted" bool NOT NULL DEFAULT false,
  "creation_time" timestamp(6) DEFAULT now(),
  "update_time" timestamp(6) DEFAULT now(),
  CONSTRAINT unique_policy_ng_name UNIQUE ("name")
);