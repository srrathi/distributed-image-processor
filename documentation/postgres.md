# Setting up PostgreSQL
## 2. Setting up PostgreSQL Locally
### **2.1 Using Docker**
To install and run PostgreSQL locally, we will run it in its Docker container through the following command:

```bash
docker run -d --name postgresCont -p 5432:5432 -e POSTGRES_PASSWORD=12345678 postgres
```

- **Container Name:** postgresCont
- **Port Mapping:** Binds the container's 5432 port with the local machine's 5432 port.
- **Password:** Sets up the default password for the user postgres.

### **2.2 Database Setup**
After running the container, execute the following command to enter the PostgreSQL shell:

```bash
docker exec -it postgresCont bash
```

Once inside the container run below command to enter postgres shell

```bash
psql -h localhost -U postgres
```

After it run below command to create a database in our postgres database

```sql
CREATE DATABASE ip_jobs;
```

### **2.3 User Setup**
Create a new PostgreSQL user named srrathi with an encrypted password:

```sql
CREATE USER srrathi WITH ENCRYPTED PASSWORD '12345678';
```

### **2.4 Granting Permissions**
Grant all privileges to the user srrathi for the ip_jobs database:

```sql
GRANT ALL PRIVILEGES ON DATABASE ip_jobs TO srrathi;
```

Exit the PostgreSQL shell:

```sql
\q
```

### **2.5 Wrapping Up**
By following these steps, you have set up a PostgreSQL database locally named ip_jobs with a user srrathi having necessary permissions.

Continue to the next steps for dumping CSV data, configuring endpoints, and creating Postman requests.

