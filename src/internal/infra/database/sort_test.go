package database

import "testing"

func TestSortParseIgnoresUnknownFields(t *testing.T) {
	orders := SortParse("id,user_name:asc,sleep(1):desc", map[string]string{
		"id":        "id",
		"user_name": "user_name",
	})
	if len(orders) != 2 {
		t.Fatalf("orders len=%d, want 2", len(orders))
	}
	if orders[0].Field != "id" || orders[0].Asc {
		t.Fatalf("orders[0]=%+v, want id desc", orders[0])
	}
	if orders[1].Field != "user_name" || !orders[1].Asc {
		t.Fatalf("orders[1]=%+v, want user_name asc", orders[1])
	}
}

func TestPageNumSizeBounds(t *testing.T) {
	pageNum, pageSize := PageNumSize("0", "0")
	if pageNum != 0 || pageSize != 10 {
		t.Fatalf("default page=(%d,%d), want (0,10)", pageNum, pageSize)
	}
	pageNum, pageSize = PageNumSize("3", "1000")
	if pageNum != 2 || pageSize != 200 {
		t.Fatalf("max page=(%d,%d), want (2,200)", pageNum, pageSize)
	}
}
