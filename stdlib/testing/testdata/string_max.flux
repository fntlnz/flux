import "testing"

inData = "
#datatype,string,long,dateTime:RFC3339,dateTime:RFC3339,dateTime:RFC3339,string,string,string,string,string,string,string
#group,false,false,false,false,false,false,true,true,true,true,true,true
#default,_result,,,,,,,,,,,
,result,table,_start,_stop,_time,_value,_field,_measurement,device,fstype,host,path
,,0,2018-05-22T19:53:24.421470485Z,2018-05-22T19:54:24.421470485Z,2018-05-22T19:54:16Z,13F2,used_percent,disk,disk1,apfs,host.local,/
,,0,2018-05-22T19:53:24.421470485Z,2018-05-22T19:54:24.421470485Z,2018-05-22T19:53:56Z,2COTDe,used_percent,disk,disk1,apfs,host.local,/
,,0,2018-05-22T19:53:24.421470485Z,2018-05-22T19:54:24.421470485Z,2018-05-22T19:54:06Z,cLnSkNMI,used_percent,disk,disk1,apfs,host.local,/
,,0,2018-05-22T19:53:24.421470485Z,2018-05-22T19:54:24.421470485Z,2018-05-22T19:53:26Z,a,used_percent,disk,disk1,apfs,host.local,/
,,0,2018-05-22T19:53:24.421470485Z,2018-05-22T19:54:24.421470485Z,2018-05-22T19:53:46Z,b,used_percent,disk,disk1,apfs,host.local,/
,,0,2018-05-22T19:53:24.421470485Z,2018-05-22T19:54:24.421470485Z,2018-05-22T19:53:36Z,k9ngm,used_percent,disk,disk1,apfs,host.local,/
"
outData = "error: invalid use of function: *functions.MaxSelector has no implementation for type string
"

t_string_max = (table=<-) =>
  table
    |> range(start:2018-05-22T19:54:16Z)
    |> max()
testFn = testing.test

testFn(name: "string_max",
            input: testing.loadStorage(csv: inData),
            want: testing.loadMem(csv: outData),
            testFn: t_string_max)