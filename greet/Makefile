### @par Copyright: 
### Copyright (c) by Aruba Networks.
### 
### ALL RIGHTS RESERVED. These coded instructions and program statements are
### copyrighted works and confidential proprietary information of ScaleFlux, Inc.
### They may not be modified, copied, reproduced, distributed, or disclosed to
### third parties in any manner, medium, or form, in whole or in part.
### 

#
# Specify the top level directory.
#
TOP_DIR         := ..

# The following are defined by the top level makefile
EXTRA_CFLAGS    := -g
EXTRA_LDFLAGS   :=
SUBDIRS         := 
SHARE_LIB       := libfoo.so
SHARE_INCLUDE   :=

SRCS            := foo.c boo.c 
EXECS           := 
EXTRA_NODE_ITEMS:=

# Preprocessing Makefile
include $(TOP_DIR)/mk/Makefile.pre

$(TARGET_OBJ_DIR)/libfoo.so :  $(patsubst %.o,$(TARGET_OBJ_DIR)/%.o,$(OBJS))
	$(CC) $(CFLAGS) -shared $(LIB_DIR) $^ -o $@

protoc:
	protoc greetpb/greet.proto  --go_out=plugins=grpc:.

run_server:
	go run greet_server/server.go

run_client:
	go run greet_client/client.go
# Postprocessing Makefile
include $(TOP_DIR)/mk/Makefile.post
