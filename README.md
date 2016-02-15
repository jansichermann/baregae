Assumptions

If you use cloudsql, create an initial `migrations` table:
```
CREATE TABLE migrations (id int auto_increment primary key, migration varchar(255), rollback(255));
```
From here on out, do database changes only through migrations via the code.
