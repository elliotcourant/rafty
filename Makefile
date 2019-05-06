.PHONY: protos

protos:
	@echo generating protos...
	@protoc -I=./ --go_out=plugins=grpc:./ ./raft.proto