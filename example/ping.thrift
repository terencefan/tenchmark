struct Bar {
   1: optional i64 v1,
   2: optional string v2,
}

service Ping {

    void ping()

    void foo(1: i16 v_i16,
             2: bool v_bool,
             3: i32 v_i32,
             4: string v_str,
             5: list<i16> v_list,
             6: set<string> v_set,
             7: map<i64, double> v_map
             8: Bar v_st,
             9: map<string, Bar> v_st_map)
    void foo1(1: i16 v_i16)
    i16 foo2(1: i16 v_i16)
}
