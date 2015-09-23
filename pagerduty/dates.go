package pagerduty

import (
	"fmt"
	"time"
)

const dateFormat = `"2006-01-02"`

// Date is a thin wrapper around the time.Time type that parses dates of the
// format 'YYYY-MM-DD'.
type Date struct {
	*time.Time
}

func (d Date) String() string {
	return d.Time.Format(dateFormat)
}

func (d *Date) MarshalJSON() ([]byte, error) {
	str := fmt.Sprintf("%q", d.Time.Format(dateFormat))
	return []byte(str), nil
}

func (d *Date) UnmarshalJSON(data []byte) error {
	date, err := time.Parse(dateFormat, string(data))
	if err != nil {
		return err
	}

	d.Time = &date
	return nil
}
