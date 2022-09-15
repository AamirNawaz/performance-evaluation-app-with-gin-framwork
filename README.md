# performance-evaluation-app-with-go-gin-framework

performance evaluation app with golang and gofiber using mongodb

# /**\*\*\*\***\*\*\*\***\*\*\*\*** How to Setup application \***\*\*\*\*\*\*\***/

Step1:
go to the root directory of your project.

    Type cmd :
        go mod init performance_evaluation-app

Step2:

       Type cmd:
            go mod tidy

Step3:
Make sure to run the redis cache on your system.

    Lastly, start the Redis server like so:
    open ubuntu WSL type command in the root:
        type cmd:
                sudo service redis-server start

Connect to Redis
ou can test that your Redis server is running by connecting with the Redis CLI:

    type cmd:
        redis-cli
        127.0.0.1:6379> ping
        PONG

if you have windows operating system just install redis with WSL.
https://redis.io/docs/getting-started/installation/install-redis-on-windows/

Step4:
run the go server main.go file

    type cmd:
        go run main.go

OR for hot-reloading

    typ cmd:
        air main.go
