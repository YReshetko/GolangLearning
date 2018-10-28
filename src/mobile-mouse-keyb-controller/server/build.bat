@echo off

set ROOT_PATH=..\
set ABS_ROOT_PATH=
pushd %ROOT_PATH%
set PROJECT_PATH=%CD%
popd

call protoc -I %PROJECT_PATH% --go_out=plugins=grpc:%PROJECT_PATH% proto/controller.proto