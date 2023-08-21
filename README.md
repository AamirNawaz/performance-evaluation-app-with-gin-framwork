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

first install ubuntu on your windows operating system with WSL.
https://apps.microsoft.com/store/detail/ubuntu/9PDXGNCFSCZV?hl=en-pk&gl=pk&rtc=1

after installation of susccessfully next step is install redis on ubuntu.

Goto the official docs of redis and follow the installaiton guide:
https://redis.io/docs/getting-started/installation/install-redis-on-windows/

After installation you can Make sure to run the redis cache on your system.

    installation command:
        curl -fsSL https://packages.redis.io/gpg | sudo gpg --dearmor -o /usr/share/keyrings/redis-archive-keyring.gpg

        echo "deb [signed-by=/usr/share/keyrings/redis-archive-keyring.gpg] https://packages.redis.io/deb $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/redis.list

        sudo apt-get update
        sudo apt-get install redis

    open ubuntu WSL type command in the root:
        type cmd:
                sudo service redis-server start

Connect to Redis
you can test that your Redis server is running by connecting with the Redis CLI:

    type cmd:
        > redis-cli (Enter)
        > 127.0.0.1:6379> ping (Enter)

        > out: PONG
    if you see pong on response it means it works perfectly.

Step4:
Make sure if you run on localhost also install mongodb as well.
https://www.mongodb.com/try/download/community

Then run the go server main.go file

    type cmd:
        go run main.go

OR for hot-reloading

    typ cmd:
        air main.go
