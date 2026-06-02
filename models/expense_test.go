package models

import (
	"errors"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"
	"time"

	expenseutils "expense-tracker-api/utils/expenses"

	"github.com/beego/beego/v2/server/web"
)

func setupModelStorage(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	if err := web.AppConfig.Set("users_csv_path", filepath.Join(dir, "users.csv")); err != nil {
		t.Fatalf("set users csv path: %v", err)
	}
	if err := web.AppConfig.Set("expenses_csv_path", filepath.Join(dir, "expenses.csv")); err != nil {
		t.Fatalf("set expenses csv path: %v", err)
	}
}

func sampleExpenses() []Expense {
	return []Expense{
		{ID: 1, UserID: 1, Title: "Lunch", Amount: 350.50, Category: "Food", ExpenseDate: "2025-06-10"},
		{ID: 2, UserID: 1, Title: "Bus", Amount: 50, Category: "Transport", ExpenseDate: "2025-06-09"},
		{ID: 3, UserID: 2, Title: "Book", Amount: 200, Category: "Education", ExpenseDate: "2025-06-11"},
	}
}

func TestFilterAndSortExpenses(t *testing.T) {
	tests := []struct {
		name    string
		options ExpenseListOptions
		wantIDs []int
	}{
		{
			name:    "filter category",
			options: ExpenseListOptions{Category: "Food"},
			wantIDs: []int{1},
		},
		{
			name:    "filter date range",
			options: ExpenseListOptions{DateFrom: "2025-06-10", DateTo: "2025-06-11"},
			wantIDs: []int{1, 3},
		},
		{
			name:    "sort amount desc",
			options: ExpenseListOptions{SortBy: "amount", SortOrder: "desc"},
			wantIDs: []int{1, 3, 2},
		},
		{
			name:    "sort date asc",
			options: ExpenseListOptions{SortBy: "expense_date", SortOrder: "asc"},
			wantIDs: []int{2, 1, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterAndSortExpenses(sampleExpenses(), tt.options)
			gotIDs := make([]int, 0, len(got))
			for _, expense := range got {
				gotIDs = append(gotIDs, expense.ID)
			}
			if !reflect.DeepEqual(gotIDs, tt.wantIDs) {
				t.Fatalf("ids = %v, want %v", gotIDs, tt.wantIDs)
			}
		})
	}
}

func TestSummarizeExpenses(t *testing.T) {
	tests := []struct {
		name       string
		dateFrom   string
		dateTo     string
		wantAmount float64
		wantCount  int
	}{
		{name: "all expenses", wantAmount: 600.50, wantCount: 3},
		{name: "date range", dateFrom: "2025-06-10", dateTo: "2025-06-11", wantAmount: 550.50, wantCount: 2},
		{name: "empty result", dateFrom: "2025-07-01", dateTo: "2025-07-31", wantAmount: 0, wantCount: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SummarizeExpenses(sampleExpenses(), tt.dateFrom, tt.dateTo)
			if got.TotalAmount != tt.wantAmount {
				t.Fatalf("total amount = %v, want %v", got.TotalAmount, tt.wantAmount)
			}
			if got.TotalCount != tt.wantCount {
				t.Fatalf("total count = %d, want %d", got.TotalCount, tt.wantCount)
			}
		})
	}
}

func TestExpenseStorageOperations(t *testing.T) {
	tests := []struct {
		name    string
		run     func(t *testing.T)
		wantErr error
	}{
		{
			name: "create and get expense",
			run: func(t *testing.T) {
				expense := Expense{UserID: 1, Title: "Lunch", Amount: 350.50, Category: "Food", ExpenseDate: "2025-06-10"}
				if err := CreateExpense(&expense); err != nil {
					t.Fatalf("create expense: %v", err)
				}
				got, err := GetExpenseByID(expense.ID, expense.UserID)
				if err != nil {
					t.Fatalf("get expense: %v", err)
				}
				if got.Title != expense.Title {
					t.Fatalf("title = %q, want %q", got.Title, expense.Title)
				}
			},
		},
		{
			name: "update expense",
			run: func(t *testing.T) {
				expense := Expense{UserID: 1, Title: "Lunch", Amount: 350.50, Category: "Food", ExpenseDate: "2025-06-10"}
				if err := CreateExpense(&expense); err != nil {
					t.Fatalf("create expense: %v", err)
				}
				expense.Title = "Dinner"
				expense.Amount = 500
				if err := UpdateExpense(&expense); err != nil {
					t.Fatalf("update expense: %v", err)
				}
				got, err := GetExpenseByID(expense.ID, expense.UserID)
				if err != nil {
					t.Fatalf("get expense: %v", err)
				}
				if got.Title != "Dinner" || got.Amount != 500 {
					t.Fatalf("updated expense = %#v", got)
				}
			},
		},
		{
			name: "delete expense",
			run: func(t *testing.T) {
				expense := Expense{UserID: 1, Title: "Lunch", Amount: 350.50, Category: "Food", ExpenseDate: "2025-06-10"}
				if err := CreateExpense(&expense); err != nil {
					t.Fatalf("create expense: %v", err)
				}
				if err := DeleteExpense(expense.ID, expense.UserID); err != nil {
					t.Fatalf("delete expense: %v", err)
				}
				if _, err := GetExpenseByID(expense.ID, expense.UserID); !errors.Is(err, ErrExpenseNotFound) {
					t.Fatalf("get deleted err = %v, want %v", err, ErrExpenseNotFound)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupModelStorage(t)
			tt.run(t)
		})
	}
}

func TestExpenseReadHelpers(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "get all skips invalid rows",
			run: func(t *testing.T) {
				records := [][]string{
					{"bad", "1", "Lunch", "350.50", "Food", "", "2025-06-10", time.Now().UTC().Format(time.RFC3339)},
					{"1", "1", "Lunch", "350.50", "Food", "", "2025-06-10", time.Now().UTC().Format(time.RFC3339)},
				}
				if err := expenseutils.WriteExpensesCSV(records); err != nil {
					t.Fatalf("write expenses: %v", err)
				}
				got, err := GetAllExpenses()
				if err != nil {
					t.Fatalf("get all: %v", err)
				}
				if len(got) != 1 {
					t.Fatalf("len = %d, want 1", len(got))
				}
			},
		},
		{
			name: "get expenses by user with options",
			run: func(t *testing.T) {
				for _, expense := range sampleExpenses() {
					expense.CreatedAt = time.Now().UTC()
					if err := expenseutils.AppendExpenseCSV(expenseToRecord(expense)); err != nil {
						t.Fatalf("append expense: %v", err)
					}
				}
				got, err := GetExpensesByUserIDWithOptions(1, ExpenseListOptions{Category: "Food"})
				if err != nil {
					t.Fatalf("get by user with options: %v", err)
				}
				if len(got) != 1 || got[0].Category != "Food" {
					t.Fatalf("expenses = %#v", got)
				}
			},
		},
		{
			name: "summary by user",
			run: func(t *testing.T) {
				for _, expense := range sampleExpenses() {
					expense.CreatedAt = time.Now().UTC()
					if err := expenseutils.AppendExpenseCSV(expenseToRecord(expense)); err != nil {
						t.Fatalf("append expense: %v", err)
					}
				}
				got, err := GetExpenseSummaryByUserID(1, "", "")
				if err != nil {
					t.Fatalf("summary by user: %v", err)
				}
				if got.TotalCount != 2 {
					t.Fatalf("count = %d, want 2", got.TotalCount)
				}
			},
		},
		{
			name: "next expense id",
			run: func(t *testing.T) {
				expense := Expense{UserID: 1, Title: "Lunch", Amount: 350.50, Category: "Food", ExpenseDate: "2025-06-10"}
				if err := CreateExpense(&expense); err != nil {
					t.Fatalf("create expense: %v", err)
				}
				if got := GetNextExpenseID(); got != expense.ID+1 {
					t.Fatalf("next id = %d, want %d", got, expense.ID+1)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupModelStorage(t)
			tt.run(t)
		})
	}
}

func TestExpenseRecordHelpers(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	tests := []struct {
		name    string
		record  []string
		wantErr bool
	}{
		{name: "valid record", record: []string{"1", "2", "Lunch", "350.50", "Food", "", "2025-06-10", now.Format(time.RFC3339)}},
		{name: "short record", record: []string{"1"}, wantErr: true},
		{name: "invalid id", record: []string{"bad", "2", "Lunch", "350.50", "Food", "", "2025-06-10", now.Format(time.RFC3339)}, wantErr: true},
		{name: "invalid user id", record: []string{"1", "bad", "Lunch", "350.50", "Food", "", "2025-06-10", now.Format(time.RFC3339)}, wantErr: true},
		{name: "invalid amount", record: []string{"1", "2", "Lunch", "bad", "Food", "", "2025-06-10", now.Format(time.RFC3339)}, wantErr: true},
		{name: "invalid created at", record: []string{"1", "2", "Lunch", "350.50", "Food", "", "2025-06-10", "bad"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := expenseFromRecord(tt.record)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	record := expenseToRecord(Expense{ID: 7, UserID: 3, Title: " Lunch ", Amount: 10, Category: "Food", Note: " note ", ExpenseDate: "2025-06-10", CreatedAt: now})
	if record[0] != strconv.Itoa(7) || record[2] != "Lunch" || record[5] != "note" {
		t.Fatalf("record = %#v", record)
	}
}

func TestExpenseFailures(t *testing.T) {
	tests := []struct {
		name    string
		run     func() error
		wantErr error
	}{
		{
			name: "invalid category on create",
			run: func() error {
				expense := Expense{UserID: 1, Title: "Lunch", Amount: 350.50, Category: "Bad", ExpenseDate: "2025-06-10"}
				return CreateExpense(&expense)
			},
			wantErr: ErrInvalidExpenseCategory,
		},
		{
			name: "missing expense on update",
			run: func() error {
				expense := Expense{ID: 99, UserID: 1, Title: "Lunch", Amount: 350.50, Category: "Food", ExpenseDate: "2025-06-10", CreatedAt: time.Now().UTC()}
				return UpdateExpense(&expense)
			},
			wantErr: ErrExpenseNotFound,
		},
		{
			name:    "missing expense on delete",
			run:     func() error { return DeleteExpense(99, 1) },
			wantErr: ErrExpenseNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupModelStorage(t)
			if err := tt.run(); !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsAllowedCategory(t *testing.T) {
	tests := []struct {
		name     string
		category string
		want     bool
	}{
		{name: "allowed", category: "Food", want: true},
		{name: "not allowed", category: "InvalidCat", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAllowedCategory(tt.category); got != tt.want {
				t.Fatalf("IsAllowedCategory(%q) = %v, want %v", tt.category, got, tt.want)
			}
		})
	}
}
