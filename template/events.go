package template

type Common struct {
	FirstName string
}

// For events : Account_added , AccountEnabled, AccountDisabled, AccountDeleted,
type AccountCommon struct {
	AccountNumber string `key:"ACCOUNT_NUMBER"`
	AccountName   string `key:"ACCOUNT_NAME"`
}

type AccountConnected struct {
	AccountNumber    string `key:"ACCOUNT_NUMBER"`
	ConnectionStatus string `key:"CONNECTION_STATUS"`
}

type AccountConnectionError struct {
	AccountNumber    string `key:"ACCOUNT_NUMBER"`
	ConnectionStatus string `key:"CONNECTION_STATUS"`
	ConnectionError  string `key:"CONNECTION_ERROR"`
}

type CopierCreated struct {
	CopierMaster string `key:"COPIER_MASTER"`
	CopierSlave  string `key:"COPIER_SLAVE"`
	CopierType   string `key:"COPIER_TYPE"`
	CopierRisk   string `key:"COPIER_RISK"`
}

// CopierCommon struct For events : Copier Enabled , Copier Disabled , Copier Modified , Copier Deleted
type CopierCommon struct {
	CopierMaster string `key:"COPIER_MASTER"`
	CopierSlave  string `key:"COPIER_SLAVE"`
}

type TradeCopiedSuccessfully struct {
	CopierMaster       string `key:"COPIER_MASTER"`
	CopierSlave        string `key:"COPIER_SLAVE"`
	CopierMasterTicket string `key:"COPIER_MASTER_TICKET"`
	CopierSlaveTicket  string `key:"COPIER_SLAVE_TICKET"`
	CopierMasterSymbol string `key:"COPIER_MASTER_SYMBOL"`
	CopierSlaveSymbol  string `key:"COPIER_SLAVE_SYMBOL"`
}

type TradeModifiedSuccessfully struct {
	CopierMaster string `key:"COPIER_MASTER"`
	CopierSlave  string `key:"COPIER_SLAVE"`
}

type TradeCopyFailure struct {
	CopyMaster  string `key:"COPIER_MASTER"`
	CopySlave   string `key:"COPIER_SLAVE"`
	CopierError string `key:"COPIER_ERROR"`
}
