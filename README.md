# ftl/cabrillo

This little Go library to handle the [Cabrillo](https://wwrof.org/cabrillo/) file format for amateur radio contest log files.

## Use as a Go Library

To include `cabrillo` into your own projects as a library:

```shell
go get github.com/ftl/cabrillo
```
### Read a Cabrillo log file

```go
f, err := os.Open("mycabrillo.log")
if err != nil {
    panic(err)
}
log, err := cabrillo.Read(f)
if err != nil {
    panic(err)
}
```

### Write a Cabrillo log file

```go
log := cabrillo.NewLog()
log.Contest = "CQ-WW-CW"
log.Callsign = callsign.MustParse("DL0ABC")
log.Operators = []callsign.Callsign{callsign.MustParse("DL1ABC")}
log.Host = callsign.MustParse("DL1ABC")
log.Location = "DX"
log.Category.Operator = cabrillo.SingleOperator
log.Category.Assisted = cabrillo.Assisted
log.Category.Band = cabrillo.BandAll
log.Category.Power = cabrillo.HighPower
log.Category.Mode = cabrillo.ModeCW
log.Category.Transmitter = cabrillo.OneTransmitter
log.ClaimedScore = 12345
log.Club = "Bavarian Contest Club"
log.Name = "Hans Hamster"
log.Email = "hans.hamster@example.com"
log.Address.Text = "Beispielstraße 1"
log.Address.City = "Musterstadt"
log.Address.Postalcode = "12345"
log.Address.StateProvince = "Bavaria"
log.Address.Country = "Germany"
log.CreatedBy = "Golang Cabrillo Example"
log.Soapbox = "this is just an example that shows how to write Cabrillo logs in Golang"

// qsos is where you keep your QSO data in your internal representation
qsoData := make([]QSO, 0, len(qsos)) 
for _, qso := range qsos {
    // convertQSOToCabrillo converts your internal represenation to cabrillo.QSO
    qsoData = append(qsoData, convertQSOToCabrillo(qso)) 
}
log.QSOData = qsoData

f, err := os.Create("mycabrillo.log")
if err != nil {
    panic(err)
}
err = cabrillo.Write(f, log, false)
if err != nil {
    panic(err)
}
```

## License
This software is published under the [MIT License](https://www.tldrlegal.com/l/mit).

Copyright [Florian Thienel](http://thecodingflow.com/)
