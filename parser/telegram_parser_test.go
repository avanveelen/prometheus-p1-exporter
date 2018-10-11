package parser

import (
	"prometheus-p1-exporter/conn"
	"testing"
)

func TestParseESMR5Format(t *testing.T) {
	lines, readErr := conn.ReadTelegram(&conn.ESMR5TelegramReaderOptions, &conn.MockSource{})
	if readErr != nil {
		t.Fatal(readErr)
	}
	telegram, parseErr := ParseTelegram(&XS210ESMR5TelegramFormat, lines)
	if parseErr != nil {
		t.Fatal(parseErr)
	}
	if telegram.ElectricityUsageHigh != 1225.59 {
		t.Error(telegram.ElectricityUsageHigh)
	}
	if telegram.ElectricityUsageLow != 1179.186 {
		t.Error(telegram.ElectricityUsageLow)
	}
	if telegram.ElectricityReturnedHigh != 0.0 {
		t.Error(telegram.ElectricityReturnedHigh)
	}
	if telegram.ElectricityReturnedLow != 0.016 {
		t.Error(telegram.ElectricityReturnedLow)
	}
	if telegram.CurrentElectricityUsageHigh != 0.0 {
		t.Error(telegram.CurrentElectricityUsageHigh)
	}
	if telegram.CurrentElectricityUsageLow != 0.2 {
		t.Error(telegram.CurrentElectricityUsageLow)
	}
	if telegram.PowerFailuresLong != 2 {
		t.Error(telegram.PowerFailuresLong)
	}
	if telegram.PowerFailuresShort != 57 {
		t.Error(telegram.PowerFailuresShort)
	}
	if telegram.ActiveTariff != 2 {
		t.Error(telegram.ActiveTariff)
	}
	if telegram.GasUsage != 1019.003 {
		t.Error(telegram.GasUsage)
	}
	if telegram.Timestamp != 181009214805 {
		t.Error(telegram.Timestamp)
	}
}