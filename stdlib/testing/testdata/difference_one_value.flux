import "testing"

inData = "
#datatype,string,long,dateTime:RFC3339,dateTime:RFC3339,dateTime:RFC3339,double,string,string,string,string
#group,false,false,false,false,false,false,true,true,true,true
#default,_result,,,,,,,,,
,result,table,_start,_stop,_time,_value,_field,_measurement,cpu,host
,,0,2018-05-22T19:53:24.421470485Z,2018-05-22T19:54:24.421470485Z,2018-05-22T19:53:26Z,1,usage_guest,cpu,cpu-total,host.local
"
outData = "
#datatype,string,long,dateTime:RFC3339,dateTime:RFC3339,dateTime:RFC3339,double,string,string,string,string
#group,false,false,false,false,false,false,true,true,true,true
#default,_result,0,,,,,usage_guest,cpu,cpu-total,host.local
,result,table,_start,_stop,_time,_value,_field,_measurement,cpu,host
"

difference_one_value = (table=<-) =>
  table
    |> range(start:2018-05-22T19:53:26Z)
    |> difference(nonNegative:true)

testFn = testing.test

testFn(name: "difference_one_value",
            input: testing.loadStorage(csv: inData),
            want: testing.loadMem(csv: outData),
            testFn: difference_one_value)
