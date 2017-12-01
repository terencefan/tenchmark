namespace php RevenueOrder.Thrift

// 定义了ErrorCode的枚举类
enum ErrorCode {
    SYSTEM_ERROR = 0
}

// 系统错误
exception SystemException {
1: required ErrorCode code,
   2: required string name,
   3: optional string message,
}

// 业务逻辑产生的错误
exception CodeException {
1: required ErrorCode code,
   2: required string name,
   3: optional string message,
}

enum OrderType {
    NORMAL = 0
        MAIN = 1
        TIE = 2
}

enum OrderStatus {
    TOPAY = 1
        PAID = 3
        DEALT = 4
        DEALT_FAILED = 5
        REFUND = 6
        CLOSED = 8
}

enum OrderPaymentDirection {
    PAY = 1
        REFUND = 2
}

struct TOrder {
1: optional i64 id,
   2: optional i16 appId,
   3: optional i32 siteId,
   4: optional OrderType orderType,
   5: required i32 uid,
   6: optional i32 consumerUid,
   7: required i32 type,
   8: optional OrderStatus status,
   9: required i32 price,
   10: optional i32 originalPrice,
   11: optional i32 money,
   12: optional i32 startTime,
   13: optional i32 endTime,
   14: optional i32 createdTime,
   15: optional i32 modifiedTime,
   16: optional i32 paidTime,
   17: optional i32 refundedTime,

   // 扩展订单信息，中间空一些序号为了后续方便增加扩展
   31: optional map<string, string> orderItem,
   32: optional map<string, string> promotionInfo,
   33: optional list<map<string, string>> giftInfo,
   34: optional map<string, string> payInfo,
   35: optional map<string, string> salesInfo,
   36: optional map<string, string> dealInfo,
   37: optional string callback,

   // 搭售订单信息
   51: optional list<TOrder> tieOrders,
   52: optional list<i64> tieOrderIds,
}

struct TCreateOrderResult {
1: required i64 orderId,
       2: optional list<i64> tieOrderIds,
}

struct TOrderPayment {
1: required i32 amount,
       2: required i32 money,
       3: optional i16 paySrc,
       4: optional i32 paySrcSeq,
       5: optional OrderPaymentDirection direction,
       11: optional i32 createdTime,
       12: optional i32 modifiedTime,
}

// 描述这个服务提供了哪些接口
service RevenueOrder {

    void ping()
        throws (1: SystemException sys_exc,
                2: CodeException code_exc)

        void foo(1: i16 v_i16,
                2: bool v_bool,
                3: i32 v_i32,
                4: string v_str,
                5: list<i16> v_list,
                6: set<string> v_set,
                7: map<i64, double> v_map
                8: TOrder v_st,
                9: map<string, TOrder> v_st_map
                )
        throws (1: SystemException sys_exc,
                2: CodeException code_exc)

        TCreateOrderResult createOrder(1: i16 appId,
                2: TOrder tOrder,
                3: string val_str,
                4: double val_dou,
                5: list<i32> val_list,
                6: set<i32> val_set,
                7: map<i32, string> val_map,
                8: i32 val_i32,
                9: i64 val_i64,
                10: byte val_byte,
                11: bool val_bool,
                12: list<list<i32>> val_spec
                )
        throws (1: SystemException sys_exc,
                2: CodeException code_exc)

        bool paid(1: i16 appId,
                2: i32 siteId,
                3: i64 orderId,
                4: list<TOrderPayment> paymentList)
        throws (1: SystemException sys_exc,
                2: CodeException code_exc)

        bool dealt(1: i16 appId,
                2: i32 siteId,
                3: i64 orderId,
                4: bool success,
                5: map<string, string> dealInfo)
        throws (1: SystemException sys_exc,
                2: CodeException code_exc)

        bool refund(1: i16 appId,
                2: i32 siteId,
                3: i64 orderId)
        throws (1: SystemException sys_exc,
                2: CodeException code_exc)

        TOrder get(1: i16 appId,
                2: i32 siteId,
                3: i64 orderId)
        throws (1: SystemException sys_exc,
                2: CodeException code_exc)

        list<TOrder> getMulti(1: i16 appId,
                2: i32 siteId,
                3: list<i64> orderIds)
        throws (1: SystemException sys_exc,
                2: CodeException code_exc)

        list<TOrder> query(1: i16 appId,
                2: i32 siteId,
                3: i32 userId,
                4: i16 type,
                5: i32 start,
                6: i32 limit)
        throws (1: SystemException sys_exc,
                2: CodeException code_exc)

        i32 getFreeOrderCount(1: i16 appId,
                2: i32 siteId,
                3: i32 userId,
                4: i16 type)
        throws (1: SystemException sys_exc,
                2: CodeException code_exc)

        list<TOrder> getDealFailedOrders(1: i16 appId,
                2: i32 startTime,
                3: i32 endTime)
        throws (1: SystemException sys_exc,
                2: CodeException code_exc)
}
