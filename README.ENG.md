# Telgram verification bot
[中国人](README.md)

### Play With Docker

1. **Deployment environment, click `Try in PWD` button**

[![Try in PWD](https://github.com/play-with-docker/stacks/raw/master/assets/images/button.png)](https://labs.play-with-docker.com/?stack=https://raw.githubusercontent.com/jqs7/zwei/master/stack.yml)

2. **Open the Docker Swarm management interface**
![](images/open.jpg)
    
3. **Login to the Docker Swarm management interface**
![](images/login.jpg)

4. **Select the verification robot Docker Swarm service -- pwd_zwei**
![](images/select.jpg)

5. **EDIT VERIFICATION BOT SERVICE**
![](images/edit.jpg)

6. **Modify the environment variable `ZWEI_TOKEN` in the verification service to your Telegram bot token, and save**
![](images/modify.jpg)



### Quick Start for development

```bash
# Postgres database installation needs to install docker for mac/win first
make all
```
