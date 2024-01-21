# Setting up RabbitMQ
## 1. Installing RabbitMQ - Setup User & Virtual Host & Permissions

### **1.1 Using RabbitMQ Docker**
To set up RabbitMQ on your local machine, the easiest way is to run it using the official Docker image.

```bash
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.12-management
```

- **Port 5672:** Enables AMQP connections.
- **Port 15672:** Hosts the Admin UI and management UI.

### **1.2 Creating a User**
Use RabbitMQCLI to create a new user. The command syntax is add_user [username] [password].

```bash
docker exec rabbitmq rabbitmqctl add_user srrathi 12345678
```

For admin access, add the administrator tag to the new user.

```bash
docker exec rabbitmq rabbitmqctl set_user_tags srrathi administrator
```

Remove the default guest user for security.

```bash
docker exec rabbitmq rabbitmqctl delete_user guest
```

### **1.3 Virtual Host and Permissions**
Create a virtual host (vhost) using the add_vhost command.

```bash
docker exec rabbitmq rabbitmqctl add_vhost jobs
```

Add permissions to the user for the created vhost using the set_permissions command.

```bash
docker exec rabbitmq rabbitmqctl set_permissions -p jobs srrathi ".*" ".*" ".*"
```

### **1.4 Setting up Exchange**
Create an exchange named jobs_events within the vhost. Specify the vhost, username, and password of the administrator. Use the durable=true flag to persist restarts.

```bash
docker exec rabbitmq rabbitmqadmin declare exchange --vhost=jobs name=jobs_events type=topic -u srrathi -p 12345678 durable=true
```

Give the user permission to send on this exchange. Set permissions on a specific topic using the **set_topic_permissions** command.

```bash
docker exec rabbitmq rabbitmqctl set_topic_permissions -p jobs srrathi jobs_events "^jobs.*" "^jobs.*"

```

Restart RabbitMQ to apply changes.

```bash
docker restart rabbitmq
```

This completes the setup of RabbitMQ with user creation, virtual host, and exchange configuration. Move on to the next steps for setting up PostgreSQL, dumping CSV data, configuring endpoints, and Postman requests.