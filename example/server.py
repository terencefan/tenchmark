#!/usr/bin/env python
# -*- coding: utf-8 -*-

# Author: stdrickforce (Tengyuan Fan)
# Email: <stdrickforce@gmail.com> <fantengyuan@baixing.com>

import thriftpy

from thriftpy.rpc import make_server

ping_thrift = thriftpy.load(
    "./ping.thrift", module_name='ping_thrift'
)


class Dispatcher(object):

    def __init__(self):
        self.logs = {}

    def ping(self):
        print("ping request")

    def foo(self, v_i16, v_bool, v_i32, v_str, v_list, v_set, v_map, v_st, v_st_map):
        print("foo: %d" % v_i16)
        # print(v_i16, v_bool, v_i32, v_str, v_list, v_set, v_map, v_st)
        # for key, val in v_st_map.items():
            # print(key, val)
    def foo1(self, v_i16):
        print("foo1: %d" % v_i16)
    def foo2(self, v_i16):
        print("foo2: %d" % v_i16)
        return 0

server = make_server(
    ping_thrift.Ping,
    Dispatcher(),
    "0.0.0.0",
    6000
)
print("server is listening on :6000")
server.serve()
