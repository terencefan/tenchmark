import thriftpy
from thriftpy.rpc import make_server
import time

revenue_order_thrift = thriftpy.load(
        "./RevenueOrder.thrift",
    module_name='revenue_order_thrift')


class Dispatcher(object):
    def __init__(self):
        self.logs = {}

    def ping(self):
        print 'ping request'

    def createOrder(self, appId, tOrder, val_str, val_dou, val_list, val_set, val_map, val_i32, val_i64, val_byte, val_bool, val_spec):
        print appId
        print tOrder
        print val_str
        print val_dou
        print val_list
        print val_set
        print val_map
        print val_i32
        print val_i64
        print val_byte
        print val_bool
        print val_spec
        res = revenue_order_thrift.TCreateOrderResult()
        res.orderId = 11
        res.tieOrderIds = [1,2,2]
        return res

    def foo(self, v_i16, v_bool, v_i32, v_str, v_list, v_set, v_map, v_st, v_st_map):
        print v_i16, v_bool, v_i32, v_str, v_list, v_set, v_map, v_st
        for key, val in v_st_map.items():
            print key, val

server = make_server(revenue_order_thrift.RevenueOrder,
                     Dispatcher(), "0.0.0.0", 8888)
server.serve()
