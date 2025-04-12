# counter-api

## Config

| Flag      | Environment Variable        | Default Value    | Notes
| :-------- | :-------------------------- | :--------------- | :-
| `db`      | `COUNTERAPI_DB_PATH`        | `database/`      | Path of the database, the process must have read/write permissions
| `backups` | `COUNTERAPI_BACKUPS_PATH`   | `backups/`       | Directory for database backups, the process must have read/write permissions
| `logs`    | `COUNTERAPI_LOGS_PATH`      | `logs/`          | Directory for log files, the process must have read/write permissions
| `address` | `COUNTERAPI_LISTEN_ADDRESS` | `127.0.0.1:8000` | Server address to listen on
| `timeout` | `COUNTERAPI_TIMEOUT`        | `15s`            | Server timeout
