# README

###
```
USE master
CREATE DATABASE xo
GO
```

```
CREATE TABLE [user] (
    id INT PRIMARY KEY IDENTITY (1, 1),
    subject VARCHAR (256) NOT NULL,
    created_date DATETIME DEFAULT GETDATE(),
    changed_date DATETIME DEFAULT GETDATE(),
    deleted_date DATETIME
);
```
