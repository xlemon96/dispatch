CREATE TABLE `task`
(
    `id`	        bigint	      NOT NULL AUTO_INCREMENT,
    `task_state`	varchar(32)	  NOT NULL,
    `task_type`		varchar(32)	  NOT NULL,
    `created_time`	datetime      NOT NULL,
    `update_time`	datetime      NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE `worker`
(
    `id`	        bigint	      NOT NULL AUTO_INCREMENT,
    `name`			varchar(32)	  NOT NULL,
    `host_ip`		varchar(32)	  NOT NULL,
    `port`			varchar(32)	  NOT NULL,
    `worker_state`	varchar(32)	  NOT NULL,
    `disable`		tinyint	  	  NOT NULL,
    `created_time`	datetime      NOT NULL,
    `update_time`	datetime      NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE `dag_instance`
(
    `id`	        bigint	      NOT NULL AUTO_INCREMENT,
    `task_id`		bigint	      NOT NULL,
    `depends`		varchar(256)  NOT NULL,
    `host_ip`		varchar(32)	  NOT NULL,
    `port`			varchar(32)	  NOT NULL,
    `dag_type`	    varchar(32)	  NOT NULL,
    `output`	    varchar(256)  NOT NULL,
    `dag_state`	    varchar(32)	  NOT NULL,
    `created_time`	datetime      NOT NULL,
    `update_time`	datetime      NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;