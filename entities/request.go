package entities

type CreateAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

type CreateUserRequest struct {
	Username        string `json:"username"  binding:"required,min=6,max=16,alphanum"`
	FullName        string `json:"full_name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6,max=16"`
	PasswordConfirm string `json:"password_confirm" binding:"required,eqfield=Password"`
}

type UpdateUserRequest struct {
	FullName    string `json:"full_name" binding:"alpha,required"`
	Email       string `json:"email" binding:"email"`
	OldPassword string `json:"old_password" binding:"min=6,max=16,required_with=NewPassword"`
	NewPassword string `json:"new_password" binding:"min=6,max=16,require"`
}

type CreateTransferRequest struct {
	FromAccountID int64 `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64 `json:"to_account_id" binding:"required,min=1"`
	Amount        int64 `json:"amount" binding:"required,gte=1"`
}

type LoginUserRequest struct {
	Username string `json:"username" binding:"required,min=6,max=16,alphanum"`
	Password string `json:"password" binding:"required,min=6,max=16"`
}

type GetAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type DeleteAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type RestoreAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type RenewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type GetTransferRequest struct {
	AccountID int64 `uri:"id" binding:"required,min=1"`
}
