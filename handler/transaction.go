package handler

import (
	"bekasiberbagi/response"
	"bekasiberbagi/transaction"
	"bekasiberbagi/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	service transaction.Service
}

func NewTransactionHandler(service transaction.Service) *transactionHandler {
	return &transactionHandler{service}
}

func (h *transactionHandler) GetCampaignTransaction(c *gin.Context) {
	var input transaction.GetCampaignTransactionInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := response.APIResponseFailed("Error uri", http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	input.User = currentUser

	transactions, err := h.service.GetTransactionByCampaignId(input)
	if err != nil {
		response := response.APIResponseFailed(err.Error(), http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := response.APIResponseSuccess("Campaign Transactions", http.StatusOK, transaction.FormatCampaignTransactions(transactions))
	c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) GetUserTransactions(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(user.User)
	userId := currentUser.ID

	transactions, err := h.service.GetTransactionByUserId(userId)

	if err != nil {
		response := response.APIResponseFailed(err.Error(), http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := response.APIResponseSuccess("Campaign Transactions", http.StatusOK, transaction.FormatUserTransactions(transactions))
	c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) CreateTransaction(c *gin.Context) {
	var inputCreateTransaction transaction.CreateTransactionInput

	err := c.ShouldBindJSON(&inputCreateTransaction)

	if err != nil {
		response := response.APIResponseValidationFailed("Cant create transaction", http.StatusUnprocessableEntity, err)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	inputCreateTransaction.User = currentUser

	newTransaction, err := h.service.CreateTransaction(inputCreateTransaction)

	if err != nil {
		data := gin.H{"error": err.Error()}
		response := response.APIResponseFailedWithData("Failed create transaction", http.StatusUnprocessableEntity, data)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := response.APIResponseSuccess("Create transaction success", http.StatusOK, transaction.FormatTransaction(newTransaction))
	c.JSON(http.StatusOK, response)
}
