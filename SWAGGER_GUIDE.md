# 🔄 Guia Completo: Geração de Swagger/OpenAPI em Go

## 📋 Formas Confiáveis para Seu Projeto

Para seu projeto com **Gin + Clean Architecture**, existem **3 opções viáveis**:

---

## 1️⃣ **SWAG (Recomendado) ⭐**

### Por que é confiável?
✅ Mais popular na comunidade Go  
✅ Sincroniza automaticamente com código  
✅ Usa comentários/annotations próximos ao código  
✅ Gera `swagger.json` automaticamente  
✅ Detecta mudanças facilmente  

### Instalação

```bash
go install github.com/swaggo/swag/cmd/swag@latest
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

### Implementação Passo a Passo

#### Passo 1: Adicionar comentários nos handlers

```go
// handler/account.go

// CreateAccount godoc
// @Summary      Criar nova conta
// @Description  Cria uma nova conta bancária com CPF ou CNPJ
// @Tags         Accounts
// @Accept       json
// @Produce      json
// @Param        account body domain.Account true "Dados da conta"
// @Success      201 {object} domain.Account "Conta criada com sucesso"
// @Failure      400 {object} pkg.ErrorResponse "Erro de validação"
// @Router       /v1/accounts [post]
func (h *AccountHandler) CreateAccount(c *gin.Context) {
    // ... seu código
}

// GetAccount godoc
// @Summary      Obter conta por ID
// @Description  Retorna uma conta específica pelo seu ID
// @Tags         Accounts
// @Accept       json
// @Produce      json
// @Param        accountId path int true "ID da conta"
// @Success      200 {object} domain.Account "Conta encontrada"
// @Failure      404 {object} pkg.ErrorResponse "Conta não encontrada"
// @Router       /v1/accounts/{accountId} [get]
func (h *AccountHandler) GetAccount(c *gin.Context) {
    // ... seu código
}
```

#### Passo 2: Adicionar comentários nos modelos de domínio

```go
// domain/account.go

// Account representa uma conta bancária
// @Description Informações da conta bancária
type Account struct {
    ID             int64  `json:"account_id" example:"123"`              // ID da conta
    DocumentNumber string `json:"document_number" example:"12345678901"` // CPF ou CNPJ
}

// Transaction representa uma transação bancária
// @Description Informações da transação
type Transaction struct {
    ID              int64         `json:"transaction_id" example:"456"`
    AccountID       int64         `json:"account_id" example:"123"`
    OperationTypeID OperationType `json:"operation_type_id" example:"1"`
    Amount          float64       `json:"amount" example:"100.50"`
    EventDate       time.Time     `json:"event_date" example:"2024-01-01T12:00:00Z"`
}
```

#### Passo 3: Adicionar comentários na função main

```go
// cmd/main.go

// @title           Account Management API
// @version         1.0
// @description     API de gerenciamento de contas bancárias
// @termsOfService  http://swagger.io/terms/

// @contact.name    API Support
// @contact.url     http://www.swagger.io/support
// @contact.email   support@swagger.io

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @host            localhost:8080
// @BasePath        /
// @schemes         http https

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @externalDocs.description OpenAPI
// @externalDocs.url https://swagger.io/resources/open-api/

func main() {
    // Seu código
}
```

#### Passo 4: Gerar documentação

```bash
# Gerar swagger.json e swagger.yaml
swag init -g cmd/main.go
```

Isso criará:
- `docs/swagger.json`
- `docs/swagger.yaml`
- `docs/docs.go`

#### Passo 5: Servir Swagger UI

```go
// cmd/main.go

import (
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
    _ "github.com/evythrossell/account-management-api/docs" // importante!
)

func main() {
    router := gin.Default()
    
    // Swagger endpoints
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    
    // Suas rotas normais
    setupRoutes(router)
    
    router.Run(":8080")
}
```

#### Passo 6: Acessar Swagger UI

```
http://localhost:8080/swagger/index.html
```

### Exemplo Completo para Seu Projeto

```go
// internal/adapter/http/handler/transaction.go

package handler

import (
    "github.com/gin-gonic/gin"
    "github.com/evythrossell/account-management-api/internal/core/domain"
    "github.com/evythrossell/account-management-api/internal/core/port"
)

type TransactionHandler struct {
    service port.TransactionService
}

// CreateTransaction godoc
// @Summary      Criar transação
// @Description  Cria uma nova transação bancária (crédito/débito)
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param        transaction body CreateTransactionRequest true "Dados da transação"
// @Success      201 {object} domain.Transaction "Transação criada"
// @Failure      400 {object} pkg.ErrorResponse "Erro de validação"
// @Failure      404 {object} pkg.ErrorResponse "Conta não encontrada"
// @Failure      500 {object} pkg.ErrorResponse "Erro interno do servidor"
// @Router       /v1/transactions [post]
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
    var req CreateTransactionRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    result, err := h.service.CreateTransaction(c.Request.Context(), req.AccountID, req.OperationTypeID, req.Amount)
    if err != nil {
        c.Error(err)
        return
    }
    
    c.JSON(201, result)
}

// GetTransaction godoc
// @Summary      Obter transação
// @Description  Retorna detalhes de uma transação específica
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param        transactionId path int64 true "ID da transação"
// @Success      200 {object} domain.Transaction "Transação encontrada"
// @Failure      404 {object} pkg.ErrorResponse "Transação não encontrada"
// @Failure      500 {object} pkg.ErrorResponse "Erro interno"
// @Router       /v1/transactions/{transactionId} [get]
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
    var req GetTransactionRequest
    if err := c.ShouldBindUri(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    result, err := h.service.GetByTransactionID(c.Request.Context(), req.TransactionID)
    if err != nil {
        c.Error(err)
        return
    }
    
    c.JSON(200, result)
}

// Request/Response types
type CreateTransactionRequest struct {
    AccountID       int64   `json:"account_id" binding:"required" example:"123"`
    OperationTypeID int16   `json:"operation_type_id" binding:"required" example:"1"`
    Amount          float64 `json:"amount" binding:"required,gt=0" example:"100.00"`
}

type GetTransactionRequest struct {
    TransactionID int64 `uri:"transactionId" binding:"required"`
}
```

---

## 2️⃣ **GO-SWAGGER (Alternativa)**

### Quando usar?
- Se preferir mais controle sobre a spec
- Projetos mais estabelecidos
- Documentação mais declarativa

### Instalação

```bash
go install github.com/go-swagger/go-swagger/cmd/swagger@latest
```

### Uso como Spec-First

```bash
# Criar spec inicial
swagger init -f json

# Gerar código a partir da spec
swagger generate server -f ./swagger.json -t gen
```

### Vantagens vs Desvantagens

| Aspecto | GO-Swagger | SWAG |
|---------|-----------|------|
| Instalação | 📦 Um tool | 📦 Dois packages |
| Curva Aprendizado | 📈 CurvaAlta | 📈 Mais fácil |
| Code-First | ⚠️ Menos natural | ✅ Muito natural |
| Spec-First | ✅ Excelente | ⚠️ Limitado |
| Comunidade | 📊 Média | 📊 Maior |

---

## 3️⃣ **MANUAL COM OPENAPI (Não Recomendado)**

Editar `openapi.json` manualmente é **arriscado** porque:
- ❌ Fácil ficar desincronizado com código
- ❌ Propenso a erros
- ❌ Difícil manter atualizado
- ❌ Sem validação automática

**Usar apenas se:**
- Projeto muito pequeno
- Equipe reduzida
- Não há atualizações frequentes

---

## 🎯 Recomendação para Seu Projeto

### ✅ **Use SWAG porque:**

1. **Seu projeto já está bem estruturado**
   - Clean Architecture facilita adicionar comentários
   - Handlers estão bem organizados

2. **Sincronização automática**
   - Comentários ficam próximos ao código
   - Mudanças são refletidas imediatamente após rodar `swag init`

3. **Integração com Gin**
   - SWAG foi feito pensando em Gin
   - Integração perfeita com `gin-swagger`

4. **Comunidade forte**
   - Milhares de exemplo
   - Suporte ativo

---

## 📝 Checklist de Implementação

- [ ] Instalar `swag` e `gin-swagger`
- [ ] Adicionar comentários godoc em todos handlers
- [ ] Adicionar comentários nos struct de request/response
- [ ] Adicionar comentários da aplicação em `main.go`
- [ ] Rodar `swag init -g cmd/main.go`
- [ ] Integrar `gin-swagger` no router
- [ ] Testar acessando `/swagger/index.html`
- [ ] Validar spec em `docs/swagger.json`
- [ ] Fazer commit dos arquivos `/docs`

---

## 🚀 Próximos Passos

1. **Implementar comentários SWAG** (30 minutos)
2. **Gerar documentação** (1 minuto)
3. **Validar especificação** (5 minutos)
4. **Publicar como CI/CD artifact** (opcional)
5. **Usar para client generation** (opcional)

---

## ⚠️ Armadilhas Comuns

### ❌ Não faça:
```bash
# Esquecer de importar o package docs
// Falta isto:  _ "github.com/evythrossell/account-management-api/docs"

// Esquecer de rodar swag init após alterações
// Sempre rodar após mudanças nos handlers

// Documentar tipos que não existem
// @Success 200 {object} NonExistentType  // ERRADO
```

### ✅ Sempre faça:
```bash
# Regenerar sempre que mudar handlers
swag init -g cmd/main.go

# Validar a spec gerada
# Abrir docs/swagger.json em swagger editor online
```

---

## 📚 Recursos Adicionais

- **SWAG Documentation**: https://github.com/swaggo/swag
- **Swagger/OpenAPI Spec**: https://swagger.io/specification/
- **Swagger Editor Online**: https://editor.swagger.io/
- **Go-Swagger Manual**: https://goswagger.io/

---

## ✅ Resumo: Qual Escolher?

| Critério | SWAG | GO-Swagger | Manual |
|----------|------|-----------|--------|
| **Confiabilidade** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐ |
| **Facilidade** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐ |
| **Sincronização** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ❌ |
| **Para Seu Projeto** | ✅ **RECOMENDADO** | Alternativa | Evitar |

