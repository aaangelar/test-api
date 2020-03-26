### How to install on dev instance? ###
Execute the ff on cmd line
``` bash
docker run --rm --privileged -it --name test-api -v $PWD:/go/src/test-api -v $HOME/.ssh:/root/.ssh -w /go/src/test-api --env-file ./.env -p 8092:8092 billyteves/alpine-golang-glide:latest bash;
run-ssh;
make;
