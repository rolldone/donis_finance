package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rolldone/donisgo/internal/storage"
	"github.com/rolldone/donisgo/internal/uuid"
	"github.com/rolldone/donisgo/plugins/donisfinance/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TransactionHandler handles transaction-related endpoints.
type TransactionHandler struct {
	db    *gorm.DB
	store storage.Store
}

// NewTransactionHandler creates a new handler.
func NewTransactionHandler(db *gorm.DB, store storage.Store) *TransactionHandler {
	return &TransactionHandler{db: db, store: store}
}

// ListTransactions godoc
// GET /admin/plugins/donisfinance/transactions
// GET /api/member/transactions
func (h *TransactionHandler) ListTransactions(c *gin.Context) {
	role := c.GetString("user_type")
	userID := c.GetString("user_id")

	f := services.TransactionFilter{
		Month:     parseInt(c.Query("month"), 0),
		Year:      parseInt(c.Query("year"), 0),
		Type:      c.Query("type"),
		Search:    strings.ToLower(c.Query("q")),
		SortBy:    c.Query("sort_by"),
		SortOrder: c.Query("sort_order"),
		Limit:     parseInt(c.Query("limit"), 50),
		Offset:    parseInt(c.Query("offset"), 0),
	}

	if role == "admin" {
		f.MemberID = c.Query("member_id")
	} else {
		f.MemberID = userID
	}

	result, err := services.ListTransactions(h.db, f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": result.Transactions, "total": result.Total})
}

// CreateTransaction godoc
// POST /api/member/transactions
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	userID := c.GetString("user_id")

	var req struct {
		AccountID   string `json:"account_id"`
		ToAccountID string `json:"to_account_id"`
		CategoryID  string `json:"category_id"`
		Amount      int64  `json:"amount"`
		Type        string `json:"type"`
		Description string `json:"description"`
		Notes       string `json:"notes"`
		Date        string `json:"date"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	var catID, acctID, toAcctID *string
	if req.CategoryID != "" {
		catID = &req.CategoryID
	}
	if req.AccountID != "" {
		acctID = &req.AccountID
	}
	if req.ToAccountID != "" {
		toAcctID = &req.ToAccountID
	}

	result, err := services.CreateTransaction(h.db, userID, acctID, catID, toAcctID, req.Amount, req.Type, req.Description, req.Notes, "", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"transaction": result})
}

// UpdateTransaction godoc
// PUT /api/member/transactions/:id
// PUT /api/admin/transactions/:id
func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("user_type")
	id := c.Param("id")

	var req struct {
		AccountID   string `json:"account_id"`
		ToAccountID string `json:"to_account_id"`
		CategoryID  string `json:"category_id"`
		Amount      int64  `json:"amount"`
		Type        string `json:"type"`
		Description string `json:"description"`
		Notes       string `json:"notes"`
		Date        string `json:"date"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	result, err := services.UpdateTransaction(h.db, id, userID, role, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transaction": result})
}

// UploadAttachment godoc
// POST /api/member/transactions/:id/attachment
// POST /admin/plugins/donisfinance/transactions/:id/attachment
func (h *TransactionHandler) UploadAttachment(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("user_type")
	txID := c.Param("id")

	// Verify transaction exists and belongs to user
	var txOwner struct {
		MemberID string
	}
	if err := h.db.Table("transactions").Select("member_id").Where("id = ?", txID).Scan(&txOwner).Error; err != nil || txOwner.MemberID == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}
	if role != "admin" && txOwner.MemberID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not your transaction"})
		return
	}

	// Read uploaded file
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	defer file.Close()

	// Validate file size (max 10MB)
	if header.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large (max 10MB)"})
		return
	}

	// Build storage path: transactions/{tx_id}/{uuid}{ext}
	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".bin"
	}
	// Only allow common image/document extensions
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true, ".pdf": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true}
	if !allowed[strings.ToLower(ext)] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file type not allowed"})
		return
	}

	storageKey := fmt.Sprintf("transactions/%s/%s%s", txID, uuid.NewString(), ext)

	if err := h.store.Put(context.Background(), storageKey, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	// Update transaction with attachment path
	if err := h.db.Table("transactions").Where("id = ?", txID).Update("attachment_path", storageKey).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update transaction"})
		return
	}

	publicURL, _ := h.store.PublicURL(context.Background(), storageKey)

	c.JSON(http.StatusOK, gin.H{
		"attachment_path": storageKey,
		"url":             publicURL,
	})
}

// GetAttachment godoc
// GET /api/member/transactions/:id/attachment
func (h *TransactionHandler) GetAttachment(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("user_type")
	txID := c.Param("id")

	var tx struct {
		MemberID       string
		AttachmentPath string
	}
	if err := h.db.Table("transactions").Select("member_id, attachment_path").Where("id = ?", txID).Scan(&tx).Error; err != nil || tx.AttachmentPath == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "attachment not found"})
		return
	}
	if role != "admin" && tx.MemberID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not your transaction"})
		return
	}

	reader, err := h.store.Get(context.Background(), tx.AttachmentPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}
	defer reader.Close()

	// Detect content type
	buf := make([]byte, 512)
	n, _ := reader.Read(buf)
	contentType := http.DetectContentType(buf[:n])

	c.DataFromReader(http.StatusOK, -1, contentType, io.MultiReader(
		strings.NewReader(string(buf[:n])),
		reader,
	), nil)
}

// DeleteTransaction godoc
// DELETE /admin/plugins/donisfinance/transactions/:id
// DELETE /api/member/transactions/:id
func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")
	role := c.GetString("user_type")

	// Get full transaction data before deleting (need it for balance reversal)
	var tx struct {
		MemberID       string
		AccountID      *string
		ToAccountID    *string
		Amount         int64
		Type           string
		AttachmentPath string
	}
	if err := h.db.Table("transactions").Select("member_id, account_id, to_account_id, amount, type, attachment_path").Where("id = ?", id).Scan(&tx).Error; err != nil || tx.MemberID == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}

	if role == "member" && tx.MemberID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not your transaction"})
		return
	}

	// Reverse balance before deleting (opposite of what CreateTransaction did)
	switch tx.Type {
	case "income":
		if tx.AccountID != nil && *tx.AccountID != "" {
			h.db.Table("accounts").Where("id = ?", *tx.AccountID).
				Update("balance", gorm.Expr("balance - ?", tx.Amount))
		}
	case "expense":
		if tx.AccountID != nil && *tx.AccountID != "" {
			h.db.Table("accounts").Where("id = ?", *tx.AccountID).
				Update("balance", gorm.Expr("balance + ?", tx.Amount))
		}
	case "transfer":
		if tx.AccountID != nil && *tx.AccountID != "" {
			h.db.Table("accounts").Where("id = ?", *tx.AccountID).
				Update("balance", gorm.Expr("balance + ?", tx.Amount))
		}
		if tx.ToAccountID != nil && *tx.ToAccountID != "" {
			h.db.Table("accounts").Where("id = ?", *tx.ToAccountID).
				Update("balance", gorm.Expr("balance - ?", tx.Amount))
		}
	}

	if err := services.DeleteTransaction(h.db, id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Delete attachment file if exists
	if tx.AttachmentPath != "" {
		h.store.Delete(context.Background(), tx.AttachmentPath)
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// ExportCSV godoc
// GET /api/admin/transactions/export?month=&year=&member_id=
func (h *TransactionHandler) ExportCSV(c *gin.Context) {
	year := parseInt(c.Query("year"), 0)
	month := parseInt(c.Query("month"), 0)
	memberID := c.Query("member_id")
	memberName := c.Query("member_name")

	csvBytes, err := services.ExportTransactionsCSV(h.db, memberID, memberName, month, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filename := fmt.Sprintf("transactions_%d_%02d.csv", year, month)
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Data(http.StatusOK, "text/csv", csvBytes)
}

// GetSummary godoc
// GET /admin/plugins/donisfinance/transactions/summary?member_id=x&month=7&year=2026
// GET /api/member/transactions/summary?month=7&year=2026
func (h *TransactionHandler) GetSummary(c *gin.Context) {
	role := c.GetString("user_type")
	userID := c.GetString("user_id")

	year := parseInt(c.Query("year"), 0)
	month := parseInt(c.Query("month"), 0)
	if year <= 0 || month <= 0 || month > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "year and month are required (1-12)"})
		return
	}

	memberID := userID
	if role == "admin" {
		if m := c.Query("member_id"); m != "" {
			memberID = m
		}
	}

	summary, err := services.GetMonthlySummary(h.db, memberID, year, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetMonthlySeries godoc
// GET /admin/plugins/donisfinance/transactions/monthly?member_id=x&months=6
// GET /api/member/transactions/monthly?months=6
func (h *TransactionHandler) GetMonthlySeries(c *gin.Context) {
	role := c.GetString("user_type")
	userID := c.GetString("user_id")

	months := parseInt(c.Query("months"), 6)
	if months <= 0 || months > 24 {
		months = 6
	}

	memberID := userID
	if role == "admin" {
		if m := c.Query("member_id"); m != "" {
			memberID = m
		}
	}

	series, err := services.GetMonthlySeries(h.db, memberID, months)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"series": series})
}

// ─── Categories ───────────────────────────────────────────────────────────────

// ListCategories godoc
// GET /admin/plugins/donisfinance/categories
// GET /api/member/categories
func (h *TransactionHandler) ListCategories(c *gin.Context) {
	catType := c.Query("type")
	results, err := services.ListCategories(h.db, catType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"categories": results})
}

// ─── Accounts ─────────────────────────────────────────────────────────────────

// ListAccounts godoc
// GET /admin/plugins/donisfinance/accounts?member_id=xxx
// GET /api/member/accounts
func (h *TransactionHandler) ListAccounts(c *gin.Context) {
	role := c.GetString("user_type")
	userID := c.GetString("user_id")

	memberID := userID
	if role == "admin" {
		if m := c.Query("member_id"); m != "" {
			memberID = m
		}
	}

	results, err := services.ListAccounts(h.db, memberID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"accounts": results})
}

// CreateAccount godoc
// POST /api/member/accounts
func (h *TransactionHandler) CreateAccount(c *gin.Context) {
	userID := c.GetString("user_id")

	var req struct {
		Name           string `json:"name"`
		Type           string `json:"type"`
		InitialBalance int64  `json:"initial_balance"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	result, err := services.CreateAccount(h.db, userID, req.Name, req.Type, req.InitialBalance)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"account": result})
}

// UpdateAccount godoc
// PUT /api/member/accounts/:id
func (h *TransactionHandler) UpdateAccount(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")

	var req struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	result, err := services.UpdateAccount(h.db, id, userID, req.Name, req.Type, nil, "")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"account": result})
}

// DeleteAccount godoc
// DELETE /api/member/accounts/:id
func (h *TransactionHandler) DeleteAccount(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")

	if err := services.DeleteAccount(h.db, id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "account deleted"})
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func parseInt(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return v
}
